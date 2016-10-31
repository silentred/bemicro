package user

import "golang.org/x/crypto/bcrypt"

// PasswordHash only uses bcrypt
func PasswordHash(password []byte, cost int) ([]byte, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return hashedPassword, err
	}

	return hashedPassword, nil
}

// PasswordVerify only use bcrypt
func PasswordVerify(password, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}
