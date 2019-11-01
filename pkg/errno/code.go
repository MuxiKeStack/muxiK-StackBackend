package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrDatabase   = &Errno{Code: 20002, Message: "Database error."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}
	ErrDelete     = &Errno{Code: 20004, Message: "Do not have this delete target"}

	// Auth errors
	ErrAuthFailed   = &Errno{Code: 20101, Message: "The sid or password was incorrect."}
	ErrTokenInvalid = &Errno{Code: 20102, Message: "The token was invalid."}

	// user errors
	ErrCreateUser   = &Errno{Code: 20201, Message: "Error in creating user."}
	ErrUpdateUser   = &Errno{Code: 20202, Message: "Error in updating user"}
	ErrUserNotFound = &Errno{Code: 20203, Message: "The user was not found."}
	ErrGetUserInfo  = &Errno{Code: 20204, Message: "Error in getting user info"}

	// table errors
	ErrTableExisting = &Errno{Code: 20401, Message: "The table is not existing "}
	ErrClassExisting = &Errno{Code: 20402, Message: "The class is not existing "}
	ErrClassIdRequired = &Errno{Code: 20403, Message: "The classId is required "}
	ErrGetClassInfo = &Errno{Code: 20404, Message: "Error occurred in getting class info"}
)
