SOURCES = app.js style.css index.html bower_components partials

app.js: app.ts controllers/ConfigVarsCtrl.ts interfaces/interfaces.ts services/ConfigVarStorage.ts controllers/ApplicationCtrl.ts controllers/EnvironmentCtrl.ts
	tsc --sourcemap --out $@ $<

appconf-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $@

rpm: app.js appconf-linux-amd64
	fpm -s dir -t rpm -n appconf --rpm-os linux -f --prefix /usr/local/appconf ./appconf-linux-amd64=appconf $(SOURCES)

.PHONY : rpm

clean:
	-rm -f app.js app.js.map appconf-1.0-1.x86_64.rpm appconf-linux-amd64
