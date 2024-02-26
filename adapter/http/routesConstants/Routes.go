package routesConstants

import (
	"net/http"
	"time"
)

const (
	GetBarRouteConst = "/"

	PostSignInRouteConst           = "/signin"
	PostSignOutRouteConst          = "/signout"
	PostSignUpRouteConst           = "/signup"

	PatchEditUserRouteConst            = "/user"
	GetLoggedUserRouteConst            = "/logged-user"
	GetUserRouteConst                  = "/user"
	GetUsersListRouteConst             = "/users"
)

const (
	BadRequestConst          = http.StatusBadRequest
	ForbiddenRequestConst    = http.StatusForbidden
	Unauthorized             = http.StatusUnauthorized
	StatusOk                 = http.StatusOK
	InternarServerErrorConst = http.StatusInternalServerError
	CreatedConst             = http.StatusCreated
	TimeoutConst             = http.StatusRequestTimeout
)

const (
	DataKeyConst               = "data"
)

const (
	DefaultLimitResponseTimeConst = 180 * time.Second
)
