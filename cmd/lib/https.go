package lib

import (
	"crypto/tls"
	"log"
	"net/http"
)

// SetupHTTPS configure le serveur pour utiliser HTTPS avec un certificat auto-signé
func SetupHTTPS(server *http.Server) {
	// Chargement du certificat et de la clé privée
	certFile := "server.crt" // Chemin vers le fichier de certificat
	keyFile := "server.key"  // Chemin vers le fichier de clé privée

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("Erreur lors du chargement du certificat: %v", err)
	}

	// Configuration TLS avec le certificat chargé
	server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Redirection de HTTP vers HTTPS
	go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
	}))
}
