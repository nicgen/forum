package lib

import (
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

func PostItOnNfty(title, msg string) {
	url := "https://ntfy.sh/" + ntfy_token
	// fmt.Println(url)
	req, _ := http.NewRequest("POST", url, strings.NewReader(msg))
	req.Header.Set("Title", title)
	// http.Post(url, "text/plain",
	// 	strings.NewReader(msg))
	http.DefaultClient.Do(req)
}
