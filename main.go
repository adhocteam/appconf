package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bmizerany/pat"
	"github.com/gorilla/handlers"
)

var (
	region = "us-east-1"
	bucket = "pet-config-vars"
)

type app struct {
	Short  string `json:"shortname"`  // short name, like "marketplace-api"
	Pretty string `json:"prettyname"` // pretty-print name, like "Marketplace API"
}

var (
	marketplaceAPI = app{"marketplace-api", "Marketplace API"}
	windowshop     = app{"windowshop", "WindowShop"}
	plancompare    = app{"plancompare", "Plan Compare"}
	haproxy        = app{"haproxy", "HAProxy"}
	petapi         = app{"petapi", "PET API"}
)

var apps = []app{marketplaceAPI, windowshop, plancompare, haproxy, petapi}

func appByName(name string) app {
	switch name {
	case "marketplace-api":
		return marketplaceAPI
	case "windowshop":
		return windowshop
	case "plancompare":
		return plancompare
	case "haproxy":
		return haproxy
	case "petapi":
		return petapi
	default:
		return app{}
	}
}

type env string

var (
	prod  = env("prod")
	imp1b = env("imp1b")
	imp1a = env("imp1a")
	test  = env("test")
	dev   = env("dev")
)

var envs = []env{prod, imp1b, imp1a, test, dev}

type cfgvar struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

func listApps(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(struct {
		Apps []app `json:"apps"`
	}{
		apps,
	}); err != nil {
		log.Printf("encoding JSON: %v", err)
		http.Error(w, http.StatusText(500), 500)
	}
}

func listEnvs(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(struct {
		App  app   `json:"app"`
		Envs []env `json:"envs"`
	}{
		appByName(r.FormValue(":app")),
		envs,
	}); err != nil {
		log.Printf("encoding JSON: %v", err)
		http.Error(w, http.StatusText(500), 500)
	}
}

func getS3Object(svc *s3.S3, key string) (*s3.GetObjectOutput, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	return svc.GetObject(params)
}

func varNameFromS3Key(key string) string {
	return path.Base(key)
}

func listVars(w http.ResponseWriter, r *http.Request) {
	a := appByName(r.FormValue(":app"))
	e := env(r.FormValue(":env"))

	svc := s3.New(session.New(), &aws.Config{Region: aws.String(region)})
	prefix := a.Short + "/" + string(e)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		log.Printf("listing objects: %v", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	vars := make([]cfgvar, 0, len(resp.Contents))

	for _, obj := range resp.Contents {
		if strings.HasSuffix(*obj.Key, "/") {
			log.Printf("skipping directory %q", *obj.Key)
			continue
		}
		varObj, err := getS3Object(svc, *obj.Key)
		if err != nil {
			log.Printf("getting object %s: %v", *obj.Key, err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		if varObj.DeleteMarker != nil && *varObj.DeleteMarker {
			log.Printf("object %s was delete marker, skipping", *obj.Key)
			continue
		}
		v := cfgvar{Name: varNameFromS3Key(*obj.Key)}
		valbytes, err := ioutil.ReadAll(varObj.Body)
		if err != nil {
			log.Printf("reading object body: %v", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
		v.Val = string(valbytes)
		vars = append(vars, v)
	}

	if err := json.NewEncoder(w).Encode(struct {
		App  app      `json:"app"`
		Env  env      `json:"env"`
		Vars []cfgvar `json:"vars"`
	}{
		a,
		e,
		vars,
	}); err != nil {
		log.Printf("encoding JSON: %v", err)
		http.Error(w, http.StatusText(500), 500)
	}
}

var envVarNameRE = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z_0-9]*`)

func createVar(w http.ResponseWriter, r *http.Request) {
	a := appByName(r.FormValue(":app"))
	e := env(r.FormValue(":env"))

	svc := s3.New(session.New(), &aws.Config{Region: aws.String(region)})
	prefix := a.Short + "/" + string(e)

	if !envVarNameRE.MatchString(r.FormValue("name")) {
		log.Printf("invalid env var name %q", r.FormValue("name"))
		http.Error(w, http.StatusText(400), 400)
		return
	}

	v := cfgvar{Name: r.FormValue("name"), Val: r.FormValue("val")}

	key := prefix + "/" + v.Name

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(v.Val),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		log.Printf("putting s3 object %s: %v", key, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	log.Printf("created env var \"%s=%s\", version: %v", v.Name, v.Val, resp.VersionId)

	w.WriteHeader(201)
}

func updateVar(w http.ResponseWriter, r *http.Request) {
	a := appByName(r.FormValue(":app"))
	e := env(r.FormValue(":env"))

	svc := s3.New(session.New(), &aws.Config{Region: aws.String(region)})
	prefix := a.Short + "/" + string(e)

	if !envVarNameRE.MatchString(r.FormValue(":name")) {
		log.Printf("invalid env var name %q", r.FormValue("name"))
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// TODO: test for existing first

	v := cfgvar{Name: r.FormValue(":name"), Val: r.FormValue("val")}

	key := prefix + "/" + v.Name

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   strings.NewReader(v.Val),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		log.Printf("putting s3 object %s: %v", key, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	log.Printf("updated env var \"%s=%s\", version: %v", v.Name, v.Val, resp.VersionId)
}

func deleteVar(w http.ResponseWriter, r *http.Request) {
	a := appByName(r.FormValue(":app"))
	e := env(r.FormValue(":env"))

	svc := s3.New(session.New(), &aws.Config{Region: aws.String(region)})
	prefix := a.Short + "/" + string(e)

	if !envVarNameRE.MatchString(r.FormValue(":name")) {
		log.Printf("invalid env var name %q", r.FormValue("name"))
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// TODO: test for existing first

	v := cfgvar{Name: r.FormValue(":name"), Val: r.FormValue("val")}

	key := prefix + "/" + v.Name

	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	resp, err := svc.DeleteObject(params)
	if err != nil {
		log.Printf("deleting s3 object %s: %v", key, err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	log.Printf("deleted env var %q (key: %s), version: %v", v.Name, key, resp.VersionId)
}

func main() {
	listenAddr := flag.String("l", ":8080", "address to listen on")

	flag.Parse()

	m := pat.New()
	m.Get("/a/apps", http.HandlerFunc(listApps))
	m.Get("/a/:app", http.HandlerFunc(listEnvs))
	m.Get("/a/:app/:env", http.HandlerFunc(listVars))
	m.Post("/a/:app/:env", http.HandlerFunc(createVar))
	m.Put("/a/:app/:env/:name", http.HandlerFunc(updateVar))
	m.Add("DELETE", "/a/:app/:env/:name", http.HandlerFunc(deleteVar))

	http.Handle("/a/", m)
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Printf("listening on %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, handlers.CombinedLoggingHandler(os.Stderr, http.DefaultServeMux)))
}
