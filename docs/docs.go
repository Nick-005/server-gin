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
        "/auth/user": {
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
                            "$ref": "#/definitions/main.RequestNewToken"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает актуальный и новый токен для пользователя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!",
                        "schema": {
                            "$ref": "#/definitions/main.TokenForUser"
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные пользователя!",
                        "schema": {
                            "$ref": "#/definitions/main.InfoError"
                        }
                    },
                    "401": {
                        "description": "Возвращает ошибку, если не удалось найти пользователя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса!",
                        "schema": {
                            "$ref": "#/definitions/main.SimpleError"
                        }
                    }
                }
            }
        },
        "/emp": {
            "post": {
                "description": "Позволяет создать работодателя в системе. Будет возвращен ID и токен для работодателя!",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "employer"
                ],
                "summary": "Создать работодателя",
                "parameters": [
                    {
                        "description": "Данные работодателя",
                        "name": "EmpData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.RequestEmployee"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает ID (И попозже будет Token) работодателя.",
                        "schema": {
                            "$ref": "#/definitions/main.NewEmployer"
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить body-request!",
                        "schema": {
                            "$ref": "#/definitions/main.InfoError"
                        }
                    },
                    "401": {
                        "description": "Возвращает ошибку, если не добавить работодателя с корректными данными. Конкретная ошибка будет в результате запроса!",
                        "schema": {
                            "$ref": "#/definitions/main.SimpleError"
                        }
                    }
                }
            }
        },
        "/emp/:id": {
            "get": {
                "description": "Позволяет получить данные работодателя по его ID. Будет возвращен ID вакансии!",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "vacancy"
                ],
                "summary": "Получить данные работодателя по его ID",
                "parameters": [
                    {
                        "description": "Данные вакансии",
                        "name": "VacData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Vacancy_Body"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает ID вакансии.",
                        "schema": {
                            "$ref": "#/definitions/main.NewVacancy"
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить body-request!",
                        "schema": {
                            "$ref": "#/definitions/main.InfoError"
                        }
                    },
                    "401": {
                        "description": "Возвращает ошибку, если не удалось добавить вакансию с переданными данными. Конкретная ошибка будет в результате запроса!",
                        "schema": {
                            "$ref": "#/definitions/main.SimpleError"
                        }
                    }
                }
            }
        },
        "/emp/vacs/id": {
            "get": {
                "description": "Позволяет получить массив данных о всех вакансиях, которые есть у работодателя. Для этого нужно передать ID работодателя!",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "employer"
                ],
                "summary": "Получить все вакансии работодателя",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID работодателя",
                        "name": "EmpID",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает массив актуальных вакансий от одного работодателя.",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/sqlite.ResponseVac"
                            }
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить ID",
                        "schema": {
                            "$ref": "#/definitions/main.InfoError"
                        }
                    },
                    "401": {
                        "description": "Возвращает ошибку, если не удалось получить список всех вакансий! Конкретная ошибка будет в результате запроса!",
                        "schema": {
                            "$ref": "#/definitions/main.SimpleError"
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
                            "$ref": "#/definitions/main.RequestAdd"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает актуальный токен для пользователя, а также ID пользователя. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!",
                        "schema": {
                            "$ref": "#/definitions/main.AddNewUser"
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить body, который отвечает за данные пользователя!",
                        "schema": {
                            "$ref": "#/definitions/main.InfoError"
                        }
                    },
                    "401": {
                        "description": "Возвращает ошибку, если не удалось добавить пользователя в БД, который соответствовал бы данным, которые были получены сервером в результате этого запроса или не удалось создать для него токен! Конкретная ошибка будет в результате запроса!",
                        "schema": {
                            "$ref": "#/definitions/main.SimpleError"
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
                            "$ref": "#/definitions/main.RequestResponse"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает ID отклика. Если произошла ошибка - статус будет 'Err' и будет возвращен текст ошибки! Также будет известно, где именно произошла ошибка!",
                        "schema": {
                            "type": "integer"
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить request body. К ответу прикрепляется ID, который получил сервер, а также где именно произошла ошибка.",
                        "schema": {
                            "$ref": "#/definitions/main.InfoError"
                        }
                    }
                }
            }
        },
        "/user/otkliks/{id}": {
            "get": {
                "description": "Возвращает список всех откликов для определенного пользователя по его ID",
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
                            "$ref": "#/definitions/main.AllUserResponseOK"
                        }
                    },
                    "404": {
                        "description": "Возвращает ошибку, если не удалось преобразовать передаваемый параметр (ID) через URL.",
                        "schema": {
                            "$ref": "#/definitions/main.SimpleError"
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
                "description": "Позволяет создать новую вакансию в системе. Будет возвращен ID вакансии!",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "vacancy"
                ],
                "summary": "Создать вакансию",
                "parameters": [
                    {
                        "description": "Данные вакансии",
                        "name": "VacData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.Vacancy_Body"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Возвращает ID вакансии.",
                        "schema": {
                            "$ref": "#/definitions/main.NewVacancy"
                        }
                    },
                    "400": {
                        "description": "Возвращает ошибку, если не удалось распарсить body-request!",
                        "schema": {
                            "$ref": "#/definitions/main.InfoError"
                        }
                    },
                    "401": {
                        "description": "Возвращает ошибку, если не удалось добавить вакансию с переданными данными. Конкретная ошибка будет в результате запроса!",
                        "schema": {
                            "$ref": "#/definitions/main.SimpleError"
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
    },
    "definitions": {
        "main.AddNewUser": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "uid": {
                    "type": "integer"
                }
            }
        },
        "main.AllUserResponseOK": {
            "type": "object",
            "properties": {
                "otkliks": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "main.InfoError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "info": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "main.NewEmployer": {
            "type": "object",
            "properties": {
                "emp_id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "main.NewVacancy": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "vacancyID": {
                    "type": "integer"
                }
            }
        },
        "main.RequestAdd": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                }
            }
        },
        "main.RequestEmployee": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "inn": {
                    "type": "string"
                },
                "nameOrg": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                }
            }
        },
        "main.RequestNewToken": {
            "type": "object",
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "main.RequestResponse": {
            "type": "object",
            "properties": {
                "UID": {
                    "type": "integer"
                },
                "vac_id": {
                    "type": "integer"
                }
            }
        },
        "main.SimpleError": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "main.TokenForUser": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "main.Vacancy_Body": {
            "type": "object",
            "properties": {
                "about": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "emp_id": {
                    "type": "integer"
                },
                "exp": {
                    "type": "integer"
                },
                "is_visible": {
                    "type": "boolean"
                },
                "location": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "vac_name": {
                    "type": "string"
                }
            }
        },
        "sqlite.ResponseVac": {
            "type": "object",
            "properties": {
                "ID": {
                    "type": "integer"
                },
                "about": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "emp_id": {
                    "type": "integer"
                },
                "exp": {
                    "type": "string"
                },
                "is_visible": {
                    "type": "boolean"
                },
                "location": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "vac_name": {
                    "type": "string"
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
