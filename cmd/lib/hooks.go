package lib

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	ntfy_token string
)

// used to initialize the `ntfy_token` global variable
func init() {
	// load env file
	LoadEnv(".env")
	ntfy_token = os.Getenv("NTFY_token")
}

func PostItOnNfty(msg string) {
	url := "https://ntfy.sh/" + ntfy_token
	fmt.Println(url)
	http.Post(url, "text/plain",
		strings.NewReader(msg))
}
