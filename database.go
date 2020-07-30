package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	// "html/template"

	_ "github.com/lib/pq"
)

type SiteStruct struct {
	Site string `json: "site"`
}

type Date struct {
	Year  int
	Month int
	Day   int
}

type Post struct {
	RowId          int
	PostId         int
	Rating         int
	Source         string
	Category       string
	AgridataSource string
	Website        string
	Nsfw           bool
	Misc           string
	AddDate        string
	PublishDate    string
}

type QueryStruct struct {
	Mood string `json: "mood"`
	Site string `json: "site"`
	Upto string `json: "upto"`
}

type ResultsQueryStruct struct {
	Mood string `json: "mood"`
	Site string `json: "site"`
	Num  int    `json: "num"`
}

type MatchStruct struct {
	RowId1 int     `json: "rowid1"`
	RowId2 int     `json: "rowid2"`
	Score1 float64 `json: "score1"`
	Score2 float64 `json: "score2"`
}

var db *sql.DB

// NOT UPDATED, DONT USE, WILL BREAK
func printEntry(rowId int) {
	// Pretty print a single entry from the database

	sqlStatement := `
	SELECT * FROM posts
	WHERE row_id=$1`

	var postId, rating int
	var source, category, agridataSource, website string
	var nsfw bool

	// Replace 3 with an ID from your database or another random
	// value to test the no rows use case.
	row := db.QueryRow(sqlStatement, rowId)
	switch err := row.Scan(&rowId, &postId, &rating, &source, &category,
		&agridataSource, &website, &nsfw); err {
	case sql.ErrNoRows:
		fmt.Println("No row with that id")
	case nil:
		fmt.Println(rowId, postId, rating, source, category, agridataSource,
			website, nsfw)
	default:
		panic(err)
	}
}

// //NOT UPDATED, DONT USE, WILL BREAK
// func printPost(postId int) {
// 	// Pretty print all entries corresponding to a post from the database
//
// 	sqlStatement := `
// 	SELECT * FROM posts
// 	WHERE post_id=$1`
//
// 	var rowId, rating int
// 	var source, category, agridataSource, website string
// 	var nsfw bool
//
// 	rows, err := db.Query(sqlStatement, postId)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	defer rows.Close()
// 	for rows.Next() {
// 		err = rows.Scan(&rowId, &postId, &rating, &source, &category,
// 			&agridataSource, &website, &nsfw)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fmt.Println(rowId, postId, rating, source, category, agridataSource,
// 			website, nsfw)
// 	}
//
// 	err = rows.Err()
// 	if err != nil {
// 		panic(err)
// 	}
// }

//NOT UPDATED, DONT USE, WILL BREAK
func printAll() {
	// pretty print all posts in the database

	fmt.Println("Database: ")

	sqlStatement := `SELECT * FROM posts`

	var rowId, postId int
	var rating float64
	var source, category, agridataSource, website, addDate, publishDate string
	var nsfw bool

	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&rowId, &postId, &rating, &source, &category,
			&agridataSource, &website, &nsfw, &addDate, &publishDate)
		if err != nil {
			panic(err)
		}
		fmt.Println(rowId, postId, rating, source, category, agridataSource,
			website, nsfw)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}
}

func addPost(initialRating int, source string, agridataSource string,
	website string, categories []string, nsfw bool, misc, addDate,
	publishDate string) int {
	// Add a post to the database, returns the new post's post_id (common to all
	// entries of this post)
	sqlStatement := `
	SELECT rating FROM posts WHERE agridata_source=$1
	`
	var rating float64

	err := db.QueryRow(sqlStatement, agridataSource).Scan(&rating)
	if err == sql.ErrNoRows {
		sqlStatement = `
		INSERT INTO posts (rating, source, category, agridata_source, 
		website, nsfw, misc, add_date, publish_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING row_id, post_id`
		rowId, postId := 0, 0
		err := db.QueryRow(sqlStatement, initialRating, source, "none",
			agridataSource, website, nsfw, misc, addDate, publishDate).Scan(&rowId,
			&postId)
		if err != nil {
			panic(err)
		}

		// for idx, category := range categories {
		// 	if idx != 0 {
		// 		sqlStatement = `
		// 		INSERT INTO posts (post_id, rating, source, category,
		// 						agridata_source, website, nsfw, misc)
		// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		// 		RETURNING row_id`
		// 		err = db.QueryRow(sqlStatement, postId, initialRating, source,
		// 			category, agridataSource, website, nsfw, misc).Scan(&rowId)
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 	}
		// }
		return postId
	} else {
		log.Println("Post you are trying to add is already in the database!")
		return -1
	}
}

func removeEntry(rowId int) {
	// Remove a post by post_id or entry by row_id from the database

	sqlStatement := `
	DELETE FROM posts
	WHERE row_id=$1`

	_, err := db.Exec(sqlStatement, rowId)
	if err != nil {
		panic(err)
	}
}

func removePost(postId int) {
	// Remove a post by post_id or entry by row_id from the database

	sqlStatement := `
	DELETE FROM posts
	WHERE post_id=$1`

	_, err := db.Exec(sqlStatement, postId)
	if err != nil {
		panic(err)
	}
}

func recordMatch(id1 int, id2 int, score1 float64, score2 float64) {
	// Given two row id's and the real scores from the match, updates their ELO
	// ratings
	sqlStatement := `SELECT rating FROM posts WHERE row_id=$1`
	var rating1 float64
	err := db.QueryRow(sqlStatement, id1).Scan(&rating1)
	if err != nil {
		panic(err)
	}

	var rating2 float64
	err = db.QueryRow(sqlStatement, id2).Scan(&rating2)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Rating 1 Before: %f Rating 2 Before: %f", rating1, rating2)

	expectation1 := 1 / (1 + math.Pow(10, (rating2-rating1)/400))
	expectation2 := 1 / (1 + math.Pow(10, (rating1-rating2)/400))

	newRating1 := rating1 + 20*(score1-expectation1)
	newRating2 := rating2 + 20*(score2-expectation2)

	sqlStatement = `
	UPDATE posts
	SET rating=$2
	WHERE row_id=$1`

	_, err = db.Exec(sqlStatement, id1, newRating1)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(sqlStatement, id2, newRating2)
	if err != nil {
		panic(err)
	}
}

/*
Ping with no data to delete the entire posts table
*/
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	deleteOrder := "DROP TABLE posts;"
	_, err := db.Exec(deleteOrder)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Cleared")
	printAll()
	w.WriteHeader(http.StatusOK)
}

/*
INPUT JSON {
	'rowid1': id of first post
	'rowid2': id of second post
	'score1': score of first post
	'score2': score of second post
}
OUTPUT NONE
*/
// call this each time the user makes a decision on a comparison to update elo
// scores
func matchHandler(w http.ResponseWriter, r *http.Request) {
	result := &MatchStruct{}
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	recordMatch(result.RowId1, result.RowId2, result.Score1, result.Score2)

	fmt.Println("Match successfully recorded")
	w.WriteHeader(http.StatusOK)
}

/*
INPUT JSON {
	'category': happy, sad, etc, category you want the post pulled from
}
OUTPUT JSON {
	'row_id1': first row_id,
	'row_id2': second row_id,
	'post_id1': first post_id
	'post_id2': second post-id,
	'rating1': first rating,
	'rating2': second rating,
	'source1': first source,
	'source2': second source,
	'agridata_source1': first agridata_source,
	'agridata_source2': second agridata_source,
	'website1': first website,
	'website2': second website
	'misc1': first miscellaneous data
	'misc2': second miscellaneous data
	'add_date1': date 1st post was added to db
	'add_date2': date 2nd post was added to db
	'publish_date1': date 1st post was published
	'publish_date2': date 2nd post was published
}
*/
// call this each time next is clicked to get comparison ID/ rows for each
// comparison (returns 1 tuple of posts)
func postsHandler(w http.ResponseWriter, r *http.Request) {
	result := &QueryStruct{}
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		panic(err)
	}
	var mood, site string
	mood = result.Mood
	site = result.Site
	sqlQuery := ""
	var rows *sql.Rows

	if mood == "none" {
		sqlQuery = `SELECT * FROM posts
		WHERE website=$1
		ORDER BY RANDOM()
		LIMIT 2
		`
		rows, err = db.Query(sqlQuery, site)
		if err != nil {
			panic(err)
		}
	} else {

		sqlQuery = `SELECT * FROM posts
		WHERE category=$1 AND website=$2
		ORDER BY RANDOM()
		LIMIT 2
		`
		rows, err = db.Query(sqlQuery, mood, site)
		if err != nil {
			panic(err)
		}
	}

	var rowId1, rowId2, postId1, postId2, rating1, rating2 int
	var source1, source2, agridataSource1, agridataSource2, website1,
		website2, misc1, misc2, addDate1, addDate2, publishDate1,
		publishDate2 string
	var nsfw1, nsfw2 bool

	row_num := 1
	defer rows.Close()
	for rows.Next() {
		if row_num == 1 {
			err = rows.Scan(&rowId1, &postId1, &rating1, &source1, &mood,
				&agridataSource1, &website1, &nsfw1, &misc1, &addDate1, &publishDate1)
		} else {
			err = rows.Scan(&rowId2, &postId2, &rating2, &source2, &mood,
				&agridataSource2, &website2, &nsfw2, &misc2, &addDate2,
				&publishDate2)
		}
		if err != nil {
			fmt.Println("RIGHT HERE")
			panic(err)
		}
		row_num += 1
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"row_id1":          rowId1,
		"row_id2":          rowId2,
		"post_id1":         postId1,
		"post_id2":         postId2,
		"rating1":          rating1,
		"rating2":          rating2,
		"source1":          source1,
		"source2":          source2,
		"agridata_source1": agridataSource1,
		"agridata_source2": agridataSource2,
		"website1":         website1,
		"website2":         website2,
		"nsfw1":            nsfw1,
		"nsfw2":            nsfw2,
		"misc1":            misc1,
		"misc2":            misc2,
		"add_date1":        addDate1,
		"add_date2":        addDate2,
		"publish_date1":    publishDate1,
		"publish_date2":    publishDate2})
	if err != nil {
		panic(err)
	}
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	result := &ResultsQueryStruct{}
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		panic(err)
	}
	var mood, site string
	mood = result.Mood
	site = result.Site
	// num := result.Num
	var rows *sql.Rows

	if mood == "none" {
		sqlQuery := `SELECT * FROM posts
		WHERE website=$1
		ORDER BY rating DESC
		`
		rows, err = db.Query(sqlQuery, site)
		if err != nil {
			panic(err)
		}
	} else {

		sqlQuery := `SELECT * FROM posts
		WHERE category=$1 AND website=$2
		ORDER BY rating DESC
		`
		rows, err = db.Query(sqlQuery, mood, site)
		if err != nil {
			panic(err)
		}
	}

	var rowId, postId int
	var rating float64
	var source, agridataSource, website, misc, addDate, publishDate string
	var nsfw bool
	var rowIds, postIds []int
	var ratings []float64
	var sources, agridataSources, websites, miscs, addDates,
		publishDates []string
	var nsfws []bool

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&rowId, &postId, &rating, &source, &mood,
			&agridataSource, &website, &nsfw, &misc, &addDate, &publishDate)
		if err != nil {
			panic(err)
		}
		rowIds = append(rowIds, rowId)
		postIds = append(postIds, postId)
		ratings = append(ratings, rating)
		sources = append(sources, source)
		agridataSources = append(agridataSources, agridataSource)
		websites = append(websites, website)
		nsfws = append(nsfws, nsfw)
		miscs = append(miscs, misc)
		addDates = append(addDates, addDate)
		publishDates = append(publishDates, publishDate)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"row_ids":          rowIds,
		"post_ids":         postIds,
		"ratings":          ratings,
		"sources":          sources,
		"agridata_sources": agridataSources,
		"websites":         websites,
		"nsfws":            nsfws,
		"miscs":            miscs,
		"add_dates":        addDates,
		"publish_dates":    publishDates})
	if err != nil {
		panic(err)
	}
}
