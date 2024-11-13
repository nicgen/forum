package handlers

import (
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strconv"
)

func FiltersHandler(w http.ResponseWriter, r *http.Request) {
	// Getting the databa
}

func Category(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()

	// Getting the number of values stored in the Categories tables
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Categories").Scan(&count)
	if err != nil {
		// Handle error
	}

	// numbercat := map[string]bool{}
	formValues := r.URL.Query()

	// for i := 0; i < count; i++ {
	// 	compteur := strconv.Itoa(i)
	// 	notempty := formValues.Get("Category" + compteur)
	// 	if len(notempty) != 0 {
	// 		numbercat["Category"+compteur] = true
	// 	} else {
	// 		numbercat["Category"+compteur] = false
	// 	}
	// }
	category := ""
	for i := 0; i < count; i++ {
		compteur := strconv.Itoa(i)
		notempty := formValues.Get("Category" + compteur)
		if len(notempty) != 0 {
			category = notempty
		}
	}

	// Category_ID1 := formValues.Get("Category_ID1")
	// Category_ID2 := formValues.Get("Category_ID2")
	// Category_ID3 := formValues.Get("Category_ID3")
	// Category_ID4 := formValues.Get("Category_ID4")

	// Users posts Request
	var posts []*models.Post
	var err_post error

	// Making the database request
	state_category := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts WHERE Category_ID = ?`
	rows, err_post := db.Query(state_category, category)

	if err_post != nil {

	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {

		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
	}

	// db.Exec()
	// db.QueryRow()
	// db.Query(state_category, Category_ID)
	data := lib.Cook(w, r)
	data["Posts"] = posts
}

func Like(w http.ResponseWriter, r *http.Request) {
	// Getting the database data
	db := lib.GetDB()
	minlike, maxlike := 0, 0
	nothing := false
	formValues := r.URL.Query()
	numberlike := formValues.Get("Like")
	switch numberlike {
	case "1-10":
		minlike = 1
		maxlike = 10

	case "11-50":
		minlike = 11
		maxlike = 50

	case "51-100":
		minlike = 51
		maxlike = 100

	case "+100":
		minlike = 101

	default:
		nothing = true
	}
	// Users posts Request
	var posts []*models.Post
	var err_post error

	// Making the database request
	state_like := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts`
	rows, err_post := db.Query(state_like)

	if err_post != nil {

	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {

		}
		if (post.Like > minlike && post.Like < maxlike) || nothing {
			posts = append(posts, &post)
		}
	}

	if err := rows.Err(); err != nil {
	}
	data := lib.Cook(w, r)
	data["Posts"] = posts
}

func Dislike(w http.ResponseWriter, r *http.Request) {
	// Getting the database data
	db := lib.GetDB()
	mindislike, maxdislike := 0, 0
	nothing := false
	formValues := r.URL.Query()
	numberdislike := formValues.Get("Dislike")
	switch numberdislike {
	case "1-10":
		mindislike = 1
		maxdislike = 10

	case "11-50":
		mindislike = 11
		maxdislike = 50

	case "51-100":
		mindislike = 51
		maxdislike = 100

	case "+100":
		mindislike = 101

	default:
		nothing = true
	}
	// Users posts Request
	var posts []*models.Post
	var err_post error

	// Making the database request
	state_dislike := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts`
	rows, err_post := db.Query(state_dislike)

	if err_post != nil {

	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {

		}
		if (post.Dislike > mindislike && post.Dislike < maxdislike) || nothing {
			posts = append(posts, &post)
		}
	}

	if err := rows.Err(); err != nil {
	}
	data := lib.Cook(w, r)
	data["Posts"] = posts
}

func CreationDate(w http.ResponseWriter, r *http.Request) {
	// db := lib.GetDB()
	// start := time.Now().Format("2006-01-02 15:04:05")

	// Time := strings.Split(start, " ")
	// Date := strings.Split(Time[0], "-")
	// Hours := strings.Split(Time[1],":")

	// Year := Date[0]
	// Month := Date[1]
	// Day := Date[2]

	// Hour := Hours[0]
	// Min := Hours[1]
	// PENSE BETE : ATTENTION AU DEBUT DE MOIS -7 && meme chose pour debut d'année (janvier)

}
