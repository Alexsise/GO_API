package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Mod struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

var modCollection  = []Mod {
	{"1", "Bite", "+330% Critical Chance\n+220% Critical Damage"},
	{"2", "Ulfrun's Endurance", "Ulfrun's Descent Augment: During Ulfrun’s attack, enemies that die from Slash Status within 20m restore Voruna’s charges."},
	{"3", "Fracturing Crush", "Crush Augment: Crush gains +50% casting speed. The armor of surviving enemies decreases by 75% and they are unable to move for 7s."},
}

func main () {
	router := gin.Default()

	router.GET("/mods", getMods)
	router.GET("/mods/:id", getModByID)
    router.POST("/mods", createMod)
    router.PUT("/mods/:id", updateMod)
    router.DELETE("/mods/:id", deleteMod)
	router.Run(":8080")
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
