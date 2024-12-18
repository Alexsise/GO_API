package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Mod struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var modCollection = []Mod{
	{"1", "Bite", "+330% Critical Chance\n+220% Critical Damage"},
	{"2", "Ulfrun's Endurance", "Ulfrun's Descent Augment: During Ulfrun’s attack, enemies that die from Slash Status within 20m restore Voruna’s charges."},
	{"3", "Fracturing Crush", "Crush Augment: Crush gains +50% casting speed. The armor of surviving enemies decreases by 75% and they are unable to move for 7s."},
}

var jwtKey = []byte("i_love_cabbage")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

var userStorage = map[string]string{
	"admin": "adminpass",
	"user": "userpass",
}

var userRoles = map[string]string{
	"admin": "admin",
	"user": "user",
}

func main() {
	router := gin.Default()

	router.POST("/login", login)
	router.POST("/register", register)

	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/mods", getMods)
		protected.GET("/mods/:id", getModByID)
		protected.POST("/mods", adminMiddleware(), createMod)
		protected.PUT("/mods/:id", adminMiddleware(), updateMod)
		protected.DELETE("/mods/:id", adminMiddleware(), deleteMod)
	}

	router.Run(":8080")
}

func register(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	if _, exists := userStorage[creds.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"message": "user already exists"})
		return
	}

	userStorage[creds.Username] = creds.Password
	userRoles[creds.Username] = "user" // Default role is "user"
	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func login(c *gin.Context) {
	var creds Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	storedPassword, exists := userStorage[creds.Username]
	if !exists || storedPassword != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	role := userRoles[creds.Username]
	token, err := generateToken(creds.Username, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func generateToken(username, role string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden: admin access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func getMods(c *gin.Context) {
	c.JSON(http.StatusOK, modCollection)
}

func getModByID(c *gin.Context) {
	id := c.Param("id")

	for _, mod := range modCollection {
		if mod.ID == id {
			c.JSON(http.StatusOK, mod)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "mod not found"})
}

func createMod(c *gin.Context) {
	var newMod Mod

	if err := c.BindJSON(&newMod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	modCollection = append(modCollection, newMod)
	c.JSON(http.StatusCreated, newMod)
}

func updateMod(c *gin.Context) {
	id := c.Param("id")
	var updatedMod Mod

	if err := c.BindJSON(&updatedMod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	for i, mod := range modCollection {
		if mod.ID == id {
			modCollection[i] = updatedMod
			c.JSON(http.StatusOK, updatedMod)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "mod not found"})
}

func deleteMod(c *gin.Context) {
	id := c.Param("id")

	for i, mod := range modCollection {
		if mod.ID == id {
			modCollection = append(modCollection[:i], modCollection[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "mod deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "mod not found"})
}
