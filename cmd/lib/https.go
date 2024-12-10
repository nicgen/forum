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
		MinVersion:   tls.VersionTLS12, // Force au minimum TLS 1.2
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		},
		PreferServerCipherSuites: true, // Privilégier les suites du serveur
	}

	// Redirection de HTTP vers HTTPS
	go func() {
		err := http.ListenAndServe(":81", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
		}))
		if err != nil {
			log.Fatalf("Erreur lors de l'écoute HTTP: %v", err)
		}
	}()
}
