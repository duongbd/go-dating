package hashedPassword

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
func Generate(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
