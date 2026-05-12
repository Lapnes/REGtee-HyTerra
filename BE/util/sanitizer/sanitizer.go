package sanitizer

// import (
// 	"regexp"
// 	"strings"
// 	"unicode"
// )

// func isASCII(s string) bool {
// 	for i := 0; i < len(s); i++ {
// 		if s[i] > unicode.MaxASCII {
// 			return false
// 		}
// 	}
// 	return true
// }

// func SanitizeUsername(Username string) error {
// 	Username = strings.ToLower(Username)

// 	if len(Username) < constant.UsernameMinLength || len(Username) > constant.UsernameMaxLength {
// 		return constant.ErrUsernameLength
// 	}

// 	UsernameRegex, _ := regexp.Compile(constant.UsernameRegex)
// 	allowedCharacters := UsernameRegex.MatchString(Username)
// 	if !allowedCharacters || !isASCII(Username) {
// 		return constant.ErrUsernameUnallowed
// 	}

// 	if !unicode.IsLetter(rune(Username[0])) {
// 		return constant.ErrUsernamePrefix
// 	}

// 	if strings.HasSuffix(Username, ".") || strings.HasSuffix(Username, "-") || strings.HasSuffix(Username, "_") {
// 		return constant.ErrUsernameSuffix
// 	}

// 	dotCount := strings.Count(Username, ".")
// 	dashCount := strings.Count(Username, "-")
// 	underscoreCount := strings.Count(Username, "_")
// 	if dotCount > 1 || dashCount > 1 || underscoreCount > 1 {
// 		return constant.ErrUsernameSymbolCount
// 	}

// 	return nil
// }

// func SanitizeEmail(email string) error {
// 	email = strings.ToLower(email)
// 	if len(email) > constant.EmailMaxLength {
// 		return constant.ErrorLength(-1, constant.EmailMaxLength)
// 	}

// 	emailRegex, _ := regexp.Compile(constant.EmailRegex)
// 	allowedCharacters := emailRegex.MatchString(email)
// 	if !allowedCharacters || !isASCII(email) {
// 		return constant.EmailUnallowed
// 	}

// 	return nil
// }

// func SanitizePassword(password string) error {
// 	if len(password) < constant.PasswordMinLength || len(password) > constant.PasswordMaxLength {
// 		return constant.ErrPasswordLength
// 	}

// 	hasLetter := false
// 	hasNumber := false
// 	hasUnderscore := false
// 	hasDash := false
// 	hasPeriod := false
// 	hasSpace := false

// 	for _, i := range password {
// 		if i > unicode.MaxASCII {
// 			return constant.ErrPasswordUnallowed
// 		}

// 		if unicode.IsLetter(i) {
// 			hasLetter = true
// 		} else if unicode.IsDigit(i) {
// 			hasNumber = true
// 		} else if unicode.IsSpace(i) {
// 			hasSpace = true
// 		}
// 	}

// 	if strings.IndexByte(password, '_') != -1 {
// 		hasUnderscore = true
// 	} else if strings.IndexByte(password, '-') != -1 {
// 		hasDash = true
// 	} else if strings.IndexByte(password, '.') != -1 {
// 		hasPeriod = true
// 	}

// 	if !hasLetter || !hasNumber {
// 		return constant.ErrPasswordWeak
// 	}
// 	if hasUnderscore || hasDash || hasPeriod || hasSpace {
// 		return constant.ErrPasswordUnallowed
// 	}

// 	return nil
// }
