package response

import "errors"

// Server errors.
var (
	ErrTimeout               = errors.New("timeout exceeded")
	ErrorNotFound            = errors.New("Not Found")
	ErrorBadRequest          = errors.New("Bad resquest")
	ErrorBadEmailOrPassword  = errors.New("Bad email or password")
	ErrorUnauthorized        = errors.New("Unauthorized")
	ErrorInternalServerError = errors.New("Interal server error")
	ErrorUUIDNotFound        = errors.New("uid not found")
	ErrorParsingUser         = errors.New("failed to parse user")
	ErrInvalidID             = errors.New("invalid id")
	ErrCouldNotInsert        = errors.New("Error could not insert")
	ErrInvalidEmail          = errors.New("Error invalid email")
	ErrCantFollowYou         = errors.New("Error you can't follow you")
)
