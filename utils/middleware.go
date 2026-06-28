package utils

import (
	"context"
	"net/http"
	"strings"
)

// Obtiene el rol
func SoloAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtenemos el rol inyectado por MiddlewareJWT
		rol, ok := r.Context().Value("rol").(string)

		if !ok || rol != "admin" {
			http.Error(w, `{"error": "Acceso denegado: requiere rol administrador"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// MiddlewareJWT protege rutas de la API verificando el token
func MiddlewareJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encabezadoAuth := r.Header.Get("Authorization")
		if encabezadoAuth == "" {
			http.Error(w, `{"error": "Falta el token de autorización"}`, http.StatusUnauthorized)
			return
		}

		// El formato debe ser "Bearer <token>"
		partes := strings.Split(encabezadoAuth, " ")
		if len(partes) != 2 || partes[0] != "Bearer" {
			http.Error(w, `{"error": "Formato de token inválido"}`, http.StatusUnauthorized)
			return
		}

		reclamos, err := ValidarToken(partes[1])
		if err != nil {
			http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}

		// Guardamos los datos del usuario en el contexto de la petición
		ctx := context.WithValue(r.Context(), "usuario_id", reclamos.UsuarioID)
		ctx = context.WithValue(ctx, "rol", reclamos.Rol)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
