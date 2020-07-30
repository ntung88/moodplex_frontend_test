package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var repeatFillDB = true

const refillInterval = 5
const dbLimit = 10000

func initDb(db *sql.DB) {
	//Initializes posts table given the pointer to a database
	createMoodType := `DO $$ BEGIN
		CREATE TYPE mood AS ENUM ('happy', 'funny', 'informative', 
			'motivational', 'sad', 'cute', 'educational', 'angry', 'uplifting', 
		    'scary', 'artistic', 'news', 'romantic', 'none');
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$`
	createWebsiteType := `DO $$ BEGIN
		CREATE TYPE website AS ENUM ('twitter', 'reddit', 'youtube', 'imgur',
			'hackernews');
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$`
	createOrder := `CREATE TABLE IF NOT EXISTS posts (
	row_id SERIAL UNIQUE,
	post_id SERIAL,
	rating FLOAT,
	source TEXT,
	category mood,
	agridata_source TEXT,
	website website,
	nsfw boolean,
	misc TEXT,
	add_date timestamp with time zone,
	publish_date timestamp with time zone
	);`
	_, err := db.Exec(createMoodType)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(createWebsiteType)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(createOrder)
	if err != nil {
		panic(err)
	}
}

func wipeAndFillDB() {
	deleteTable := "DROP TABLE IF EXISTS posts;"
	_, err := db.Exec(deleteTable)
	if err != nil {
		panic(err)
	}
	deleteMoodEnum := "DROP TYPE IF EXISTS mood;"
	_, err = db.Exec(deleteMoodEnum)
	if err != nil {
		panic(err)
	}
	deleteWebsiteEnum := "DROP TYPE IF EXISTS website;"
	_, err = db.Exec(deleteWebsiteEnum)
	if err != nil {
		panic(err)
	}

	initDb(db)

	if err != nil {
		panic("Cannot initialize database")
	}

	fmt.Println("Successfully Connected")
	for _, src := range dataSources {
		src.fillDB()
	}
	log.Println("Database filled!")
	refillDB()
}

func refillDB() {
	ticker := time.NewTicker(refillInterval * time.Minute)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			var ct int
			ctRow := db.QueryRow("SELECT COUNT(*) from posts;")
			err := ctRow.Scan(&ct)
			if err != nil {
				panic("Cannot determine how many rows in the database")
			}
			if repeatFillDB {
				if ct >= dbLimit {
					wipeAndFillDB()
				} else {
					// Hacker News
					hackernewsSrc.fillDB()
					log.Println("Successfully refilled database!")
				}
			} else {
				close(quit)
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func main() {
	// templates := template.Must(template.ParseFiles("frontend/index.html"))

	// Handle for css for html
	// http.Handle("/static/", http.StripPrefix("/static/",
	// http.FileServer(http.Dir("static"))))

	// Connects to the database

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("Cannot open the database")
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic("Cannot connect to the database")
	}

	// go wipeAndFillDB()

	r := mux.NewRouter()
	r.HandleFunc("/match", matchHandler).Methods("POST")
	r.HandleFunc("/posts", postsHandler).Methods("POST")
	r.HandleFunc("/delete", deleteHandler).Methods("POST")
	r.HandleFunc("/getresults", resultsHandler).Methods("POST")

	r.HandleFunc("/source/{site}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data := SiteStruct{Site: vars["site"]}
		tmpl := template.Must(template.ParseFiles("frontend/redirect.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/category/{mood}/{site}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data := QueryStruct{
			Mood: vars["mood"],
			Site: vars["site"]}
		tmpl := template.Must(template.ParseFiles("frontend/category.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/results/{mood}/{site}/{upto}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data := QueryStruct{
			Mood: vars["mood"],
			Site: vars["site"],
			Upto: vars["upto"]}
		tmpl := template.Must(template.ParseFiles("frontend/results.html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	r.HandleFunc("/secretresults/{mood}/{site}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		data := QueryStruct{
			Mood: vars["mood"],
			Site: vars["site"]}
		tmpl := template.Must(template.ParseFiles("frontend/results_" + vars["site"] + ".html"))
		err := tmpl.Execute(w, data)
		if err != nil {
			panic(err)
		}
	})

	staticDir := "/frontend"

	// Create the route
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("." + staticDir)))

	fmt.Println("Listening on : " + os.Getenv("PORT") + " ...")
	err = http.ListenAndServe(":"+os.Getenv("PORT"),
		handlers.CORS(
			handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type",
				"Authorization"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD",
				"OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}))(r))
	if err != nil {
		panic(err)
	}
}
