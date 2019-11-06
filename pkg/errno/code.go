package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase   = &Errno{Code: 20002, Message: "Database error."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}
	ErrGetQuery   = &Errno{Code: 20004, Message: "Error occurred while getting query. "}
	ErrGetParam   = &Errno{Code: 20005, Message: "Error occurred while getting path params. "}
	ErrDelete     = &Errno{Code: 20006, Message: "Error occurred while deleting sth. "}

	// Auth errors
	ErrAuthFailed   = &Errno{Code: 20101, Message: "The sid or password was incorrect."}
	ErrTokenInvalid = &Errno{Code: 20102, Message: "The token was invalid."}

	// user errors
	ErrCreateUser   = &Errno{Code: 20201, Message: "Error occurred in creating user."}
	ErrUpdateUser   = &Errno{Code: 20202, Message: "Error occurred in updating user"}
	ErrUserNotFound = &Errno{Code: 20203, Message: "The user was not found."}
	ErrGetUserInfo  = &Errno{Code: 20204, Message: "Error in getting user info"}

	// comment errors
	ErrNotLiked       = &Errno{Code: 20301, Message: "User has not liked yet. "}
	ErrEvaluationList = &Errno{Code: 20302, Message: "Error occurred while getting evaluation list. "}
	ErrCommentList    = &Errno{Code: 20303, Message: "Error occurred while getting comment list. "}
	//ErrGetEvaluationInfo = &Errno{Code: 20304, Message: "Error occurred while getting evaluation info. "}

	// table errors
	ErrTableExisting = &Errno{Code: 20401, Message: "The table is not existing. "}
	ErrClassExisting = &Errno{Code: 20402, Message: "The class is not existing. "}
	ErrGetTableInfo  = &Errno{Code: 20403, Message: "Error occurred in getting table info. "}
	ErrGetClassInfo  = &Errno{Code: 20404, Message: "Error occurred in getting class info."}

	// message errors
	ErrGetMessage = &Errno{Code: 20501, Message: "Error occurred in getting message list"}
)
