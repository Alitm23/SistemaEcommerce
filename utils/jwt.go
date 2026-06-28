package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Clave
var jwtKey = []byte("sistemaecommerce78968")

// Claims personalizados para incluir el Rol
type ReclamosEcommerce struct {
	UsuarioID int    `json:"usuario_id"`
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}

// GenerarToken crea un token válido por 24 horas
func GenerarToken(usuarioID int, rol string) (string, error) {
	tiempoExpiracion := time.Now().Add(24 * time.Hour) // Establecemos la expiración del token
	reclamos := &ReclamosEcommerce{                    // en la variable reclamos se almacenan los datos del usuario y el rol
		UsuarioID: usuarioID,
		Rol:       rol,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tiempoExpiracion),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, reclamos)
	return token.SignedString(jwtKey)
}

// ValidarToken verifica la firma y extrae los datos
func ValidarToken(tokenString string) (*ReclamosEcommerce, error) {
	reclamos := &ReclamosEcommerce{}
	// ParseWithClaims analiza el token y llena la estructura reclamos
	token, err := jwt.ParseWithClaims(tokenString, reclamos, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) { // Si el token ha expirado, retornamos un error específico
			return nil, errors.New("el token ha expirado")
		}
		return nil, errors.New("token inválido")
	}

	if !token.Valid {
		return nil, errors.New("token inválido")
	}

	return reclamos, nil
}
