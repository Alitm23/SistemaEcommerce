package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/Alitm23/SistemaEcommerce/models"

	"github.com/gorilla/mux"
)

// Función que renderiza el formulario de registro de un nuevo usuario
func MostrarRegistro(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/usuario/registro.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Recibe los datos del formulario y crea un nuevo usuario en la BD
func ProcesarRegistro(w http.ResponseWriter, r *http.Request) {
	// Parsear los datos enviados por el formulario HTML
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Construir el usuario con el constructor del modelo
	// El constructor aplica el hash a la contraseña internamente
	usuario, err := models.RegistrarUsuario(
		r.FormValue("nombre"),
		r.FormValue("email"),
		r.FormValue("password"),
		"cliente",
		r.FormValue("direccion"),
		r.FormValue("telefono"),
	)
	if err != nil {
		http.Error(w, "Datos inválidos: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Persistir el usuario en la base de datos
	if err := usuario.Registrar(); err != nil {
		http.Error(w, "Error al registrar el usuario: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirigir al login tras el registro exitoso
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Renderizar el formulario de inicio de sesión
func MostrarLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/usuario/login.html")
	if err != nil {
		http.Error(w, "Error al cargar el formulario", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// Validar las credenciales del usuario contra la base de datos
func ProcesarLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("contrasenia")

	// Buscar el usuario por email — retorna (usuario, bool)
	usuario, ok := models.BuscarPorEmail(email)
	if !ok {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	// Verificar la contraseña usando bcrypt a través del método del modelo
	if err := usuario.Autenticar(password); err != nil {
		http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
		return
	}

	// Credenciales correctas — redirigir al inicio
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Función obtiene todos los usuarios y los envía al template de lista
func ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	// Consultar todos los usuarios desde la base de datos
	usuarios, err := models.ListarUsuarios()
	if err != nil {
		http.Error(w, "Error al obtener los usuarios: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/usuario/lista.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	// Enviar la lista de usuarios al template para renderizar
	tmpl.Execute(w, usuarios)
}

// Función para cargar los datos del usuario y renderiza el formulario de edición.
func MostrarEdicion(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID desde la URL: /usuarios/{id}/editar
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Buscar el usuario — si no existe retorna false
	usuario, ok := models.BuscarPorID(id)
	if !ok {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/usuario/editar.html")
	if err != nil {
		http.Error(w, "Error al cargar el template", http.StatusInternalServerError)
		return
	}

	// Pasar el usuario al template para prellenar el formulario
	tmpl.Execute(w, usuario)
}

// actualiza los datos del usuario con los valores del formulario.
func ProcesarEdicion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error al procesar el formulario", http.StatusBadRequest)
		return
	}

	// Cargar el usuario actual desde la BD antes de modificarlo
	usuario, ok := models.BuscarPorID(id)
	if !ok {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Actualizar los campos con los valores del formulario
	usuario.Nombre = r.FormValue("nombre")
	usuario.Email = r.FormValue("email")
	usuario.Direccion = r.FormValue("direccion")
	usuario.Telefono = r.FormValue("telefono")

	// Cambiar el rol solo si se envió un valor — CambiarRol valida internamente
	if nuevoRol := r.FormValue("rol"); nuevoRol != "" {
		if err := usuario.CambiarRol(nuevoRol); err != nil {
			http.Error(w, "Rol inválido: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Cambiar la contraseña solo si el campo no viene vacío
	if nuevaPassword := r.FormValue("password"); nuevaPassword != "" {
		if err := usuario.CambiarPw(nuevaPassword); err != nil {
			http.Error(w, "Error al cambiar la contraseña: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Persistir los cambios en la base de datos
	if err := usuario.Actualizar(); err != nil {
		http.Error(w, "Error al actualizar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

// Función para eliminar un usuario de la base de datos por su ID
func EliminarUsuario(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	// Solo necesitamos el ID para eliminar — no hace falta cargar el usuario completo
	usuario := &models.Usuario{ID: id}
	if err := usuario.Eliminar(); err != nil {
		http.Error(w, "Error al eliminar: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}
