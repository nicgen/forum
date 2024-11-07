package handlers

import (
	"net/http"
)

// ? Function to load add an error to the data map
func ErrorMessage(w http.ResponseWriter, data map[string]interface{}, errtype string) map[string]interface{} {
	// Defining the error map that will contain the error we need to send
	errorMap := map[string]interface{}{
		"LoginMail":        "null",
		"LoginPassword":    "null",
		"RegisterUsername": "null",
		"RegisterEmail":    "null",
		"EmailFormat":      "null",
		"RegisterPassword": "null",
		"PasswordMatch":    "null",
	}

	// Changing the error depending on the errtype variable given
	switch errtype {
	case "LoginMail":
		errorMap["LoginMail"] = "No account registered with this email"
	case "LoginPassword":
		errorMap["LoginPassword"] = "Invalid password"
	case "RegisterUsername":
		errorMap["RegisterUsername"] = "Username already taken"
	case "RegisterEmail":
		errorMap["RegisterEmail"] = "Email already taken"
	case "EmailFormat":
		errorMap["EmailFormat"] = "Invalid Email format"
	case "RegisterPassword":
		errorMap["RegisterPassword"] = "Invalid password, password must contain 8 characters, at least 1 special character and a number"
	case "PasswordMatch":
		errorMap["PasswordMatch"] = "Password doesn't match"
	}

	// Adding the errorMap to the data map
	data["Error"] = errorMap

	// Returning the map with the error added
	return data
}

// ? Function to add error to the Error map for the index
func ErrorServer(w http.ResponseWriter, errtype string) {
	// Storing the error into a map
	data := map[string]interface{}{
		"ErrorServer": errtype,
	}

	// Load the page with the Error sent on the map
	renderTemplate(w, "layout/index", "page/errorServer", data)
}
