package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Alitm23/SistemaEcommerce/models"
)

// ApiVerCarrito devuelve los artículos actuales en el carrito del usuario (Servicio Web 4)
func ApiVerCarrito(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//extrae el usuario_id del contexto de la petición, que fue establecido por el middleware JWT
	usuarioID, ok := r.Context().Value("usuario_id").(int)
	// r.Context() devuelve un contexto que contiene información sobre la petición HTTP actual. El middleware JWT coloca el usuario_id en este contexto después de validar el token JWT.
	// Aquí, se intenta extraer ese valor y convertirlo a un entero.
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no autorizado"})
		return
	}

	// Busca el carrito activo de este usuario
	carritoID, err := models.ObtenerCarritoActivo(usuarioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener el carrito"})
		return
	}

	// devolver los items del carrito con sus nombres y detalles
	items, err := models.ObtenerItemsCarrito(carritoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al obtener los items"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

// ApiAgregarItem recibe un JSON del frontend para añadir una joya (Servicio Web 6)
func ApiAgregarItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Validamos la identidad del usuario a través del token JWT
	usuarioID, ok := r.Context().Value("usuario_id").(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Usuario no autorizado"})
		return
	}

	// Crea un struct temporal solo para atrapar los datos que envía el frontend
	var peticion struct {
		ProductoTallaID int `json:"producto_talla_id"`
		Cantidad        int `json:"cantidad"`
	}

	// Decodifica el JSON recibido en la variable peticion
	if err := json.NewDecoder(r.Body).Decode(&peticion); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato JSON inválido"})
		return
	}

	carritoID, err := models.ObtenerCarritoActivo(usuarioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al verificar el carrito activo"})
		return
	}

	nuevoItem := models.ItemCarrito{
		CarritoID:       carritoID,
		ProductoTallaID: peticion.ProductoTallaID,
		Cantidad:        peticion.Cantidad,
	}

	// Agregamos el item al carrito y notificamos al canal de analítica
	if err := models.AgregarItem(nuevoItem); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated) // Recibido y creado correctamente
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Joya añadida exitosamente"})
}
