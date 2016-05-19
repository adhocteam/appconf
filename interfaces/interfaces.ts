module settings {
    export interface Var {
        name: string;
        val: string;
        original?: string;
        delete?: boolean;
    }

    export interface IConfigVarsScope extends ng.IScope {
        app: Application;
        env: string;
        vars: Var[];
        newvar?: Var;
        vm: ConfigVarsCtrl;
        loading: boolean;
    }

    export interface IConfigVarStorage {
        loadAll(app: Application, env: string, callback: (error: string, vars: Var[]) => void)
        create(app: Application, env: string, v: Var, callback: (error: string, vars: Var[]) => void)
        update(app: Application, env: string, v: Var, callback: (error: string, vars: Var[]) => void)
        delete(app: Application, env: string, name: string, callback: (error: string, vars: Var[]) => void)
    }

    export interface Application {
        shortname: string;
        prettyname: string;
    }

    export interface IApplicationScope extends ng.IScope {
        apps: Application[];
    }

    export interface IEnvironmentScope extends ng.IScope {
        app: Application;
        envs: string[];
    }
}
