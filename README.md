# Sistema de E-commerce de Joyería y Accesorios Artesanales
### Datos del proyecto
| | |
|---|---|
| **Nombre del Proyecto** | Sistema de Gestión E-commerce |
| **Estudiante** | Nataly Tituaña |
| **Fecha** | 24 de mayo de 2026 |
| **Materia** | Programación Orientada a Objetos |
| **Docente** | Ing. Milton Palacios 

## Introducción
Este proyecto es un Sistema de E-commerce destinado a la venta de Joyería y Accesorios Artesanales, implementado en el lenguaje de programación Go. Desarrollado como parte de la materia de Programación Orientada a Objetos, el sistema busca aplicar los conceptos de diseño de software para gestionar de forma eficiente una tienda en línea.

## Objetivo del Sistema
Desarrollar una plataforma digital de e-commerce que permita a los clientes explorar productos artesanales, gestionar su carrito de compras y simular el proceso de pago. A su vez, provee a los administradores las herramientas necesarias para controlar el inventario y gestionar el catálogo de productos.

## Principales Funcionalidades
* **Gestión de Productos:** Permite listar productos con stock disponible (mayor a cero) y visualizar detalles específicos como nombre, precio, descripción, categoría y stock actual.
* **Gestión del Carrito de Compras:** Los usuarios pueden agregar productos especificando la cantidad. Si el producto ya está en el carrito, la cantidad se incrementa automáticamente. Permite realizar el *checkout*, cambiando el estado a "pagado" y generando una orden.
* **Usuarios y Autenticación:** Cuenta con un sistema de registro e inicio de sesión seguro que devuelve un token JWT. Soporta roles diferenciados: Administrador (acceso total) y Cliente (acceso limitado a su carrito y pedidos).
* **Simulación de Pagos:** Implementa un procesamiento simulado que aprueba el pago automáticamente si existe stock suficiente, actualizando el estado del pedido (pendiente, pagado o cancelado).

## Módulos del Sistema

1. **Catálogo e Inventario:** Encargado del registro, actualización, eliminación y consulta de productos. También maneja la reducción de stock al confirmar un pago.
2. **Carrito:** Gestiona la agregación y eliminación de ítems, el cálculo del total de la compra y el proceso de checkout.
3. **Ordenes y Pagos:** Módulo responsable de generar la orden de compra, procesar el pago simulado y gestionar los cambios de estado de la orden.
4. **Usuarios:** Controla el registro, el inicio de sesión y la gestión básica de los roles de acceso.

## Paquetes Utilizados
El proyecto está desarrollado en **Go (Golang)** y hace uso de las siguientes dependencias:
**Paquetes Externos:**
* **Driver PostgreSQL (`gorm.io/driver/postgres`):** Controlador para la conexión con la base de datos relacional PostgreSQL.
* **Bcrypt (`golang.org/x/crypto/bcrypt`):** Algoritmo de hash para encriptar contraseñas y evitar que se guarden en texto plano.
* **Gin (`github.com/gin-gonic/gin`):** Framework web de alto rendimiento para el enrutamiento eficiente de la API.
* **GORM (`gorm.io/gorm`):** ORM utilizado para mapear structs a tablas, construir consultas y generar migraciones automáticas.
* **JWT (`github.com/golang-jwt/jwt/v5`):** Utilizado para la generación y validación de tokens garantizando una autenticación segura.
* **Godotenv (`github.com/joho/godotenv`):** Para la carga y configuración de variables de entorno desde un archivo `.env`.

**Paquetes Nativos de Go:**
* `encoding/json`: Para la conversión de structs a formato JSON en la comunicación de la API REST.
* `net/http`: Servidor HTTP base de la aplicación.
* `html/template`: Para el renderizado de plantillas HTML en el servidor.
