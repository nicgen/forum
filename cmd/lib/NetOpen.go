package lib

import (
	"log"
	"os/exec"
	"runtime"
)

// OpenBrowser ouvre l'URL dans le navigateur par défaut
func OpenBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		log.Printf("Unsupported platform, please open the URL manually: %s", url)
		return
	}
	if err != nil {
		log.Printf("Erreur lors de l'ouverture du navigateur: %v", err)
	}
}
