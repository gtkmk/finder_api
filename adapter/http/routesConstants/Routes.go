package routesConstants

import (
	"net/http"
	"time"
)

const (
	GetBarRouteConst = "/"

	PostSignInRouteConst  = "/signin"
	PostSignOutRouteConst = "/signout"
	PostSignUpRouteConst  = "/signup"

	PatchEditUserRouteConst       = "/user"
	PatchFirstAccessRouteConst    = "/user/first-access"
	PatchForgotPasswordRouteConst = "/user/forgot-password"
	PatchResetPasswordRouteConst  = "/user/reset-password"
	PostForgotPasswordRouteConst  = "/user/forgot-password"
	GetLoggedUserRouteConst       = "/logged-user"
	GetUserDetailsRouteConst      = "/user"
	GetUsersListByNameRouteConst  = "/users"
	PatchUserInfoRouteConst       = "/user/info"

	PostCreatePostRouteConst   = "/post"
	PostFindAllPostsRouteConst = "/posts"
	PostEditPostRouteConst     = "/post"
	DeletePostRouteConst       = "/post"
	FindPostParamsRouteConst   = "/post-params"
	PostAnimalFoundRouteConst  = "/post/animal-found"

	PostLikeRouteConst = "/like"

	PostCreateCommentRouteConst   = "/comment"
	PatchEditCommentRouteConst    = "/comment"
	DeleteCommentRouteConst       = "/comment"
	FindFindAllCommentsRouteConst = "/comments"

	FindDocumentBase64RouteConst       = "/document/imageBase64"
	UpdateChangeProfileImageRouteConst = "/document/changeProfileImage"

	PostFollowRouteConst = "/follow"

	// === Route marker ===
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
	DataKeyConst    = "data"
	MessageKeyConst = "message"
)

const (
	DefaultLimitResponseTimeConst = 180 * time.Second
)
