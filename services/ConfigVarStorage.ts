module settings {
    "use strict";

    // export class LocalConfigVarStorage implements IConfigVarStorage {
    //     STORAGE_ID = "pet-config-vars-dev";
        
    //     load(app: Application, env: string, callback: (error: string, vars: Var[]) => void) {
    //         callback("", JSON.parse(localStorage.getItem(this.STORAGE_ID) || "[]"));
    //     }

    //     save(app: Application, env: string, vars: Var[], callback: (error: string) => void) {
    //         localStorage.setItem(this.STORAGE_ID, angular.toJson(env));
    //     }
    // }

    interface IConfigVarsResponse {
        app: Application;
        env: string;
        vars: Var[];
    }

    export class RemoteConfigVarStorage implements IConfigVarStorage {
        public static $inject = ["$http"];
            
        constructor(private $http: ng.IHttpService) {}

        loadAll(app: Application, env: string, callback: (error: string, vars: Var[]) => void) {
            this.$http.get(`/a/${app.shortname}/${env}`).success((data: IConfigVarsResponse) => {
                callback("", data.vars);
            });
        }

        create(app: Application, env: string, v: Var, callback: (error: string, vars: Var[]) => void) {
            this.$http({
                method: "POST",
                url: `/a/${app.shortname}/${env}`,
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                data: {name: v.name, val: v.val},
                transformRequest: (obj) => {
                    var params = [];
                    for (var p in obj) {
                        params.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
                    }
                    return params.join("&");
                }
            }).success(() => {
                this.loadAll(app, env, callback);
            });
        }

        update(app: Application, env: string, v: Var, callback: (error: string, vars: Var[]) => void) {
            this.$http({
                method: "PUT",
                url: `/a/${app.shortname}/${env}/${v.name}`,
                headers: {"Content-Type": "application/x-www-form-urlencoded"},
                data: {val: v.val},
                transformRequest: (obj) => {
                    var params = [];
                    for (var p in obj) {
                        params.push(encodeURIComponent(p) + "=" + encodeURIComponent(obj[p]));
                    }
                    return params.join("&");
                }
            }).success(() => {
                this.loadAll(app, env, callback);
            });
        }

        delete(app: Application, env: string, name: string, callback: (error: string, vars: Var[]) => void) {
            this.$http.delete(`/a/${app.shortname}/${env}/${name}`).success(() => {
                this.loadAll(app, env, callback);
            });
        }
    }
}
