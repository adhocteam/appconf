/// <reference path="typings/tsd.d.ts" />
/// <reference path="interfaces/interfaces.ts" />
/// <reference path="services/ConfigVarStorage.ts" />
/// <reference path="controllers/ConfigVarsCtrl.ts" />
/// <reference path="controllers/ApplicationCtrl.ts" />
/// <reference path="controllers/EnvironmentCtrl.ts" />

module settings {
    "use strict";

    var mod = angular.module("settings", [
        "ngRoute",
        "ngAnimate"
    ])
        .service("configVarStorage", RemoteConfigVarStorage)
        .controller("ApplicationCtrl", ApplicationCtrl)
        .controller("EnvironmentCtrl", EnvironmentCtrl)
        .controller("ConfigVarsCtrl", ConfigVarsCtrl)
        .config(["$routeProvider", ($routeProvider) => {
            $routeProvider.
            when("/", {
                templateUrl: "partials/app-list.html",
                controller: "ApplicationCtrl"
            }).
            when("/:app", {
                templateUrl: "partials/env-list.html",
                controller: "EnvironmentCtrl"
            }).
            when("/:app/:env", {
                templateUrl: "partials/configvars.html",
                controller: "ConfigVarsCtrl"
            }).
            otherwise({redirectTo: "/"})
        }]);
}
