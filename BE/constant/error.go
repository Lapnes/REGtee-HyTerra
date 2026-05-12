package constant

import (
	"errors"
	"fmt"
)

var (
	ErrInvalid = errors.New("invalid")

	ErrUsernameLength      = errors.New(fmt.Sprintf("must be %d to %d characters long", 6, 60))
	ErrUsernameUnallowed   = errors.New("username cannot contain special characters/ can only contains lowercase characters")
	ErrUsernamePrefix      = errors.New("username must start with a letter")
	ErrUsernameSuffix      = errors.New("username cannot start or end with period, dash, and underscore")
	ErrUsernameSymbolCount = errors.New("cannot contain more than one dot, underscore, or dash")

	ErrExistEmail   = errors.New("email exist")
	ErrUserNotFound = errors.New("account not exist")

	EmailUnallowed = errors.New("email can only use letters, numbers, underscore, and period")

	ErrPasswordLength    = errors.New(fmt.Sprintf("must be %d to %d characters long", 5, 60))
	ErrPasswordUnallowed = errors.New("cannot contain unallowed characters")
	ErrPasswordWeak      = errors.New("must contain letter and number")

	ErrNotFound = errors.New("data not found")

	ErrVerificationInterval = errors.New("too many verification request")

	ErrUserExist = errors.New("user exist")

	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenMismatch = errors.New("token does not match")
	ErrTokenVerified = errors.New("token has been verified")

	ErrWrongCredentials = errors.New("wrong password")
	ErrExpiredSession   = errors.New("expired session")
	ErrInvalidAuth      = errors.New("invalid auth")

	ErrWrongPassword             = errors.New("wrong password")
	ErrUnauthorizedContestAccess = errors.New("unauthorized contest access")

	ErrInvalidExtension = errors.New("invalid file extension")

	ErrFileSize = errors.New("file size exceeds the maximum allowed")

	ErrRabbitMQMaxRetry = errors.New("rabbitmq max retry reached")

	ErrContestEnded = errors.New("contest ended")
)

func ErrorLength(min, max int) (err error) {
	if min < 0 && max < 0 {
		min, max = -min, -max
	}
	if min > max && max >= 0 {
		min, max = max, min
	}
	var errMsg string
	if min == max {
		errMsg = fmt.Sprintf("must be %d characters", min)
	} else if max < 0 {
		errMsg = fmt.Sprintf("must be more than %d characters", min)
	} else if min < 0 {
		errMsg = fmt.Sprintf("must be less than %d characters", max)
	} else {
		errMsg = fmt.Sprintf("must be between %d and %d characters", min, max)
	}

	err = errors.New(errMsg)
	return
}
