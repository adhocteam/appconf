app.js: app.ts controllers/ConfigVarsCtrl.ts interfaces/interfaces.ts services/ConfigVarStorage.ts controllers/ApplicationCtrl.ts controllers/EnvironmentCtrl.ts
	tsc --sourcemap --out $@ $<

appconf-linux-amd64: main.go
	GOPATH=$$GOPATH:$$HOME/work/marketplace:$$HOME/work/marketplace/vendor GOOS=linux GOARCH=amd64 go build -o $@ $<

SOURCES = app.js style.css index.html bower_components partials

rpm: app.js appconf-linux-amd64
	fpm -s dir -t rpm -n appconf --rpm-os linux -f --prefix /usr/local/appconf ./appconf-linux-amd64=appconf $(SOURCES)

.PHONY : release
