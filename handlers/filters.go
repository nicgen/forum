package handlers

import (
	"database/sql"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strconv"
	"time"
)

// const (
// 	likeMin1    = 1
// 	likeMax1    = 10
// 	likeMin2    = 11
// 	likeMax2    = 50
// 	likeMin3    = 51
// 	likeMax3    = 100
// 	dislikeMin1 = 1
// 	dislikeMax1 = 10
// 	dislikeMin2 = 11
// 	dislikeMax2 = 50
// 	dislikeMin3 = 51
// 	dislikeMax3 = 100
// )

func FiltersHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()

	formValues := r.URL.Query()

	categories := formValues.Get("Category")
	numberlike := formValues.Get("Like")
	numberdislike := formValues.Get("Dislike")
	period := formValues.Get("Period")

	println("Like: ", numberlike)
	println("Dislike: ", numberdislike)
	println("Period: ", period)

	// Getting the number of values stored in the Categories tables
	// var count int
	// err := db.QueryRow("SELECT COUNT(*) FROM Categories").Scan(&count)
	// if err != nil {
	// 	// Handle error
	// }

	// // filter category
	// category := ""
	// for i := 0; i < count; i++ {
	// 	compteur := strconv.Itoa(i)
	// 	notempty := formValues.Get("Category" + compteur)
	// 	if len(notempty) != 0 {
	// 		category = notempty
	// 	}
	// }

	// Filter like
	// likemin, likemax := 0, 0
	// switch numberlike {
	// case "tous les likes":
	// 	likemin = 0
	// 	likemax = 0
	// case "1-10":
	// 	likemin = likeMin1
	// 	likemax = likeMax1
	// case "11-50":
	// 	likemin = likeMin2
	// 	likemax = likeMax2
	// case "51-100":
	// 	likemin = likeMin3
	// 	likemax = likeMax3
	// case "plus-de-100":
	// 	likemin = 101
	// }

	// filter dislike
	// dislikemin, dislikemax := 0, 0
	// switch numberdislike {
	// case "tous les dislikes":
	// 	dislikemin = 0
	// 	dislikemax = 0
	// case "1-10":
	// 	dislikemin = dislikeMin1
	// 	dislikemax = dislikeMax1
	// case "11-50":
	// 	dislikemin = dislikeMin2
	// 	dislikemax = dislikeMax2
	// case "51-100":
	// 	dislikemin = dislikeMin3
	// 	dislikemax = dislikeMax3
	// case "plus-de-100":
	// 	dislikemin = 101
	// }

	//filter period
	// lastperiod := ""
	// switch period {
	// case "week":
	// 	lastperiod = "-7 days"
	// case "month":
	// 	lastperiod = "start of month"
	// case "year":
	// 	lastperiod = "start of year"
	// }
	var posts []*models.Post
	var err_post error

	// println(likemin)
	// println(likemax)
	// println(dislikemin)
	// println(dislikemax)
	// println(categories)
	// println(period)

	state_filters :=
		`SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID
	FROM Posts
	WHERE 
  (Category_ID = ?)
	 `
	rows, err_post := db.Query(state_filters, categories)

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
	data := lib.DataTest(w, r)
	data["Posts"] = posts
	data = lib.ErrorMessage(w, data, "none")

	lib.RenderTemplate(w, "layout/index", "page/index", data)
	// Redirect User to the home page
	// http.Redirect(w, r, "/", http.StatusSeeOther)
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

	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID); err != nil {

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
	}
}
