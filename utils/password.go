package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword genera un hash seguro para la contraseña usando bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err // Retorna el hash y el error si lo hay
}

// CheckPassword compara un hash con una contraseña para verificar su validez
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
