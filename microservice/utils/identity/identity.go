package identity

import "github.com/google/uuid"

func NewUUIDV4() string {
	return uuid.New().String()
}

func IsValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}

func IsNotValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err != nil
}
