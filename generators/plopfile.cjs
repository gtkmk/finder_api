module.exports = (plop) => {
    plop.addHelper('ifCond', function (v1, operator, v2, options) {
        switch (operator) {
            case '===':
                return (v1 === v2) ? options.fn(this) : options.inverse(this);
            case '!==':
                return (v1 !== v2) ? options.fn(this) : options.inverse(this);
            case '<':
                return (v1 < v2) ? options.fn(this) : options.inverse(this);
            case '<=':
                return (v1 <= v2) ? options.fn(this) : options.inverse(this);
            case '>':
                return (v1 > v2) ? options.fn(this) : options.inverse(this);
            case '>=':
                return (v1 >= v2) ? options.fn(this) : options.inverse(this);
            case '&&':
                return (v1 && v2) ? options.fn(this) : options.inverse(this);
            case '||':
                return (v1 || v2) ? options.fn(this) : options.inverse(this);
            default:
                return options.inverse(this);
        }
    });

    function generateTemplateConfig() {
        return {
            description: `Create Handler`,
            prompts: [
                {
                    type: 'input',
                    name: 'model',
                    message: `What is your model?`,
                },
                {
                    type: 'list',
                    name: 'requestType',
                    message: 'What is your request type?',
                    choices: ['Find', 'Create', 'Update', 'Delete'],
                },
                {
                    type: 'input',
                    name: 'fileName',
                    message: `What is your file name?`,
                },
            ],
            actions: function(data) {
                return [
                    {
                        type: 'add',
                        path: `../adapter/http/handlers/{{camelCase model}}Handler/{{requestType}}{{pascalCase fileName}}Handler.go`,
                        templateFile: `templates/handler/templateHandler${data.requestType}.go.hbs`
                    },
                    {
                        type: 'add',
                        path: `../adapter/http/handlers/{{camelCase model}}Handler/{{requestType}}{{pascalCase fileName}}Handler.go`,
                        path: `../core/usecase/{{camelCase model}}/{{requestType}}{{pascalCase fileName}}{{pascalCase model}}.go`,
                        templateFile: `templates/useCase/templateUseCase${data.requestType}.go.hbs`
                    },
                    {
                        type: 'append',
                        path: `../adapter/http/routesConstants/Routes.go`,
                        pattern: `// === Route marker ===`,
                        template: `	{{pascalCase requestType}}{{pascalCase fileName}}RouteConst = "/{{camelCase model}}/{{camelCase fileName}}"`,
                    },
                    {
                        type: 'append',
                        path: `../adapter/http/routes/{{camelCase model}}/{{pascalCase model}}Routes.go`,
                        pattern: `// === Route constants marker ===`,
                        template: `	{{requestType}}{{pascalCase model}}{{pascalCase fileName}}Const string = "{{requestType}}{{pascalCase model}}"`,
                    },
                    {
                        type: 'append',
                        path: `../adapter/http/routes/{{camelCase model}}/{{pascalCase model}}Routes.go`,
                        pattern: `// === Register route marker ===`,
                        templateFile: `templates/route/templateRegisterRouteMarker.go.hbs`
                    },
                    {
                        type: 'append',
                        path: `../adapter/http/routes/{{camelCase model}}/{{pascalCase model}}Routes.go`,
                        pattern: `// === Register handler marker ===`,
                        templateFile: `templates/route/templateDefineHandlerRoute{{requestType}}.go.hbs`
                    },
                ]
            }
        };
    }

    plop.setGenerator('handler', generateTemplateConfig());
};
