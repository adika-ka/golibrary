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
        "/authors": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authors"
                ],
                "summary": "Получить список всех авторов",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Author"
                            }
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authors"
                ],
                "summary": "Добавить нового автора",
                "parameters": [
                    {
                        "description": "Имя автора",
                        "name": "author",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/entities.Author"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entities.Author"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
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
        "/authors/top": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "authors"
                ],
                "summary": "Получить топ авторов по количеству выданных книг",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Author"
                            }
                        }
                    }
                }
            }
        },
        "/books": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Получить список всех книг",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Book"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Добавляет книгу с указанным названием и ID автора",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "books"
                ],
                "summary": "Добавить новую книгу",
                "parameters": [
                    {
                        "description": "Данные книги",
                        "name": "book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.BookCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entities.Book"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
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
        "/loans/borrow": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "loans"
                ],
                "summary": "Выдать книгу пользователю (создать запись аренды)",
                "parameters": [
                    {
                        "description": "Данные аренды (user_id, book_id)",
                        "name": "loan",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.CreateLoanRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Книга успешно выдана"
                    },
                    "400": {
                        "description": "Неверный запрос или книга уже выдана",
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
        "/loans/return": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "loans"
                ],
                "summary": "Принять книгу обратно в библиотеку (закрыть активный заем)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пользователя",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "ID книги",
                        "name": "book_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Книга успешно возвращена"
                    },
                    "400": {
                        "description": "Неверный запрос или книга не числится в аренде",
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
        "/users": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Получить список всех пользователей",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.User"
                            }
                        }
                    }
                }
            }
        },
        "/users/borrowed-books": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "users"
                ],
                "summary": "Получить список книг, взятых пользователем в аренду (активные займы)",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID пользователя",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список книг, которые пользователь не вернул",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.Book"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.BookCreateRequest": {
            "type": "object",
            "properties": {
                "author_id": {
                    "type": "integer",
                    "example": 52
                },
                "title": {
                    "type": "string",
                    "example": "Some Book Title"
                }
            }
        },
        "controller.CreateLoanRequest": {
            "type": "object",
            "properties": {
                "book_id": {
                    "type": "integer"
                },
                "user_id": {
                    "type": "integer"
                }
            }
        },
        "entities.Author": {
            "type": "object",
            "properties": {
                "books": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Book"
                    }
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "entities.Book": {
            "type": "object",
            "properties": {
                "author": {
                    "$ref": "#/definitions/entities.Author"
                },
                "author_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "entities.User": {
            "type": "object",
            "properties": {
                "Email": {
                    "type": "string"
                },
                "ID": {
                    "type": "integer"
                },
                "Name": {
                    "type": "string"
                },
                "RentedBooks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entities.Book"
                    }
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Library API",
	Description:      "API для управления библиотекой.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
