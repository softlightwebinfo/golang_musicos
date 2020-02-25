package libs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HelloWorldResponse
//
// This is the structure used to respond with simple messages
//
// swagger:model HelloWorldResponse
type HelloWorldResponse struct {
	Message      string `json:"message" xml:"message`
	ProvidedName string `json:"name" xml:"name`
}
type ErrorResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

func JSONResponse(w http.ResponseWriter, code int, output interface{}) {
	// Convert our interface to JSON
	response, _ := json.Marshal(output)
	// Set the content type to json for browsers
	w.Header().Set("Content-Type", "application/json")
	// Our response code
	w.WriteHeader(code)

	_, _ = w.Write(response)
}
func JSONResponseOk(c *gin.Context, output interface{}) {
	c.JSON(200, output)
}
func JSONResponseError(c *gin.Context, message string, error interface{}) {
	e := ErrorResponse{
		Message: message,
		Error:   error,
	}
	c.JSON(400, e)
}
func JSONResponseNotFound(c *gin.Context, message string, error interface{}) {
	e := ErrorResponse{
		Message: message,
		Error:   error,
	}
	c.JSON(404, e)
}
func JSONResponseUnauthorized(c *gin.Context, message string, error interface{}) {
	e := ErrorResponse{
		Message: message,
		Error:   error,
	}
	c.JSON(401, e)
}
