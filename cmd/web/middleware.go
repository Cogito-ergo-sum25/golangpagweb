package main

import (
	"net/http"
	"github.com/justinas/nosurf"
)

// NoSurf is the csrf protection middleware
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves session data for current request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// Auth es el middleware que requiere que un usuario esté autenticado.
// Se ejecutará después de SessionLoad.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verificamos si la llave "user_id" existe en la sesión actual.
		if !session.Exists(r.Context(), "user_id") {
			// Si no existe, el usuario NO ha iniciado sesión.
			// Guardamos un mensaje para mostrar en la página de login.
			session.Put(r.Context(), "error", "Por favor, inicia sesión para acceder.")
			
			// Redirigimos al usuario a la página de login.
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return // Importante: detenemos la ejecución aquí.
		}
		
		// Si la sesión existe, el usuario está autenticado.
		// Le permitimos continuar a la siguiente ruta (ej. /inventario).
		next.ServeHTTP(w, r)
	})
}
