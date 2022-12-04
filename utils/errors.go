package utils

type Resterr struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func BadRequest(message string) *Resterr {
	return &Resterr{
		Message: message,
		Status:  400,
		Error:   "Bad Request",
	}
}

func NotFound(message string) *Resterr {
	return &Resterr{
		Message: message,
		Status:  404,
		Error:   "Not Found",
	}
}

func InternalErr(message string) *Resterr {
	return &Resterr{
		Message: message,
		Status:  500,
		Error:   "Internal Server Error",
	}
}
