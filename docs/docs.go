// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Marlen E. Satriani",
            "email": "marlendotedots@gmail.com"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/MartinHeinz/go-project-blueprint/blob/master/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Login as a user",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userRequest.Login"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userResource.User"
                        }
                    }
                }
            }
        },
        "/auth/login/google": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Login with Google",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userRequest.Google"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userResource.User"
                        }
                    }
                }
            }
        },
        "/auth/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Register a new user",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userRequest.Register"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/userResource.User"
                        }
                    }
                }
            }
        },
        "/order": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order"
                ],
                "summary": "Create Order",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/orderRequest.Order"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/orderResource.Order"
                        }
                    }
                }
            }
        },
        "/order/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order"
                ],
                "summary": "Get Order By ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderResource.Order"
                        }
                    }
                }
            }
        },
        "/order/{id}/cancel": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order"
                ],
                "summary": "Cancel Order",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Order ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderResource.Order"
                        }
                    }
                }
            }
        },
        "/orders": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order"
                ],
                "summary": "Get Orders Me",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderResource.OrderList"
                        }
                    }
                }
            }
        },
        "/orders/all": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order"
                ],
                "summary": "Get Orders All",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderResource.OrderList"
                        }
                    }
                }
            }
        },
        "/orders/user/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Order"
                ],
                "summary": "Get Orders By User",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/orderResource.OrderList"
                        }
                    }
                }
            }
        },
        "/profile/me": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get the current logged-in user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userResource.User"
                        }
                    }
                }
            }
        },
        "/profile/update": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update the current logged-in user",
                "parameters": [
                    {
                        "description": "Body",
                        "name": "_",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/userRequest.Update"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/userResource.User"
                        }
                    }
                }
            }
        },
        "/user/{id}": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get a user by ID",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userResource.Base"
                        }
                    }
                }
            }
        },
        "/users": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get all users",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userResource.Base"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "orderRequest.Order": {
            "type": "object",
            "required": [
                "addressId",
                "collectDate",
                "estimateDate",
                "orderType",
                "price",
                "serviceType"
            ],
            "properties": {
                "addressId": {
                    "type": "integer"
                },
                "collectDate": {
                    "type": "string"
                },
                "estimateDate": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "note": {
                    "type": "string"
                },
                "orderType": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "serviceType": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "transactionId": {
                    "type": "integer"
                }
            }
        },
        "orderResource.Base": {
            "type": "object",
            "properties": {
                "addressId": {
                    "type": "integer"
                },
                "collectDate": {
                    "type": "string"
                },
                "createdAt": {
                    "type": "string"
                },
                "estimateDate": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "note": {
                    "type": "string"
                },
                "orderType": {
                    "type": "string"
                },
                "price": {
                    "type": "number"
                },
                "serviceType": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "transactionId": {
                    "type": "integer"
                },
                "updatedAt": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "orderResource.Hypermedia": {
            "type": "object",
            "properties": {
                "cancel": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "create": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "self": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                }
            }
        },
        "orderResource.Order": {
            "type": "object",
            "properties": {
                "_links": {
                    "$ref": "#/definitions/orderResource.Hypermedia"
                },
                "message": {
                    "type": "string"
                },
                "order": {
                    "$ref": "#/definitions/orderResource.Base"
                }
            }
        },
        "orderResource.OrderList": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                },
                "orders": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/orderResource.Base"
                    }
                }
            }
        },
        "userRequest.Google": {
            "type": "object",
            "properties": {
                "fcmToken": {
                    "type": "string"
                },
                "idToken": {
                    "type": "string"
                }
            }
        },
        "userRequest.Login": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "fcmToken": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "userRequest.Register": {
            "type": "object",
            "required": [
                "email",
                "firstName",
                "lastName",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 130,
                    "minLength": 3
                }
            }
        },
        "userRequest.Update": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 130,
                    "minLength": 3
                }
            }
        },
        "userResource.Base": {
            "type": "object",
            "properties": {
                "createdAt": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "image": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "userResource.User": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/userResource.Base"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Washit API",
	Description:      "Swagger for washit app.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
