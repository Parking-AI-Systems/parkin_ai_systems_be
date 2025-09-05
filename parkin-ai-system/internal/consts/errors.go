package consts

import (
	"fmt"
	"net/http"
)

var (
	// Auth errors
	CodeInvalidToken       = customCode{code: 201, message: "Invalid token", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeTokenExpired       = customCode{code: 202, message: "Token expired", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeUnauthorized       = customCode{code: 203, message: "Unauthorized", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeNotAdmin           = customCode{code: 204, message: "Not admin", detail: nil, httpStatus: http.StatusForbidden}
	CodeNotOwner           = customCode{code: 205, message: "Not owner", detail: nil, httpStatus: http.StatusForbidden}
	CodeInvalidInput       = customCode{code: 206, message: "Invalid input", detail: nil, httpStatus: http.StatusBadRequest}
	CodeHashPasswordFailed = customCode{code: 207, message: "Failed to hash password", detail: nil, httpStatus: http.StatusInternalServerError}
	// User errors
	CodeUserNotFound      = customCode{code: 301, message: "User not found", detail: nil, httpStatus: http.StatusNotFound}
	CodeEmailExists       = customCode{code: 302, message: "Email already exists", detail: nil, httpStatus: http.StatusConflict}
	CodeIncorrectPassword = customCode{code: 303, message: "Incorrect password", detail: nil, httpStatus: http.StatusUnauthorized}
	CodeMaxUsersReached   = customCode{code: 304, message: "Max users reached", detail: nil, httpStatus: http.StatusBadRequest}

	// Vehicle errors
	CodeVehicleNotFound    = customCode{code: 401, message: "Vehicle not found", detail: nil, httpStatus: http.StatusNotFound}
	CodeMaxVehiclesReached = customCode{code: 402, message: "Max vehicles reached (5)", detail: nil, httpStatus: http.StatusBadRequest}
	CodeLicensePlateExists = customCode{code: 403, message: "License plate already exists", detail: nil, httpStatus: http.StatusConflict}

	// Parking lot errors
	CodeParkingLotNotFound = customCode{code: 501, message: "Parking lot not found", detail: nil, httpStatus: http.StatusNotFound}
	CodeLocationExists     = customCode{code: 502, message: "Location already exists", detail: nil, httpStatus: http.StatusConflict}
	CodeInvalidTimeFormat  = customCode{code: 503, message: "Invalid time format", detail: nil, httpStatus: http.StatusBadRequest}

	// Favourite errors
	CodeFavouriteNotFound = customCode{code: 601, message: "Favourite not found", detail: nil, httpStatus: http.StatusNotFound}
	CodeFavouriteExists   = customCode{code: 602, message: "Already in favourites", detail: nil, httpStatus: http.StatusConflict}

	// Database errors
	CodeDatabaseError    = customCode{code: 901, message: "Database error", detail: nil, httpStatus: http.StatusInternalServerError}
	CodeFailedToCreate   = customCode{code: 902, message: "Failed to create", detail: nil, httpStatus: http.StatusInternalServerError}
	CodeFailedToUpdate   = customCode{code: 903, message: "Failed to update", detail: nil, httpStatus: http.StatusInternalServerError}
	CodeFailedToDelete   = customCode{code: 904, message: "Failed to delete", detail: nil, httpStatus: http.StatusInternalServerError}
	CodeCannotDeleteSelf = customCode{code: 905, message: "Cannot delete yourself", detail: nil, httpStatus: http.StatusBadRequest}

	CodeParkingSlotNotFound = customCode{code: 1001, message: "Parking slot not found", detail: nil, httpStatus: http.StatusNotFound}
	CodeAlreadyFavorited    = customCode{code: 1002, message: "Parking lot already favorited", detail: nil, httpStatus: http.StatusConflict}
	CodeFavoriteNotFound    = customCode{code: 1003, message: "Favorite not found", detail: nil, httpStatus: http.StatusNotFound}
)

type customCode struct {
	code       int
	message    string
	detail     interface{}
	httpStatus int
}

// Code returns the integer number of current error code.
func (c customCode) Code() int {
	return c.code
}

// Message returns the brief message for current error code.
func (c customCode) Message() string {
	return c.message
}

// Detail returns the detailed information of current error code,
// which is mainly designed as an extension field for error code.
func (c customCode) Detail() interface{} {
	return c.detail
}

// String returns current error code as a string.
func (c customCode) String() string {
	if c.detail != nil {
		return fmt.Sprintf(`%d:%s %v`, c.code, c.message, c.detail)
	}
	if c.message != "" {
		return fmt.Sprintf(`%d:%s`, c.code, c.message)
	}
	return fmt.Sprintf(`%d`, c.code)
}

func (c customCode) HttpStatus() int {
	return c.httpStatus
}
