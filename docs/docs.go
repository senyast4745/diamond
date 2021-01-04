// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/diamond": {
            "post": {
                "description": "method to get some diamonds from a mine. If mine is empty, this method delete this mine.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Get Diamonds from Mine.",
                "parameters": [
                    {
                        "description": "the name of the mine from which we want to extract diamonds",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.mineName"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.mineCount"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "get the status of server.",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Show the status of server.",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/mine": {
            "get": {
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Show all mines.",
                "responses": {
                    "200": {
                        "description": "list of all mines",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/repository.Mine"
                            }
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add new mine.",
                "parameters": [
                    {
                        "description": "new mine model",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/repository.Mine"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "description": "gets all the diamonds from the mine and closes it",
                "consumes": [
                    "*/*"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Closes mine.",
                "parameters": [
                    {
                        "description": "the name of the mine to be closed",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.mineName"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.mineCount"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.mineCount": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                }
            }
        },
        "api.mineName": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        },
        "repository.Mine": {
            "type": "object",
            "properties": {
                "diamond_count": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:7171",
	BasePath:    "/",
	Schemes:     []string{"http"},
	Title:       "Diamond And Mine",
	Description: "This is diamond play.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}