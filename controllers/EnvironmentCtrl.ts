module settings {
    interface IEnvsResponse {
        app: Application;
        envs: string[];
    }

    export class EnvironmentCtrl {
        public static $inject = [
            "$scope",
            "$routeParams",
            "$http"
        ];

        constructor(private $scope: IEnvironmentScope, $routeParams: ng.route.IRouteService, $http: ng.IHttpService) {
            $http.get(`/a/${$routeParams['app']}`).success((data: IEnvsResponse) => {
                $scope.app = data.app;
                $scope.envs = data.envs;
            });
        }
    }
}
