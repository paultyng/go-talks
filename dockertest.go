package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest"
)

var db *sql.DB

func main() {
	// startslide1 OMIT
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=secret"}) // HL
	if err != nil {
		log.Fatalf("Could not start resource: %// startslide1 OMITs", err)
	}
	// endslide1 OMIT

	// startslide2 OMIT
	// exponential backoff-retry, because the application in the container might
	// not be ready to accept connections yet
	if err := pool.Retry(func() error { // HL
		var err error
		// startslide3 OMIT
		db, err = sql.Open(
			"mysql",
			fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp")), // HL
		)
		// endslide3 OMIT
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// endslide2 OMIT

	// startslide4 OMIT
	_, err = db.Query("select 1")
	if err != nil {
		log.Fatalf("Could not query db: %s", err)
	}

	log.Print("Successfully queried DB!")
	// endslide4 OMIT

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
