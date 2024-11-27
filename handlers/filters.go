package handlers

import (
	"database/sql"
	"fmt"
	"forum/cmd/lib"
	"forum/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func FiltersHandler(w http.ResponseWriter, r *http.Request) {
	db := lib.GetDB()

	formValues := r.URL.Query()

	categories := formValues.Get("Category")
	numberlike := formValues.Get("Like")
	numberdislike := formValues.Get("Dislike")
	period := formValues.Get("Period")

	fmt.Println(period)
	fmt.Println(numberdislike)
	fmt.Println(numberlike)

	var posts []*models.Post
	var err_post error
	// sqlite
	state_filters :=
		`SELECT ID, Category_ID, Title, Text, Like, Dislike, CreatedAt, User_UUID, ImagePath
	FROM Posts
	WHERE 
   (
    CASE 
      WHEN ? = 'tous les likes' THEN Like >= 0
      WHEN ? = '1-10' THEN Like BETWEEN 1 AND 10
      WHEN ? = '11-50' THEN Like BETWEEN 11 AND 50
      WHEN ? = '51-100' THEN Like BETWEEN 51 AND 100
      WHEN ? = 'plus-de-100' THEN Like >= 101
      ELSE True
    END
  )
	AND
	(
    CASE 
      WHEN ? = 'tous les dislikes' THEN Dislike >= 0
      WHEN ? = '1-10' THEN Dislike BETWEEN 1 AND 10
      WHEN ? = '11-50' THEN Dislike BETWEEN 11 AND 50
      WHEN ? = '51-100' THEN Dislike BETWEEN 51 AND 100
      WHEN ? = 'plus-de-100' THEN Dislike >= 101
      ELSE True
    END
  )
	AND
	(
	CASE
	    WHEN ? = 'today' THEN date(CreatedAt) = date('now')
        WHEN ? = 'last-7-days' THEN date(CreatedAt) BETWEEN date('now', '-7 days') AND date('now')
		WHEN ? = 'last-30-days' THEN date(CreatedAt) BETWEEN date ('now','-30 days') AND date('now')
		WHEN ? = 'last-7-days' THEN date(CreatedAt) BETWEEN date('now', '-365 days') AND date('now')
		ELSE True 
		END
)
	 `
	var likeTous, like1_10, like11_50, like51_100, likePlus100 string
	switch numberlike {
	case "tous les likes":
		likeTous = numberlike
	case "1-10":
		like1_10 = numberlike
	case "11-50":
		like11_50 = numberlike
	case "51-100":
		like51_100 = numberlike
	case "plus-de-100":
		likePlus100 = numberlike
	}

	var dislikeTous, dislike1_10, dislike11_50, dislike51_100, dislikePlus100 string
	switch numberdislike {
	case "tous les likes":
		dislikeTous = numberdislike
	case "1-10":
		dislike1_10 = numberdislike
	case "11-50":
		dislike11_50 = numberdislike
	case "51-100":
		dislike51_100 = numberdislike
	case "plus-de-100":
		dislikePlus100 = numberdislike
	}

	rows, err_post := db.Query(state_filters, categories, categories, likeTous, like1_10, like11_50, like51_100, likePlus100, dislikeTous, dislike1_10, dislike11_50, dislike51_100, dislikePlus100, period, period, period, period)

	if err_post != nil {
		lib.ErrorMessage(w, nil, err_post.Error())
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			lib.ErrorMessage(w, nil, err.Error())
		}
	}()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Category_ID, &post.Title, &post.Text, &post.Like, &post.Dislike, &post.CreatedAt, &post.User_UUID, &post.ImagePath); err != nil {
		}

		category_array := strings.Split(post.Category_ID, " - ")
		is_contained := false
		for i := 0; i < len(category_array); i++ {
			if category_array[i] == categories {
				is_contained = true
			}
		}

		if is_contained {
			posts = append(posts, &post)
		}
	}

	if err := rows.Err(); err != nil {
		lib.ErrorMessage(w, nil, err.Error())
		return
	}

	data := lib.DataTest(w, r)
	data["Posts"] = posts
	data = lib.ErrorMessage(w, data, "none")

	lib.RenderTemplate(w, "layout/index", "page/index", data)

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
