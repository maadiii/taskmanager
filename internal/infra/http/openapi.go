package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func routeDocs(g *gin.Engine, port int) {
	g.GET("/openapi.json", func(c *gin.Context) {
		spec := strings.Replace(
			openAPISpec,
			"http://localhost:8080",
			fmt.Sprintf("http://localhost:%d", port),
			1,
		)
		c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(spec))
	})
	g.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUI))
	})
}

const swaggerUI = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Task Manager API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [SwaggerUIBundle.presets.apis],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`

const openAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Task Manager API",
    "description": "Task management API using cursor pagination, Redis cache-aside, and PostgreSQL.",
    "version": "1.0.0"
  },
  "servers": [
    {
      "url": "http://localhost:8080",
      "description": "Local server"
    }
  ],
  "tags": [
    {
      "name": "Tasks",
      "description": "Create and manage tasks"
    }
  ],
  "paths": {
    "/api/v1/tasks": {
      "post": {
        "tags": ["Tasks"],
        "summary": "Create a task",
        "operationId": "createTask",
        "description": "Creates a task for the authenticated user. The current demo authorization middleware supplies the user identity and create permission.",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/CreateTaskRequest"},
              "example": {"title": "Read documentation", "description": "Review the project README"}
            }
          }
        },
        "responses": {
          "201": {"description": "Task created", "headers": {"request-id": {"$ref": "#/components/headers/RequestID"}}, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/CreateTaskResponse"}}}},
          "400": {"$ref": "#/components/responses/BadRequest"},
          "403": {"$ref": "#/components/responses/Forbidden"},
          "409": {"$ref": "#/components/responses/Conflict"},
          "500": {"$ref": "#/components/responses/InternalServerError"}
        }
      },
      "get": {
        "tags": ["Tasks"],
        "summary": "List tasks",
        "operationId": "listTasks",
        "description": "Lists tasks using cursor pagination. Results are ordered by id descending.",
        "parameters": [
          {"$ref": "#/components/parameters/Status"},
          {"$ref": "#/components/parameters/Limit"},
          {"$ref": "#/components/parameters/LastID"}
        ],
        "responses": {
          "200": {"description": "Tasks returned", "headers": {"request-id": {"$ref": "#/components/headers/RequestID"}}, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ListTasksResponse"}}}},
          "400": {"$ref": "#/components/responses/BadRequest"},
          "500": {"$ref": "#/components/responses/InternalServerError"}
        }
      }
    },
    "/api/v1/tasks/{id}": {
      "parameters": [{"$ref": "#/components/parameters/TaskID"}],
      "get": {
        "tags": ["Tasks"],
        "summary": "Get a task",
        "operationId": "getTask",
        "responses": {
          "200": {"description": "Task returned", "headers": {"request-id": {"$ref": "#/components/headers/RequestID"}}, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/GetTaskResponse"}}}},
          "400": {"$ref": "#/components/responses/BadRequest"},
          "404": {"$ref": "#/components/responses/NotFound"},
          "500": {"$ref": "#/components/responses/InternalServerError"}
        }
      },
      "put": {
        "tags": ["Tasks"],
        "summary": "Update a task",
        "operationId": "updateTask",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {"$ref": "#/components/schemas/UpdateTaskRequest"},
              "example": {"title": "Updated documentation", "status": "DONE", "priority": "HIGH"}
            }
          }
        },
        "responses": {
          "200": {"description": "Task updated", "headers": {"request-id": {"$ref": "#/components/headers/RequestID"}}, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/UpdateTaskResponse"}}}},
          "400": {"$ref": "#/components/responses/BadRequest"},
          "403": {"$ref": "#/components/responses/Forbidden"},
          "404": {"$ref": "#/components/responses/NotFound"},
          "409": {"$ref": "#/components/responses/Conflict"},
          "500": {"$ref": "#/components/responses/InternalServerError"}
        }
      },
      "delete": {
        "tags": ["Tasks"],
        "summary": "Delete a task",
        "operationId": "deleteTask",
        "responses": {
          "200": {"description": "Task deleted", "headers": {"request-id": {"$ref": "#/components/headers/RequestID"}}, "content": {"application/json": {"schema": {"$ref": "#/components/schemas/DeleteTaskResponse"}}}},
          "400": {"$ref": "#/components/responses/BadRequest"},
          "403": {"$ref": "#/components/responses/Forbidden"},
          "404": {"$ref": "#/components/responses/NotFound"},
          "500": {"$ref": "#/components/responses/InternalServerError"}
        }
      }
    }
  },
  "components": {
    "parameters": {
      "TaskID": {"name": "id", "in": "path", "required": true, "description": "Task UUID.", "schema": {"type": "string", "format": "uuid"}},
      "Status": {"name": "status", "in": "query", "required": false, "description": "Filter by task status.", "schema": {"$ref": "#/components/schemas/TaskStatus"}},
      "Limit": {"name": "limit", "in": "query", "required": false, "description": "Page size. Defaults to the repository default; allowed range is 1-100.", "schema": {"type": "integer", "minimum": 1, "maximum": 100}},
      "LastID": {"name": "last_id", "in": "query", "required": false, "description": "Cursor: the last task id from the previous page.", "schema": {"type": "string", "format": "uuid"}}
    },
    "headers": {
      "RequestID": {"description": "Unique request identifier generated by the API.", "schema": {"type": "string"}}
    },
    "schemas": {
      "CreateTaskRequest": {"type": "object", "properties": {"title": {"type": "string", "description": "Task title."}, "description": {"type": "string", "description": "Task description."}}},
      "UpdateTaskRequest": {"type": "object", "description": "All fields are optional; omitted fields remain unchanged.", "properties": {"title": {"type": "string"}, "description": {"type": "string"}, "status": {"$ref": "#/components/schemas/TaskStatus"}, "priority": {"$ref": "#/components/schemas/TaskPriority"}}},
      "Task": {"type": "object", "required": ["id", "title", "description", "status", "priority", "createdAtUnixSec"], "properties": {"id": {"type": "string", "format": "uuid"}, "title": {"type": "string"}, "description": {"type": "string"}, "status": {"$ref": "#/components/schemas/TaskStatus"}, "priority": {"$ref": "#/components/schemas/TaskPriority"}, "createdAtUnixSec": {"type": "integer", "format": "int64"}}},
      "CreateTaskResponse": {"$ref": "#/components/schemas/Task"},
      "GetTaskResponse": {"$ref": "#/components/schemas/Task"},
      "UpdateTaskResponse": {"allOf": [{"$ref": "#/components/schemas/Task"}, {"type": "object", "required": ["updatedAtUnixSec"], "properties": {"updatedAtUnixSec": {"type": "integer", "format": "int64"}}}]},
      "DeleteTaskResponse": {"type": "object", "required": ["id", "deleted"], "properties": {"id": {"type": "string", "format": "uuid"}, "deleted": {"type": "boolean"}}},
      "ListTasksResponse": {"type": "object", "required": ["tasks"], "properties": {"tasks": {"type": "array", "items": {"$ref": "#/components/schemas/Task"}}}},
      "TaskStatus": {"type": "string", "enum": ["TODO", "IN_PROGRESS", "DONE"]},
      "TaskPriority": {"type": "string", "enum": ["LOW", "MEDIUM", "HIGH"]},
      "ErrorResponse": {"type": "object", "required": ["key"], "properties": {"key": {"type": "string", "description": "Stable machine-readable error key."}}, "examples": [{"key": "BAD_REQUEST"}, {"key": "FORBIDDEN"}, {"key": "TASK_ID_NOT_FOUND"}, {"key": "TASK_TITLE_ALREADY_EXISTS"}, {"key": "INTERNAL_SERVER_ERROR"}]}
    },
    "responses": {
      "BadRequest": {"description": "Invalid URI, query, header, or JSON body. Error key is BAD_REQUEST.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "example": {"key": "BAD_REQUEST"}}}},
      "Forbidden": {"description": "The authenticated user lacks the required permission. Error key is FORBIDDEN.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "example": {"key": "FORBIDDEN"}}}},
      "NotFound": {"description": "The task does not exist for the authenticated user. Error key is TASK_ID_NOT_FOUND.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "example": {"key": "TASK_ID_NOT_FOUND"}}}},
      "Conflict": {"description": "A task with the same title already exists for the user. Error key is TASK_TITLE_ALREADY_EXISTS.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "example": {"key": "TASK_TITLE_ALREADY_EXISTS"}}}},
      "InternalServerError": {"description": "Unexpected internal or persistence error. Error key is INTERNAL_SERVER_ERROR.", "content": {"application/json": {"schema": {"$ref": "#/components/schemas/ErrorResponse"}, "example": {"key": "INTERNAL_SERVER_ERROR"}}}}
    }
  }
}`
