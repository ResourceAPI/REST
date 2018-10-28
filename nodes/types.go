package nodes

import "github.com/labstack/echo"

type Response struct {
	Success bool         `json:"success"`
	Data    *interface{} `json:"data,omitempty"`
	Error   *Error       `json:"error,omitempty"`
}

type ResponseResource struct {
	Success bool                      `json:"success"`
	Data    *[]map[string]interface{} `json:"data"`
	Error   *Error                    `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

var (
	ErrorCouldNotReadBody     = Error{Code: 1, Message: "could not read body of request: ", Status: 400}
	ErrorResourceDoesNotExist = Error{Code: 2, Message: "resource does not exist", Status: 404}
	ErrorResourceInvalid      = Error{Code: 3, Message: "resource does not meet schema: ", Status: 400}
	ErrorFilterInvalid        = Error{Code: 4, Message: "filter is not valid: ", Status: 400}
	ErrorFilterDoesntExist    = Error{Code: 5, Message: "filter type does not exist", Status: 400}
	ErrorFetchingResource     = Error{Code: 6, Message: "error fetching resource: ", Status: 400}
	ErrorDeletingResource     = Error{Code: 7, Message: "error deleting resource: ", Status: 400}
	ErrorCreatingResource     = Error{Code: 8, Message: "error creating resource: ", Status: 400}
	ErrorUpdatingResource     = Error{Code: 9, Message: "error updating resource: ", Status: 400}
)

func PrepError(c echo.Context, err Error) error {
	return c.JSON(err.Status, Response{
		Success: false,
		Error:   &err,
	})
}
