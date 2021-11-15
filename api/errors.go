package api

import "errors"

var (
	ErrBadMethod error = errors.New("api: bad method")

	ErrInvalidCategory error = errors.New("api: invalid category is selected")

	ErrNotAllowDirectFuncReg error = errors.New("api: no direct handler func registration is allowed")

	ErrRepeatGetPath  error = errors.New("api: repeated path for GET method")
	ErrRepeatPostPath error = errors.New("api: repeated path for POST method")

	ErrUnknownUserGroup error = errors.New("api: usergroup has no known access control function")

	ErrInvalidUserGroup          error = errors.New("api: invalid usergroup")
	ErrAccessControlFuncNotFound error = errors.New("api: access control function not found")
)
