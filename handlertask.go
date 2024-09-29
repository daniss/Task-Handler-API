package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"fmt"
)

func createTask(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/createTask" ,func(c *gin.Context) {
		var task Task

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := db.Create(&task)
		if result.RowsAffected == 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "Task didn't create because of a not unique name"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Task created successfully"})
		fmt.Println(task)
		
	})
}

func getTasks(r *gin.RouterGroup, db *gorm.DB) {
    r.GET("/getTasks", func(c *gin.Context) {
		var task []Task

		result := db.Find(&task)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusOK, task)
	})

	r.GET("/getTask/:id", func(c *gin.Context) {
		var task Task
		id := c.Param("id")

		result := db.First(&task, id)

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"Updated task": task})
	})
}

func updateTask(r *gin.RouterGroup, db *gorm.DB) {
    r.PUT("/updateTask/:id",func(c *gin.Context) {
		var task Task

		result := db.First(&task, c.Param("id"))

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db.Save(&task)
		c.JSON(http.StatusOK, task)
	})
}

func deleteTask(r *gin.RouterGroup, db *gorm.DB) {
    r.DELETE("/deleteTask/:id",func(c *gin.Context) {
		var task Task

		result := db.First(&task, c.Param("id"))

		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
			return
		}

		db.Delete(&task)
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})

		db.Exec("UPDATE tasks SET id = id - 1 WHERE id > ?", c.Param("id"))
	})
}

func main() {
	
	db := setupDatabase()
	r := gin.Default()

	api := r.Group("/api")
	api.Use(jwtAuthMiddleware)
	register(r, db)
	login(r, db)
	createTask(api, db)
	getTasks(api, db)
	updateTask(api, db)
	deleteTask(api, db)
	
	
	r.Run()
	
}
