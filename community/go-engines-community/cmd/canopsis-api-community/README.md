# Community API 

- [Deployment](#deployment)
- [Documentation](#documentation)

## Deployment

Use [configuration dir](../../config/api/security). Provide the path to the config through `-c` argument.

Use `-port` to change API port.

Use `-d` to enable debug mode.

Use `-docs` to enable Open API documentation on `/swagger` endpoint.

## Documentation

How to add new endpoint to Open API documentation: 

1. For `POST` and `PUT` requests add a model to a handler comment.
    ```go
    // @Param request body CreateRequest true "request"
    func (a *api) Create(c *gin.Context) {
    
    }
    ```

2. Add new response model to a handler comment.
    ```go
    // @Success 201 {object} Rule
    func (a *api) Create(c *gin.Context) {
        
    }
    ```

3. Generate [Open API schemas](../../lib/api/docs/schemas_swagger.yaml) documentation for Golang structs using `go:generate` in [main.go](./main.go).

4. Use generated schemas in [Open API v3](../../lib/api/docs/swagger.yaml) documentation.
    ```yaml
   paths:
      /rules:
        post:
          requestBody:
            required: true
            content:
              application/json:
                schema:
                  $ref: 'schemas_swagger.yaml#/definitions/rule.CreateRequest'
          responses:
            201:
              content:
                application/json:
                  schema:
                    $ref: 'schemas_swagger.yaml#/definitions/rule.Rule'
    ```
