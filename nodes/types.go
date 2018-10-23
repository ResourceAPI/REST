package nodes

import (
	"github.com/Vilsol/GoLib"
)

var (
	ErrorCouldNotReadBody     = GoLib.ErrorResponse{Code: 1, Message: "could not read body of request", Status: 400}
	ErrorResourceDoesNotExist = GoLib.ErrorResponse{Code: 2, Message: "resource does not exist", Status: 404}
	ErrorResourceInvalid      = GoLib.ErrorResponse{Code: 3, Message: "resource does not meet schema: ", Status: 400}
	ErrorFilterInvalid        = GoLib.ErrorResponse{Code: 4, Message: "filter is not valid: ", Status: 400}
	ErrorFilterDoesntExist    = GoLib.ErrorResponse{Code: 5, Message: "filter type does not exist", Status: 400}
	ErrorFetchingResource     = GoLib.ErrorResponse{Code: 6, Message: "error fetching resource: ", Status: 400}
	ErrorDeletingResource     = GoLib.ErrorResponse{Code: 7, Message: "error deleting resource: ", Status: 400}
	ErrorCreatingResource     = GoLib.ErrorResponse{Code: 8, Message: "error creating resource: ", Status: 400}
	ErrorUpdatingResource     = GoLib.ErrorResponse{Code: 9, Message: "error updating resource: ", Status: 400}
)
