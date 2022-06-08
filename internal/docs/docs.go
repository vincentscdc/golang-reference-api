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
        "contact": {
            "name": "Vincent",
            "email": "vincent.serpoul@crypto.com"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/internal/pay_later/refund": {
            "post": {
                "description": "refunds a payment",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment"
                ],
                "summary": "Refunds a payment",
                "parameters": [
                    {
                        "description": "Refund reqBody data",
                        "name": "refund_request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internalfacing.RefundRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "403": {
                        "description": "payment plan is unconfirmed",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "payment plan not found",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/internal/pay_later/user/{user_uuid}/credit_line": {
            "get": {
                "description": "returns the amount and status",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "credit_line"
                ],
                "summary": "Renders a user's credit line info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "user_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/internalfacing.CreditLineResponse"
                        }
                    }
                }
            }
        },
        "/api/internal/pay_later/user/{user_uuid}/payment_plans": {
            "post": {
                "description": "pre creates a payment plan",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment_plan"
                ],
                "summary": "Creates a pending a payment plan",
                "parameters": [
                    {
                        "description": "Pre create payment plan reqBody",
                        "name": "pre_create_payment_plan_request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/internalfacing.CreatePendingPaymentPlanRequest"
                        }
                    },
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "user_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "bad reqBody",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/internal/pay_later/user/{user_uuid}/payment_plans/{uuid}/cancel": {
            "post": {
                "description": "cancels a payment plan",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment_plan"
                ],
                "summary": "Cancels a payment plan",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "user_uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Payment Plan UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "bad reqBody",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "payment plan not belongs to user",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "payment plan is not in pending",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "payment plan not found",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/internal/pay_later/user/{user_uuid}/payment_plans/{uuid}/complete": {
            "post": {
                "description": "completes a payment plan",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment_plan"
                ],
                "summary": "Completes a payment plan",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "user_uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Payment Plan UUID",
                        "name": "uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "bad reqBody",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "payment plan not belongs to user",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "403": {
                        "description": "payment plan is not in pending",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "payment plan not found",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "internal error",
                        "schema": {
                            "$ref": "#/definitions/handlerwrap.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/pay_later/credit_line": {
            "get": {
                "description": "returns the amount and status",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "credit_line"
                ],
                "summary": "Renders a user's credit line info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "X-CRYPTO-USER-UUID",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userfacing.CreditLineResponseOKStyle"
                        }
                    }
                }
            }
        },
        "/api/pay_later/payment_plans": {
            "get": {
                "description": "returns pagination of one user's payment plans",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "payment_plan"
                ],
                "summary": "Renders a user's payment plans",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User UUID",
                        "name": "X-CRYPTO-USER-UUID",
                        "in": "header",
                        "required": true
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "description": "Start index in the list",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "maximum": 10,
                        "minimum": 0,
                        "type": "integer",
                        "description": "Number of items displayed",
                        "name": "limit",
                        "in": "query"
                    },
                    {
                        "enum": [
                            "asc",
                            "desc"
                        ],
                        "type": "string",
                        "default": "desc",
                        "description": "Order by payment.created_at asc  OR desc",
                        "name": "created_at_order",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/userfacing.PaymentPlansResponseOKStyle"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "handlerwrap.ErrorResponse": {
            "type": "object",
            "properties": {
                "error_code": {
                    "type": "string"
                },
                "error_msg": {
                    "type": "string"
                }
            }
        },
        "internalfacing.CreatePendingPaymentPlanRequest": {
            "type": "object",
            "properties": {
                "payment": {
                    "$ref": "#/definitions/internalfacing.PendingPayment"
                },
                "user_wallet_currency": {
                    "type": "string"
                }
            }
        },
        "internalfacing.CreditLineResponse": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "limit": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "internalfacing.PendingPayment": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "meta": {
                    "type": "object",
                    "properties": {
                        "amount": {
                            "type": "string"
                        },
                        "crypto_amounts": {
                            "type": "object",
                            "properties": {
                                "CRO": {
                                    "type": "string"
                                },
                                "USDC": {
                                    "type": "string"
                                }
                            }
                        },
                        "crypto_currency": {
                            "type": "string"
                        },
                        "currency": {
                            "type": "string"
                        },
                        "custom_id": {
                            "type": "string"
                        },
                        "deadline": {
                            "type": "string"
                        },
                        "is_approved": {
                            "type": "boolean"
                        },
                        "items": {
                            "type": "string"
                        },
                        "live_mode": {
                            "type": "boolean"
                        },
                        "merchant_reference": {
                            "type": "string"
                        },
                        "pay_later_amount": {
                            "type": "object",
                            "properties": {
                                "amount": {
                                    "type": "string"
                                },
                                "currency": {
                                    "type": "string"
                                }
                            }
                        },
                        "pay_later_installments": {
                            "type": "array",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "amount": {
                                        "type": "string"
                                    },
                                    "currency": {
                                        "type": "string"
                                    },
                                    "date": {
                                        "type": "string"
                                    }
                                }
                            }
                        },
                        "quotation_id": {
                            "type": "string"
                        },
                        "recipient": {
                            "type": "string"
                        },
                        "remaining_time": {
                            "type": "string"
                        },
                        "status": {
                            "type": "string"
                        },
                        "value_fiat": {
                            "type": "string"
                        }
                    }
                },
                "output": {
                    "type": "string"
                },
                "resource_id": {
                    "type": "string"
                },
                "resource_type": {
                    "type": "string"
                },
                "value": {
                    "type": "string"
                }
            }
        },
        "internalfacing.RefundRequest": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "payment_id": {
                    "type": "string"
                },
                "refund_data": {},
                "refund_id": {
                    "type": "string"
                }
            }
        },
        "userfacing.CreditLineResponse": {
            "type": "object",
            "properties": {
                "available_amount": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "total_amount": {
                    "type": "string"
                }
            }
        },
        "userfacing.CreditLineResponseOKStyle": {
            "type": "object",
            "properties": {
                "credit_info": {
                    "$ref": "#/definitions/userfacing.CreditLineResponse"
                },
                "error": {
                    "type": "string"
                },
                "error_message": {
                    "type": "string"
                },
                "ok": {
                    "type": "boolean"
                }
            }
        },
        "userfacing.PaymentPlanResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "from_currency": {
                    "type": "string"
                },
                "is_liquidated": {
                    "type": "boolean"
                },
                "next_repayment_id": {
                    "type": "string"
                },
                "outstanding_late_charge_amount": {
                    "type": "string"
                },
                "payable_amount": {
                    "type": "string"
                },
                "payments": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/userfacing.PaymentResponse"
                    }
                },
                "repayment_status_for_display": {
                    "type": "string"
                },
                "status": {
                    "type": "string"
                },
                "total_amount": {
                    "type": "string"
                },
                "total_late_charge_amount": {
                    "type": "string"
                },
                "total_paid_amount": {
                    "type": "string"
                },
                "total_refund_amount": {
                    "type": "string"
                }
            }
        },
        "userfacing.PaymentPlansResponseOKStyle": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string"
                },
                "error_message": {
                    "type": "string"
                },
                "ok": {
                    "type": "boolean"
                },
                "payment_plans": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/userfacing.PaymentPlanResponse"
                    }
                }
            }
        },
        "userfacing.PaymentResponse": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "due_at": {
                    "type": "string"
                },
                "from_amount": {
                    "type": "string"
                },
                "from_currency": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "late_charge_amount": {
                    "type": "string"
                },
                "outstanding_amount": {
                    "type": "string"
                },
                "refund_amount": {
                    "type": "string"
                },
                "settled_at": {
                    "type": "string"
                },
                "status": {
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
	Title:            "bnpl API",
	Description:      "This is the service used for buy now pay later.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
