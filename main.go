package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type CreatPost struct {
	ID   int    `db:"id" json:"id"`
	Name string `json:"name" binding:"required"`
	POST string `json:"owner" binding:"required"`
}

func main() {

	// connect to mysql
	db, err := sql.Open("mysql", "root:@Aa1234567890@tcp(127.0.0.1:3306)/blog")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// check if connection is ok
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Connected to the database")

	r := gin.Default()
	r.POST("/api/post", func(c *gin.Context) {
		var NewPost CreatPost
		c.BindJSON(&NewPost)
		c.JSON(http.StatusCreated, NewPost)

		fmt.Print(NewPost)

		query := "INSERT INTO newpost (`name`, `owner`) VALUES (? , ?)"
		_, err := db.ExecContext(context.Background(), query, NewPost.Name, NewPost.POST)

		//erro, _ := db.QueryRow(query, NewPost.Name, NewPost.POST), time.Now()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "failed to create text"})
			fmt.Println(err)
			return
		}
	})

	r.GET("/api/post/", func(c *gin.Context) {

		rows, err := db.Query("SELECT * FROM newpost")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve jobs"})
			return
		}
		defer rows.Close()
		var Posts []CreatPost

		for rows.Next() {
			var Post CreatPost
			err := rows.Scan(&Post.ID, &Post.Name, &Post.POST)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve jobs"})
				return
			}
			Posts = append(Posts, Post)
		}
		c.JSON(http.StatusOK, Posts)

	})

	r.GET("/api/post/:id", func(c *gin.Context) {
		id := c.Param("id")
		rows, err := db.Query("SELECT * FROM newpost WHERE id =?", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve jobs"})
			return
		}
		defer rows.Close()
		var Posts []CreatPost

		for rows.Next() {
			var Post CreatPost
			err := rows.Scan(&Post.ID, &Post.Name, &Post.POST)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve jobs"})
				return
			}
			Posts = append(Posts, Post)
		}
		c.JSON(http.StatusOK, Posts)

	})

	r.PUT("/api/post/:id", func(c *gin.Context) {
		id := c.Param("id")
		//var NewPost CreatPost
		var NewPost CreatPost

		//NewPost.ID = id
		c.BindJSON(&NewPost)
		//c.JSON(http.StatusUpgradeRequired, NewPost)

		fmt.Print(NewPost)

		_, err := db.Exec("UPDATE newpost  set  name=?  , owner=?  where id= ?", NewPost.Name, NewPost.POST, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create text"})
			fmt.Println(err)
			return

			c.String(http.StatusOK, " %s", id)

		}
	})

	r.DELETE("/api/post/:id", func(c *gin.Context) {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM  newpost where id= ?", id)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create text"})
			fmt.Println(err)
			return
		}
		c.String(http.StatusOK, " %s", id)
	})

	r.Run(":8081") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
