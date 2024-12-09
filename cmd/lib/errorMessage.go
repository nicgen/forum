package lib

import (
	"net/http"
)

// ? Function to load add an error to the data map
func ErrorMessage(w http.ResponseWriter, data map[string]interface{}, errtype string) map[string]interface{} {
	// Defining the error map that will contain the error we need to send
	errorMap := map[string]interface{}{
		"LoginMail":              "null",
		"LoginPassword":          "null",
		"RegisterUsername":       "null",
		"RegisterEmail":          "null",
		"EmailFormat":            "null",
		"RegisterPassword":       "null",
		"PasswordMatch":          "null",
		"PostAlreadyLiked":       "null",
		"PostAlreadyDisliked":    "null",
		"CommentAlreadyLiked":    "null",
		"CommentAlreadyDisliked": "null",
		"ReportExist":            "null",
		"ContentEmpty":           "null",
		"ImageContent":           "",
		"LogSession": "null",
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
	case "PostAlreadyLiked":
		errorMap["PostAlreadyLiked"] = "You already liked this post"
	case "PostAlreadyDisliked":
		errorMap["PostAlreadyDisliked"] = "You already disliked this post"
	case "CommentAlreadyLiked":
		errorMap["CommentAlreadyLiked"] = "You already liked this comment"
	case "CommentAlreadyDisliked":
		errorMap["CommentAlreadyDisliked"] = "You already disliked this comment"
	case "ReportExist":
		errorMap["ReportExist"] = "Report already completed"
	case "ContentEmpty":
		errorMap["ContentEmpty"] = "Modify must be at least 1 character long"
	case "ImageContent":
		errorMap["ImageContent"] = "Image size: > 20Mo"
	case "LogSession":
		errorMap["LogSession"] = "Session already active"	
		// log.Printf("ImageContent")
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
	RenderTemplate(w, "layout/index", "page/errorServer", data)
}
