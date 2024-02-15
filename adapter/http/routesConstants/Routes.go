package routesConstants

import (
	"net/http"
	"time"
)

const (
	GetBarRouteConst = "/"

	PostCompanyRoutesConst  = "/company"
	GetCompaniesRoutesConst = "/companies"
	GetCompanyRoutesConst   = "/company"

	PostEnterpriseRouteConst  = "/enterprise"
	GetEnterpriseRouteConst   = "/enterprise"
	GetEnterprisesRouteConst  = "/enterprises"
	PatchEnterpriseRouteConst = "/enterprise"

	PostPermissionGroupRouteConst     = "/permission-group"
	PostPermissionMappingRouteConst   = "/permission-mapping"
	PostPermissionRouteConst          = "/permission"
	GetAllPermissionMappingRouteConst = "/permissions-mappings"
	GetPermissionMappingRouteConst    = "/permissions-mapping"
	GetPermissionGroupRouteConst      = "/permissions-groups"
	GetPermissionRouteConst           = "/permissions"

	PostProductRouteConst       = "/product"
	PostProductTypeRouteConst   = "/product-type"
	GetAllProductTypeRouteConst = "/product-types"
	GetProductTypeRouteConst    = "/product-type"

	PostPropertyRouteConst  = "/property"
	PatchPropertyRouteConst = "/property"

	GetBestPaymentDayRouteConst         = "/payment-day"
	PostProposalDocumentsRouteConst     = "/proposal-documents"
	PostProposalPersonRouteConst        = "/proposal-person"
	PostProposalRouteConst              = "/proposal"
	PostNewProposalRouteConst           = "/new-proposal"
	PostNewProposalDifferenceRouteConst = "/new-proposal/difference"
	PutAcceptNewProposalRouteConst      = "/accept-proposal"
	PutNewProposalRouteConst            = "/new-proposal"
	PutNewProposalDifferenceRouteConst  = "/new-proposal/difference"
	DeleteNewProposalRouteConst         = "/reject-proposal"
	GetWebsocketTicketRouteConst        = "/ws/ticket"
	GetFindProposalPersonRouteConst     = "/proposal-person"
	GetFindProposalInstallmentsConst    = "/proposal-installments"
	GetFindProposalRouteConst           = "/proposal"
	GetListProposalsRouteConst          = "/proposals"
	PatchPendingDocumentsRouteConst     = "/pending-document"
	PatchProposalClientRouteConst       = "/proposal/client"
	PatchProposalGuarantorRouteConst    = "/proposal/guarantor"
	PatchRefuseGuarantorConst           = "/proposal/refuse-guarantor"
	PatchRefuseGuaranteeConst           = "/proposal/refuse-guarantee"
	PatchProposalMessagesRouteConst     = "/proposal/messages"
	PatchProposalRisksRouteConst        = "/proposal/risks"
	PatchProposalApprovedAndFormalized  = "/proposal/approved-and-formalized"
	PatchProposalResponsibleRouteConst  = "/proposal/responsible"
	PatchProposalRouteConst             = "/proposal"
	PatchProposalVehicleRouterConst     = "/proposal/vehicle"
	PostSaveScoreRoutesConst            = "/save-score"
	GetProductCategoriesRouteConst      = "/product-categories"
	PostCreateDifferenceProposalConst   = "/proposal/difference"
	GetProposalsByPersonRouteConst      = "/proposals/client"
	GetProposalHistoryRouteConst        = "/proposal/history"
	PatchProposalEmcashNotEnabledConst  = "/proposal/emcash-unabled"
	PatchProposalRetrogressRouteConst   = "/proposal/retrogress"
	GetNotFormalizedProposalsConst      = "/wallet/not-formalized-proposals"
	PatchUpdateProposalAnalyst          = "/proposal/analyst"
	PostObservationRouteConst           = "/observation"
	GetObservationsRouteConst           = "/observations"
	PatchProposalEmcashEnabledConst     = "/proposal/emcash-enabled"

	PostSignInRouteConst           = "/signin"
	PostSignOutRouteConst          = "/signout"
	PostSignUpRouteConst           = "/signup"
	PostMassRegistrationRouteConst = "/mass-registration"

	PostUserProductRouteConst = "/user/product"

	PatchEditUserRouteConst            = "/user"
	GetLoggedUserRouteConst            = "/logged-user"
	GetManagersListRouteConst          = "/managers"
	GetUserRouteConst                  = "/user"
	GetUsersListRouteConst             = "/users"
	GetManagerAndConsultantsListConst  = "/users/managers-and-consultants"
	GetUserUnitiesRouteConst           = "/user/unities"
	PatchFirstAccessRouteConst         = "/user/first-access"
	PatchForgotPasswordRouteConst      = "/user/forgot-password"
	PatchResetPasswordRouteConst       = "/user/reset-password"
	PatchUserPermissionGroupRouteConst = "/user/permission-group"
	PostForgotPasswordRouteConst       = "/user/forgot-password"
	GetExportUsersRouteConst           = "/users/export"
	GetUserLeaderOptionsConst          = "/user/leader-options"
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
	MessageKeyConst            = "message"
	ProposalKeyConst           = "proposal"
	InstallmentsKeyConst       = "installments"
	ListKeyConst               = "list"
	TicketKeyConst             = "ticket"
	PermissionMappingKeyConst  = "permission_mapping"
	PermissionMappingsKeyConst = "permission_mappings"
	ProductTypeKeyConst        = "product_type"
	EnterprisesKeyConst        = "enterprises"
	EnterpriseKeyConst         = "enterprise"
	SellerKeyConst             = "seller"
	HistoryKeyConst            = "history"
)

const (
	UnityIdConst     = "unity_id"
	UnityUserIdConst = "unity_user_id"
)

const (
	DefaultLimitResponseTimeConst = 180 * time.Second
)
