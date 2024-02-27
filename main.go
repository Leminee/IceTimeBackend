package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Break struct {
	Id   int    `json:"id"`
	Time string `json:"time"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "ice:break@tcp(localhost:3306)/break")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/break", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, time FROM break.ice_break")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var breaks []Break

		for rows.Next() {
			var brk Break
			err := rows.Scan(&brk.Id, &brk.Time)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			breaks = append(breaks, brk)
		}

		c.JSON(http.StatusOK, breaks)
	})

	r.POST("/break", func(c *gin.Context) {
		var brk Break
		if err := c.BindJSON(&brk); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO break.ice_break (time, added_at) VALUES (?, NOW())", brk.Time)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, brk)
	})

	r.Run(":8080")
}
