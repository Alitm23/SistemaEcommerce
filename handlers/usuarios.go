package handlers

import (
	"net/http"
	"net/url"
	"time"

	"github.com/Alitm23/SistemaEcommerce/models"
	"github.com/Alitm23/SistemaEcommerce/utils"
)

type authViewData struct {
	Error string
	Email string
}

// Registro maneja la solicitud de registro de un nuevo usuario.
func Registro(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		usuario := models.Usuario{
			Nombre:   r.FormValue("nombre"),
			Apellido: r.FormValue("apellido"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}
		if err := models.RegistrarUsuario(usuario); err != nil {
			renderTemplate(w, http.StatusBadRequest, []string{"templates/base.html", "templates/login.html"}, "base", authViewData{
				Error: err.Error(),
				Email: usuario.Email,
			})
			return
		}
		http.Redirect(w, r, "/login?email="+url.QueryEscape(usuario.Email), http.StatusSeeOther)
		return
	}
	renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/login.html"}, "base", authViewData{
		Email: r.URL.Query().Get("email"),
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		usuario, err := models.AutenticarUsuario(email, password)
		if err != nil {
			renderTemplate(w, http.StatusUnauthorized, []string{"templates/base.html", "templates/login.html"}, "base", authViewData{
				Error: "Credenciales invalidas",
				Email: email,
			})
			return
		}

		token, err := utils.GenerarToken(usuario.ID, usuario.Rol)
		if err != nil {
			http.Error(w, "Error al iniciar sesion", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		if usuario.Rol == "admin" {
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/productos", http.StatusSeeOther)
		return
	}

	renderTemplate(w, http.StatusOK, []string{"templates/base.html", "templates/login.html"}, "base", authViewData{})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
