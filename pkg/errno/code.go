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

	ErrSecurityCheck = &Errno{Code: 20007, Message: "Error occurred in security check"}

	// Auth errors
	ErrAuthFailed   = &Errno{Code: 20101, Message: "The sid or password was incorrect."}
	ErrTokenInvalid = &Errno{Code: 20102, Message: "The token was invalid."}

	// user errors
	ErrCreateUser    = &Errno{Code: 20201, Message: "Error occurred in creating user."}
	ErrUpdateUser    = &Errno{Code: 20202, Message: "Error occurred in updating user"}
	ErrUserNotFound  = &Errno{Code: 20203, Message: "The user was not found."}
	ErrGetUserInfo   = &Errno{Code: 20204, Message: "Error in getting user info"}
	ErrUserInfo      = &Errno{Code: 20205, Message: "The user information json cannot be null"}
	ErrAddSampleData = &Errno{Code: 20206, Message: "Error occurred in joining project"}
	ErrAddLicence    = &Errno{Code: 20207, Message: "Error occurred in grade licence getting"}

	// comment errors
	ErrNotLiked              = &Errno{Code: 20301, Message: "User has not liked yet. "}
	ErrHasLiked              = &Errno{Code: 20302, Message: "User has already liked. "}
	ErrEvaluationList        = &Errno{Code: 20303, Message: "Error occurred while getting evaluation list. "}
	ErrCommentList           = &Errno{Code: 20304, Message: "Error occurred while getting comment list. "}
	ErrGetEvaluation         = &Errno{Code: 20305, Message: "Error occurred while getting evaluation."}
	ErrGetSubCommentInfo     = &Errno{Code: 20306, Message: "Error occurred while getting subComment info"}
	ErrGetParentCommentInfo  = &Errno{Code: 20307, Message: "Error occurred while getting parent comment info"}
	ErrGetHotEvaluations     = &Errno{Code: 20308, Message: "Error occurred while getting hot evaluations"}
	ErrGetHistoryEvaluations = &Errno{Code: 20309, Message: "Error occurred while getting history evaluations"}
	ErrUpdateCourseInfo      = &Errno{Code: 20310, Message: "Error occurred while updating course's info"}
	ErrDeleteComment         = &Errno{Code: 20311, Message: "Error occurred while deleting a comment "}
	ErrHasEvaluated          = &Errno{Code: 20312, Message: "User has evaluated the course"}
	ErrWordLimitation        = &Errno{Code: 20313, Message: "Word limit exceeded"}
	ErrGetEvaluationInfo     = &Errno{Code: 20314, Message: "Error occurred while getting evaluation info"}

	// table errors
	ErrTableExisting    = &Errno{Code: 20401, Message: "The table does not exist"}
	ErrClassExisting    = &Errno{Code: 20402, Message: "The class does not exist"}
	ErrGetTableInfo     = &Errno{Code: 20403, Message: "Error occurred in getting table info. "}
	ErrGetClassInfo     = &Errno{Code: 20404, Message: "Error occurred in getting class info."}
	ErrCourseHasExisted = &Errno{Code: 20405, Message: "This course has already existed in the table."}
	ErrNewTable         = &Errno{Code: 20406, Message: "Error occurred while creating a new table "}

	// message errors
	ErrGetMessage = &Errno{Code: 20501, Message: "Error occurred in getting message list"}

	// search errors
	ErrSearchCourse = &Errno{Code: 20601, Message: "Error occured in searching courses."}

	// upload errors
	ErrGetFile    = &Errno{Code: 20701, Message: "Error occurred in getting file from FormFile()"}
	ErrUploadFile = &Errno{Code: 20702, Message: "Error occurred in uploading file to oss"}

	// course errors
	ErrHistoryCourseExisting = &Errno{Code: 20801, Message: "History Course does not exist."}
	ErrGetSelfCourses        = &Errno{Code: 20802, Message: "Error occurred in getting self courses"}
	ErrSavesDataToLocal      = &Errno{Code: 20803, Message: "Error occurred in Saving data to local"}
	ErrCreateHistoryCourse   = &Errno{Code: 20804, Message: "Error occurred in creating a new history course"}
	ErrFindUsingCourse       = &Errno{Code: 20805, Message: "Error occurred in finding the specific using course"}
	ErrUsingCourseExisting   = &Errno{Code: 20806, Message: "Using Course does not exist."}

	// report errors
	ErrCreateReport = &Errno{Code: 20901, Message: "Error occurred in creating new report."}

	// collection errors
	ErrGetCollections = &Errno{Code: 21001, Message: "Error occurred in getting collections"}
)
