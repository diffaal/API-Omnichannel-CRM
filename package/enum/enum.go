package enum

const (
	UNCLAIMED   = "UNCLAIMED"
	WAITING     = "WAITING"
	IN_PROGRESS = "IN_PROGRESS"
	ACTIVE      = "ACTIVE"
	INACTIVE    = "INACTIVE"
	CLOSED      = "CLOSED"
	MISSED      = "MISSED"
	UNPROCESSED = "UNPROCESSED"
	PROCESSED   = "PROCESSED"
)

const (
	FACEBOOK  = "FACEBOOK"
	IG        = "INSTAGRAM"
	WA        = "WHATSAPP"
	EMAIL     = "EMAIL"
	LIVE_CHAT = "LIVE_CHAT"
)

// websocket platform
const (
	WEBHOOK     = "WEBHOOK"
	OMNICHANNEL = "OMNICHANNEL"
)

const (
	AGENT    = "AGENT"
	REPORTER = "REPORTER"
)

const (
	PESAN   = "PESAN"
	MENTION = "MENTION"
)

const (
	IMAGE    = "IMAGE"
	VIDEO    = "VIDEO"
	LOCATION = "LOCATION"
)

const (
	ROLE_ADMIN_PUSAT         int = 1
	ROLE_AGENT_PUSAT         int = 2
	ROLE_ADMIN_PROVINSI      int = 3
	ROLE_DISPATCHER_PROVINSI int = 4
	ROLE_RESPONDER_PROVINSI  int = 5
	ROLE_ADMIN_KOTA          int = 6
	ROLE_DISPATCHER_KOTA     int = 7
	ROLE_RESPONDER_KOTA      int = 8
)

const (
	CATEGORY_TYPE           int = 1
	CLASSIFICATION_TYPE     int = 2
	SUBCLASSIFICATION1_TYPE int = 3
	SUBCLASSIFICATION2_TYPE int = 4
	SUBCLASSIFICATION3_TYPE int = 5
)

const (
	PROBING1_TYPE int = 1
	PROBING2_TYPE int = 2
	PROBING3_TYPE int = 3
)

const (
	// Auth Response Enum
	INVALID_EMAIL_STATUS       = "INVALID_EMAIL"
	PASSWORD_REQUIRED_STATUS   = "PASSWORD_REQUIRED"
	EMAIL_REQUIRED_STATUS      = "EMAIL_REQUIRED"
	USERNAME_REQUIRED_STATUS   = "USERNAME_REQUIRED"
	FIRSTNAME_REQUIRED_STATUS  = "FIRSTNAME_REQUIRED"
	SUPERVISOR_REQUIRED_STATUS = "SUPERVISOR_REQUIRED"
	ACCOUNT_NOT_FOUND_STATUS   = "ACCOUNT_NOT_FOUND"
	CREDENTIALS_WRONG_STATUS   = "CREDENTIALS_WRONG"
	ACCOUNT_EXISTED_STATUS     = "ACCOUNT_EXSISTS"
	INVALID_PASSWORD_STATUS    = "INVALID_PASSWORD"
	FIELD_REQUIRED_STATUS      = "FIELD_REQUIRED"
	INVALID_QUERY_STATUS       = "INVALID QUERY"

	INVALID_EMAIL_MESSAGE       = "Please Provide a valid email"
	PASSWORD_REQUIRED_MESSAGE   = "Password field must not be empty"
	EMAIL_REQUIRED_MESSAGE      = "Email field must not be empty"
	USERNAME_REQUIRED_MESSAGE   = "Username Field must filled"
	FIRSTNAME_REQUIRED_MESSAGE  = "First name field must fileld"
	SUPERVISOR_REQUIRED_MESSAGE = "Supervisor field must filled"
	ACCOUNT_NOT_FOUND_MESSAGE   = "Account not found, Please check the email and password"
	CREDENTIALS_WRONG_MESSAGE   = "Your Credentials is wrong, please filled with the correct data"
	ACCOUNT_EXISTED_MESSAGE     = "The account with the email already exists, please use another email"
	INVALID_PASSWORD_MESSAGE    = "Invalid password, Please re-enter the correct password"
	FIELD_REQUIRED_MESSAGE      = "Field must be filled"
	INVALID_QUERY_MESSAGE       = "Invalid query params, Please re-enter with the right values"

	// System Error Response
	SYSTEM_BUSY_STATUS  = "SYSTEM_BUSY"
	SYSTEM_BUSY_MESSAGE = "System Busy"

	FAILED_BIND_JSON_STATUS  = "FAILED_BIND_JSON"
	FAILED_BIND_JSON_MESSAGE = "Invalid JSON Body"

	FAILED_REQUEST_SERVICE_MESSAGE = "Failed to request API to another service"

	// Invalid Query Parameter Response
	INVALID_PARAMETER_STATUS  = "INVALID_PARAMETER"
	INVALID_PARAMETER_MESSAGE = "The parameter in the request is invalid"

	// Unauthorized Response
	UNAUTHORIZED_STATUS  = "UNAUTHORIZED"
	UNAUTHORIZED_MESSAGE = "The Access Token is Expired or Not Provided"

	FORBIDDEN_STATUS  = "FORBIDDEN"
	FORBIDDEN_MESSAGE = "You do not have permission to access the requested resource"

	DATA_NOT_FOUND_STATUS  = "DATA_NOT_FOUND"
	DATA_NOT_FOUND_MESSAGE = "The requested data does not exist"

	// Ticket Response Enum
	TITLE_REQUIRED_STATUS        = "TITLE_REQUIRED"
	REPORTER_REQUIRED_STATUS     = "REPORTER_NAME_REQUIRED"
	PHONE_REQUIRED_STATUS        = "PHONE_REQUIRED"
	CATEGORY_REQUIRED_STATUS     = "CATEGORY_REQUIRED"
	PROBLEM_TYPE_REQUIRED_STATUS = "PROBLEM_TYPE_REQUIRED"
	SEVERITY_REQUIRED_STATUS     = "SEVERITY_REQUIRED"
	DISPATCH_TO_REQUIRED_STATUS  = "DISPATCH_REQUIRED"
	USER_REQUIRED_STATUS         = "USER_REQUIRED"

	TITLE_REQUIRED_MESSAGE        = "Ticket title field must be filled"
	REPORTER_REQUIRED_MESSAGE     = "Reporter Name field must be filled"
	PHONE_REQUIRED_MESSAGE        = "Phone Number field must be filled"
	CATEGORY_REQUIRED_MESSAGE     = "Category field must be filled"
	PROBLEM_TYPE_REQUIRED_MESSAGE = "Problem Type field must be filled"
	SEVERITY_REQUIRED_MESSAGE     = "Severity field must be filled"
	DISPATCH_TO_REQUIRED_MESSAGE  = "Dispatch destination field must be filled"
	USER_REQUIRED_MESSAGE         = "User id field must be filled"

	// Severity Response Enum
	NAME_REQUIRED_STATUS  = "NAME_REQUIRED"
	NAME_REQUIRED_MESSAGE = "Name field must be filled"

	FAILED_STATUS                        = "FAILED"
	USER_DO_NOT_HAVE_CHANNEL_ACCOUNT_MSG = "USER_DO_NOT_HAVE_CHANNEL_ACCOUNT"
	PLATFORM_ID_NOT_SET_MSG              = "PLATFORM_ID_NOT_SET"
	PLATFORM_ACCESS_TOKEN_NOT_SET_MSG    = "PLATFORM_ACCESS_TOKEN_NOT_SET"
	CHANNEL_ACCOUNT_NOT_MATCH_MSG        = "CHANNEL_ACCOUNT_NOT_MATCH"

	INVALID_PLATFORM_MSG = "INVALID PLATFORM"
)
