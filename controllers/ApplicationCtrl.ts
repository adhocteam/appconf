module settings {
    interface IAppsResponse {
        apps: Application[];
    }

    export class ApplicationCtrl {
        public static $inject = ["$scope", "$http"];

        constructor(private $scope: IApplicationScope, $http: ng.IHttpService) {
            $http.get("/a/apps").success((data: IAppsResponse) => {
                $scope.apps = data.apps;
            });
        }
    }
}
