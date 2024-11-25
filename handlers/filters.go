package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strconv"
	"time"
)

func FiltersHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()

	formValues := r.URL.Query()

	// categories := formValues.Get("Category")
	numberlike := formValues.Get("Like")
	numberdislike := formValues.Get("Dislike")
	period := formValues.Get("Period")

	println("Like: ", numberlike)
	println("Dislike: ", numberdislike)
	println("Period: ", period)

	// Getting the number of values stored in the Categories tables
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Categories").Scan(&count)
	if err != nil {
		//Erreur critique: échec de la récuperation du nombre de catégorie
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving categories count. Please try again.",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}

	// filter category
	category := ""
	for i := 0; i < count; i++ {
		compteur := strconv.Itoa(i)
		notempty := formValues.Get("Category" + compteur)
		if len(notempty) != 0 {
			category = notempty
		}
	}

	// Filter like
	likemin, likemax := 0, 0
	switch numberlike {
	case "1-10":
		likemin = 1
		likemax = 10
	case "11-50":
		likemin = 11
		likemax = 50
	case "51-100":
		likemin = 51
		likemax = 100
	case "plus-de-100":
		likemin = 101
	}

	// filter dislike
	dislikemin, dislikemax := 0, 0
	switch numberdislike {
	case "1-10":
		dislikemin = 1
		dislikemax = 10
	case "11-50":
		dislikemin = 11
		dislikemax = 50
	case "51-100":
		dislikemin = 51
		dislikemax = 100
	case "plus-de-100":
		dislikemin = 101
	}

	//filter period
	lastperiod := ""
	switch period {
	case "week":
		lastperiod = "-7 days"
	case "month":
		lastperiod = "start of month"
	case "year":
		lastperiod = "start of year"
	}
	var posts []*models.Post
	var err_post error

	state_filters :=
		`SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID
	FROM posts
	WHERE 
  (Category_ID = ?)
	AND
	(Like BETWEEN ? AND ? OR (? = 0 AND ? = 0)) 
	AND 
	(Dislike BETWEEN ? AND ? OR (? = 0 AND ? = 0))
	AND 
	(CreatedAt >= DATE('now', ?) AND CreatedAt < DATE('now', 'start of month', '+1 month')) 
	 `
	rows, err_post := db.Query(state_filters, category, likemin, likemax, likemin, likemax, dislikemin, dislikemax, dislikemin, dislikemax, lastperiod)

	if err_post != nil {
		//Erreur critique: échec de la récuperation des posts
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving posts. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
			//Erreur non critique: échec de la lecture des données du post
			lib.ErrorServer(w, "Error scanning post data. Please try again later.")
			return
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		//Erreur non critique: erreur lors du traitement des lignes
		lib.ErrorServer(w, "Error processing posts. Please try again later.")
	}
	data := lib.DataTest(w, r)
	data["Posts"] = posts

	// Redirect User to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func Category(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()

	// Getting the number of values stored in the Categories tables
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM Categories").Scan(&count)
	if err != nil {
		//Erreur critique : échec de la récuperation du nombre de catégories.
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving categories count.Please try again.",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
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
		//Erreur critique : echec de la récupération des posts par categories.
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving posts by category. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
			//Erreur non critique : échec de la lecture des données du post.
			lib.ErrorServer(w, "Error scanning post data. Please try again later.")
			return
		}

		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		//Erreur non critique: Erreur lors du traitement des lignes
		lib.ErrorServer(w, "Error processing posts. Please try again later.")
		return
	}

	// db.Exec()
	// db.QueryRow()
	// db.Query(state_category, Category_ID)
	data := lib.DataTest(w, r)
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
	var rows *sql.Rows
	var err_post error

	// Making the database request
	if !nothing {
		state_like := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts WHERE Like BETWEEN ? AND ?`
		rows, err_post = db.Query(state_like, minlike, maxlike)
	} else {
		state_like := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts`
		rows, err_post = db.Query(state_like)
	}

	if err_post != nil {
		//Erreur critique : Echec de la recuperation des posts par like
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving posts by like. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
		return

	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
			//Erreur non critique : échec de la lecture des données du posts
			lib.ErrorServer(w, "Error scanning post data.Please try again")
			return

		}
		posts = append(posts, &post)
	}

	if err := rows.Err(); err != nil {
		//Erreur non critique : erreur lors du traitements des lignes
		lib.ErrorServer(w, "Error processing posts. Please try again later")
		return
	}
	data := lib.DataTest(w, r)
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
	var rows *sql.Rows
	var err_post error

	// Making the database request
	if !nothing {
		state_dislike := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts WHERE Dislike BETWEEN ? AND ?`
		rows, err_post = db.Query(state_dislike, mindislike, maxdislike)
	} else {
		state_dislike := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts`
		rows, err_post = db.Query(state_dislike)
	}

	if err_post != nil {
		//Erreur critique : échec de la récupération des posts par dislike
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving posts by dislike. Please try again later.",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
			//Erreur non critique : Echec de la lecture des données du posts
			lib.ErrorServer(w, "Error scanning post data. Please try again later.")
			return
		}
		if (post.Dislike > mindislike && post.Dislike < maxdislike) || nothing {
			posts = append(posts, &post)
		}
	}

	if err := rows.Err(); err != nil {
		//Erreur non critique: erreur lors du traitement des lignes
		lib.ErrorServer(w, "Error processing posts. Please try again later")
		return
	}
	data := lib.DataTest(w, r)
	data["Posts"] = posts
}

func CreationDate(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()
	// start := time.Now().Format("2006-01-02 15:04:05")

	// lastweek := now.AddDate(0, 0, -7)
	// lastmonth := now.AddDate(0, -1, 0)
	// lastyear := now.AddDate(-1, 0, 0)

	formValues := r.URL.Query()
	period := formValues.Get("Period")

	var posts []*models.Post
	var err_post error

	state_period := `SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID FROM Posts`
	rows, err_post := db.Query(state_period)

	if err_post != nil {
		//Erreur critique: echec de la récuperation des posts par date
		err := &models.CustomError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving posts by creation date. Please try again later",
		}
		HandleError(w, err.StatusCode, err.Message)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {
			//Erreur non critique : Echec de la lecture des données du posts
			lib.ErrorServer(w, "Error scanning post data. Please try again later")
			return
		}
		// Time := strings.Split(post.CreatedAt.Format(2024), " ")
		// Date := strings.Split(Time[0], "-")
		// Hours := strings.Split(Time[1], ":")

		now := time.Now()
		lastperiod := now
		nothing := false
		switch period {
		case "Last Week":
			lastperiod = now.AddDate(0, 0, -7)
		case "Last Month":
			lastperiod = now.AddDate(0, -1, 0)
		case "last Year":
			lastperiod = now.AddDate(-1, 0, 0)
		default:
			nothing = true
		}
		display := post.CreatedAt.Before(lastperiod)
		if !display || nothing {
			posts = append(posts, &post)
		}
	}

	if err := rows.Err(); err != nil {
		//Erreur non critique : erreur lors du traitement des lignes
		lib.ErrorServer(w, "Error processing posts. Please try again later.")
		return
	}
}
