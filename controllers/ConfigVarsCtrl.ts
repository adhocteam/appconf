module settings {
    export class ConfigVarsCtrl {
        public static $inject = [
            "$scope",
            "$routeParams",
            "$http",
            "configVarStorage"
        ];

        constructor(
            private $scope: IConfigVarsScope,
            $routeParams: ng.route.IRouteService,
            $http: ng.IHttpService,
            private configVarStorage: IConfigVarStorage
        ) {
            $scope.loading = true;

            var appShortname = $routeParams["app"];

            $http.get(`/a/${appShortname}`).success((data) => {
                $scope.app = data["app"];
                $scope.env = $routeParams["env"];

                configVarStorage.loadAll($scope.app, $scope.env, (error: string, vars: Var[]) => {
                    $scope.loading = false;
                    
                    if (error) {
                        console.error(error);
                        return;
                    }

                    $scope.vars = vars;
                });
            });

            // 'vm' stands for 'view model', a way to access these
            // methods/properties via the scope in the view templates
            $scope.vm = this;

            this.reset();
        }

        updateVars(vars: Var[]) {
            this.$scope.vars = vars;
            this.$scope.vars.sort((a, b) => {
                if (a.name > b.name) {
                    return 1;
                } else if (a.name < b.name) {
                    return -1;
                }
                return 0;
            });
        }

        reset() {
            this.$scope.newvar = this.emptyVar();
        }

        emptyVar(): Var {
            return {name: "", val: ""};
        }

        add() {
            var replaced = false;
            for (var i = 0; i < this.$scope.vars.length; i++) {
                if (this.$scope.vars[i].name === this.$scope.newvar.name) {
                    this.$scope.vars[i].val = this.$scope.newvar.val;
                    replaced = true;
                }
            }

            if (!replaced) {
                this.$scope.vars.push(this.$scope.newvar);
            }

            this.configVarStorage.create(this.$scope.app, this.$scope.env, this.$scope.newvar, (error: string, vars: Var[]) => {
                if (error) {
                    console.error(error);
                    return;
                }
                this.$scope.vars = vars;
            });

            this.reset();
        }

        startEditing(v: Var) {
            v.original = v.val;
        }

        cancelEditing(v: Var) {
            delete v.original;
        }

        isEditing(v: Var): boolean {
            return "original" in v;
        }

        delete(v: Var) {
            var found = false;
            for (var i = 0; i < this.$scope.vars.length; i++) {
                if (this.$scope.vars[i].name === v.name) {
                    this.$scope.vars.splice(i, 1);
                    found = true;
                    break;
                }
            }

            if (found) {
                this.configVarStorage.delete(this.$scope.app, this.$scope.env, v.name, (error: string, vars: Var[]) => {
                    if (error) {
                        console.error(error);
                        return;
                    }
                    this.$scope.vars = vars;
                });
            }
        }            

        update(v: Var) {
            var found = false;
            for (var i = 0; i < this.$scope.vars.length; i++) {
                if (this.$scope.vars[i].name === v.name) {
                    this.$scope.vars[i].val = v.val;
                    delete v.original;
                    found = true;
                    break;
                }
            }

            if (found) {
                this.configVarStorage.update(this.$scope.app, this.$scope.env, v, (error: string, vars: Var[]) => {
                    if (error) {
                        console.error(error);
                        return;
                    }
                    this.$scope.vars = vars;
                });
            }
        }
    }
}
