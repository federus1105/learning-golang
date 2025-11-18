package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// create user
func createUser(c *gin.Context) {
	var user struct {
		Name       string `json:"name"`
		Department string `json:"department"`
		Email      string `json:"email"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Printf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID int
	err := db.QueryRow(`INSERT INTO users (name, department, email) 
                        VALUES ($1, $2, $3) RETURNING id`,
		user.Name, user.Department, user.Email).Scan(&userID)

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": userID})
}

// update user
func updateUser(c *gin.Context) {
	id := c.Param("id")
	var user struct {
		Name       string `json:"name"`
		Department string `json:"department"`
		Email      string `json:"email"`
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`UPDATE users 
                       SET name=$1, department=$2, email=$3 
                       WHERE id=$4`,
		user.Name, user.Department, user.Email, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated"})
}

// delete user
func deleteUser(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// fetch users
func getUsersFromDB() ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT id, name, department, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var name, department, email string
		err = rows.Scan(&id, &name, &department, &email)
		if err != nil {
			return nil, err
		}

		users = append(users, gin.H{
			"id":         id,
			"name":       name,
			"department": department,
			"email":      email,
		})
	}

	return users, nil
}
