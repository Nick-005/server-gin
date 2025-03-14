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
        "/all/vac": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/emp": {
            "post": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/emp/:id": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/emp/vacs/id": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/token/check": {
            "get": {
                "description": "Позволяет выдать новый токен пользователю, если у него нету актуального 'Bearer Token' или был, но он уже не действителен.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "token"
                ],
                "summary": "Выдать новый токен пользователю",
                "parameters": [
                    {
                        "description": "Актуальные логин (почта) и пароль пользователя",
                        "name": "UserEmailNPassword",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает актуальный и новый токен для пользователя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные пользователя!",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Возвращает ошибку, если не удалось найти пользователя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса!",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/user": {
            "post": {
                "description": "Позволяет добавить нового пользователя в систему, если пользователя с такими данными не существовало!",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Создать нового пользователя",
                "parameters": [
                    {
                        "description": "Данные пользователя. А именно: Почта (email), пароль (password), name (имя), номер телефона (phoneNumber)",
                        "name": "UserData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает актуальный токен для пользователя, а также ID пользователя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные пользователя!",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "401": {
                        "description": "Возвращает ошибку, если не удалось добавить пользователя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса или не удалось создать для него токен! Конкретная ошибка будет в результате запроса!",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/user/otklik": {
            "post": {
                "description": "Создает отклик на вакансию при помощи ID пользователя и вакансии. Статус отклика автоматически присваевается \"Ожидание\"",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "vacancy"
                ],
                "summary": "Создание отклика на вакансию",
                "parameters": [
                    {
                        "description": "ID пользователя и вакансии, на которую нужно добавить отклик",
                        "name": "IDs",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает ID отклика. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!",
                        "schema": {
                            "type": "int"
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить request body. К ответу прикрепляется ID, который получил сервер, а также где именно произошла ошибка.",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/user/otkliks/{id}": {
            "get": {
                "description": "Возвращает список всех откликов для определенного пользователя по его ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Получение списка всех откликов для пользователя",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пользователя",
                        "name": "UID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает статус и массив откликов. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки!",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Возвращает ошибку, если не удалось преобразовать передаваемый параметр (ID) через URL.",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/vac": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/vac/:id": {
            "get": {
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
