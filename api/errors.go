package api

import "errors"

var (
	ErrBadMethod error = errors.New("api.RegisterApiEndpoint(): bad method")

	ErrInvalidCategory error = errors.New("api: invalid category is selected")

	ErrNotAllowDirectFuncReg error = errors.New("api.RegisterApiEndpoint(): no direct handler func registration is allowed")

	ErrRepeatGetPath  error = errors.New("api.RegisterApiEndpoint(): repeated path for GET method")
	ErrRepeatPostPath error = errors.New("api.RegisterApiEndpoint(): repeated path for POST method")
)
