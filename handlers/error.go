package handlers

import (
	"fmt"
	"forum/cmd/lib"
	"forum/models"
	"html/template"
	"log"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
)

// HandleError handles error requests and send an error response with the given status code and message
func HandleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	data := models.PageData{
		Title:  "Error",
		Header: fmt.Sprintf("Error %d", statusCode),
		// Content:   message,
		Content: map[string]template.HTML{
			"Msg_raw":    template.HTML("<h1>" + message + "</>"),
			"Msg_styled": template.HTML("<h1 style=\"text=color: blue;\">" + strconv.Itoa(statusCode) + "</h1><p>paragraph with</br>style</>"),
		},
		IsError:   true,
		ErrorCode: statusCode,
	}
	lib.RenderTemplate(w, "layout/error", "page/error", data)
}

// WithErrorHandling middleware that handles all errors and panics
func WithErrorHandling(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// log the error and stack trace
				log.Printf("Panic: %v\n%s", err, debug.Stack())

				// default status code and message
				statusCode := http.StatusInternalServerError
				message := "Internal Server Error"

				// switch to determine the status code and message
				switch e := err.(type) {
				case *models.CustomError: // handles your custom error type
					// fmt.Println(">>>>>>>>>>CUSTOM<<<<<<<<<<")
					statusCode = e.StatusCode
					message = e.Message
				case *net.OpError: // for network-related errors
					statusCode = http.StatusInternalServerError
					message = "A network error occurred"
				case string: // direct string panics
					// fmt.Println(">>>>>>>>>>STRING<<<<<<<<<<")
					switch e {
					case "bad request":
						statusCode = http.StatusBadRequest
						message = "Bad Request"
					case "not found":
						statusCode = http.StatusNotFound
						message = "Not Found"
					case "too many request":
						statusCode = http.StatusTooManyRequests
						message = "Too Many Request"
					case "see other":
						statusCode = http.StatusSeeOther
						message = "See Other"
					case "Unauthorized":
						statusCode = http.StatusUnauthorized
						message = "Status Unauthorized"
					default:
						// fmt.Println(">>>>>>>>>>STRING:Internal<<<<<<<<<<")
						statusCode = http.StatusInternalServerError
						message = "Internal Server Error 💀"
					}
				default: // handle other types of panics
					// fmt.Println(">>>>>>>>>>PANIC<<<<<<<<<<")
					statusCode = http.StatusInternalServerError
					message = "Internal Server Error"
				}

				// ? Optional: Post an alert on [NTFY](https://ntfy.sh/)
				ntfy_title := "[FORUM server]"
				ntfy_msg := message
				lib.PostItOnNfty(ntfy_title, ntfy_msg)
				// render the error page
				HandleError(w, statusCode, message)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// ! Force500Handler forces a 500 error (for testing purposes)
func Force500Handler(w http.ResponseWriter, r *http.Request) {

	// panic("This is a forced panic to test 500 error handling")

	// panic(&models.CustomError{ // custom error
	// 	StatusCode: http.StatusInternalServerError,
	// 	Message:    "Oh, snap! Internal Server Error",
	// })

	OutOfRange() // Panic: runtime error: index out of range
}

// Oh SNAP!
func OutOfRange() {
	t := []int{1, 2, 3}
	for i := 0; i < 4; i++ {
		fmt.Println(t[i])
	}
}

// ! ForceDirectError forces a direct string error (for testing purposes)
func ForceDirectError(w http.ResponseWriter, r *http.Request) {
	// panic("This is a forced panic to test the direct string error") // string(default)
	panic("bad request")
	// panic("not found")
	// panic("Internal Server Error")
}
