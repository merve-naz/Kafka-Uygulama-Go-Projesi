// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/trigger-kafka": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "manually trigger kafka with payload",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Certificate"
                ],
                "summary": "manually trigger kafka with payload",
                "parameters": [
                    {
                        "description": "certification dto",
                        "name": "certificate",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Certificate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "kafka triggered manually success",
                        "schema": {
                            "$ref": "#/definitions/handlers.RespondJson"
                        }
                    },
                    "400": {
                        "description": "invalid certificate info for trigger",
                        "schema": {
                            "$ref": "#/definitions/handlers.RespondJson"
                        }
                    },
                    "500": {
                        "description": "internal server error while trigger kafka",
                        "schema": {
                            "$ref": "#/definitions/handlers.RespondJson"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.RespondJson": {
            "type": "object",
            "properties": {
                "intent": {
                    "type": "string",
                    "example": "bbrn:::certificateservice:::/upload"
                },
                "message": {},
                "status": {
                    "type": "boolean",
                    "example": true
                }
            }
        },
        "models.Certificate": {
            "type": "object",
            "required": [
                "avatar",
                "badge_owner",
                "badge_title",
                "completed_at",
                "name",
                "registered_at",
                "registration_uid",
                "slug",
                "title",
                "url_slug"
            ],
            "properties": {
                "avatar": {
                    "type": "string"
                },
                "badge_owner": {
                    "type": "string"
                },
                "badge_title": {
                    "type": "string"
                },
                "completed_at": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "registered_at": {
                    "type": "string"
                },
                "registration_uid": {
                    "type": "string"
                },
                "slug": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "url_slug": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.12",
	Host:             "",
	BasePath:         "/api-certificates",
	Schemes:          []string{"http", "https"},
	Title:            "BB Certificate Generator Service API",
	Description:      "bb.app.certificateservice: microservice for certificate.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
