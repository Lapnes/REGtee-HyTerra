package hash

import (
	variable "backend/constant"
	"crypto/sha256"
	"encoding/hex"
)

func VerifyFileIntegrity(str string, salt string, expectedHash string) bool {
	fileContents := []byte(str + salt + variable.PRIVATE_KEY)
	hash := sha256.Sum256(fileContents)
	return hex.EncodeToString(hash[:]) == expectedHash
}

func HashPassword(str string) string {
	conc := sha256.Sum256([]byte(str))
	return hex.EncodeToString(conc[:])
}
