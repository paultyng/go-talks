package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
)

var db *sql.DB

func main() {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might
	// not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err = sql.Open(
			"mysql",
			fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")),
		)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	//startslide OMIT
	// setup a test server and http.Handler that queries the DB
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// select now() as a string
		res, err := db.Query("select cast(now() as char)") // HL
		if err != nil {
			panic(err)
		}
		defer res.Close()

		// grab the string value from the result set and write it to the response
		if res.Next() {
			var s string
			res.Scan(&s) // HL
			_, err := w.Write(([]byte)(s))
			if err != nil {
				panic(err)
			}
		}
	}))
	//endslide OMIT
	defer ts.Close()

	// execute an http get against our new server
	res, err := http.Get(ts.URL)
	if err != nil {
		log.Fatal(err)
	}
	now, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", now)

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
