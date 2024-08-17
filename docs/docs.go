// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/auth/login": {
            "post": {
                "description": "accepts json with user info and authorize him",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authorization"
                ],
                "summary": "login the user",
                "parameters": [
                    {
                        "description": "account info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.LoginRequest"
                        }
                    },
                    {
                        "description": "session info",
                        "name": "output",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.LoginResponse"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Login was successful"
                    },
                    "400": {
                        "description": "error: Invalid to insert token"
                    },
                    "403": {
                        "description": "error: Invalid email or password"
                    },
                    "422": {
                        "description": "error: Invalid password size"
                    },
                    "500": {
                        "description": "error: Invalid to create token"
                    }
                }
            }
        },
        "/api/auth/logout": {
            "delete": {
                "description": "accepts json with refresh token and delete session",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authorization"
                ],
                "summary": "delete user session",
                "parameters": [
                    {
                        "description": "session info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.LogoutRequest"
                        }
                    },
                    {
                        "description": "response info",
                        "name": "output",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.LogoutResponse"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Logout was successful"
                    },
                    "400": {
                        "description": "error: Failed to read body"
                    },
                    "500": {
                        "description": "error: Invalid to remove session"
                    }
                }
            }
        },
        "/api/auth/ping": {
            "get": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "example"
                ],
                "summary": "ping example",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "plain"
                        }
                    }
                }
            }
        },
        "/api/auth/refresh": {
            "post": {
                "description": "accept json and refresh user tokens",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authorization"
                ],
                "summary": "refresh user's tokens",
                "parameters": [
                    {
                        "description": "session info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RefreshTokensRequest"
                        }
                    },
                    {
                        "description": "response info",
                        "name": "output",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RefreshTokensResponse"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: RefreshTokens was successful"
                    },
                    "401": {
                        "description": "error: Invalid to get refresh token from cookie"
                    },
                    "500": {
                        "description": "error: Invalid to create token"
                    }
                }
            }
        },
        "/api/auth/registration": {
            "post": {
                "description": "accepts json with user info and registers him",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authorization"
                ],
                "summary": "register user",
                "parameters": [
                    {
                        "description": "account info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RegisterRequest"
                        }
                    },
                    {
                        "description": "response info",
                        "name": "output",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/requests.RegisterResponse"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "message: Registration was successful"
                    },
                    "400": {
                        "description": "error: Failed to read body"
                    },
                    "409": {
                        "description": "error: email or channel already been use"
                    },
                    "422": {
                        "description": "error: Failed create password, because it exceeds the character limit or backwards"
                    },
                    "500": {
                        "description": "error: Failed to hash password. Please, try again later"
                    }
                }
            }
        },
        "/api/upload": {
            "post": {
                "description": "accepts file and upload it",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "upload"
                ],
                "summary": "upload a JSON",
                "responses": {
                    "200": {
                        "description": "message: Uploade was successful"
                    },
                    "400": {
                        "description": "error: Only JSON file accepted"
                    }
                }
            }
        }
    },
    "definitions": {
        "requests.LoginRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "fingerprint": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "requests.LoginResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "requests.LogoutRequest": {
            "type": "object",
            "properties": {
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "requests.LogoutResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "requests.RefreshTokensRequest": {
            "type": "object",
            "properties": {
                "fingerprint": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "requests.RefreshTokensResponse": {
            "type": "object",
            "properties": {
                "accessToken": {
                    "type": "string"
                },
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "refreshToken": {
                    "type": "string"
                }
            }
        },
        "requests.RegisterRequest": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "fingerprint": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "requests.RegisterResponse": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Smart Toy API",
	Description:      "API Server for Smart Toy",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
