package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
	"github.com/tasks/GolangPythonMicroservice/Go/models"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func PokesRouter() {

	// Execution mode. Exotic condition to ensure ONLY DEBUG/RELEASE values are accepted
	if os.Getenv("POK_ENV_DEPLOY_MODE") == "DEBUG" || os.Getenv("POK_ENV_DEPLOY_MODE") == "RELEASE" {
		if os.Getenv("POK_ENV_DEPLOY_MODE") == "RELEASE" {
			gin.SetMode(gin.ReleaseMode)
		}
	} else {
		fmt.Println("*ERROR Posts Microservice Router: invalid environment 'DEPLOY_MODE' aborting!")
		return
	}
	// Router config
	router := gin.Default()

	// CORS config
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE ,OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	// Routes/Endpoints
	pokemon := router.Group("api")
	{
		pokemon.GET("/pokemon", handlerGetAllPokemon)
		pokemon.GET("/pokemon/info/:id", handlerGetPokemon)
		pokemon.GET("/pokemon/exist/id/:id", hanlderExistPokemon)
		pokemon.POST("/pokemon/register", handlerPostPokemon)
		pokemon.PUT("/pokemon/update/:id", handlerUpdatePokemon)
		pokemon.DELETE("/pokemon/delete/:id", handlerDeletePokemon)

	}

	router.Run(":" + os.Getenv("POK_ENV_API_PORT"))
}
func handlerGetAllPokemon(c *gin.Context) {

	var (
		response gin.H
		pokemons []models.ObjPokemonGet
		err      error
	)

	pokemons, err = models.FunGetAllPokemon()
	if err != nil {
		response = gin.H{"message": "error quering pokes", "response": nil, "error": "RE"}
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	if pokemons == nil {
		response = gin.H{"message": "pokes not found", "response": nil, "error": "RE"}
		c.JSON(http.StatusNotFound, response)
		return
	}

	response = gin.H{"message": "pokemons found", "error": nil, "status": http.StatusOK, "response": pokemons}
	c.JSON(http.StatusOK, response)

}

func handlerPostPokemon(c *gin.Context) {
	var (
		vginResponse gin.H
		vPokemon     models.ObjPokemonPost
		err          error
	)

	err = c.BindJSON(&vPokemon)
	if err != nil {
		err = errors.New("*Error handlerPostPokemon: couldn't bind payload provided whit ObjPokemonPost: " + err.Error())
		vginResponse = gin.H{"message": "error reading payload provided", "response": nil, "error": "RE", "status": http.StatusBadRequest}
		c.JSON(http.StatusBadRequest, vginResponse)
		return
	}

	vboolExists, err := models.FunExistPokemon("name", vPokemon.Name)
	vmapResponse := make(map[string]bool)

	if err != nil {
		err = errors.New("*ERROR handlerPostPokemon: couldn't know if pokemon exists by " + " -> " + err.Error())
		vginResponse = gin.H{"message": "internal error", "response": vmapResponse, "error": "IE", "status": http.StatusInternalServerError}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}
	if vboolExists {
		vmapResponse["key"] = false
		vginResponse = gin.H{"message": "pokemon name exist", "response": vmapResponse, "error": nil, "status": http.StatusOK}
		c.JSON(http.StatusOK, vginResponse)
		return
	}

	err = models.FunPostPokemon(vPokemon)
	if err != nil {
		err = errors.New("*Error FunPostPokemon couldn't insert pokemon in db")
		vginResponse = gin.H{"message": "internal error", "response": nil, "error": "IE", "status": http.StatusInternalServerError}
		c.JSON(http.StatusInternalServerError, vginResponse)
	}
	vginResponse = gin.H{"message": "Pokemon successfully registered", "error": nil, "status": http.StatusOK}
	c.JSON(http.StatusOK, vginResponse)
}

func handlerUpdatePokemon(c *gin.Context) {
	var (
		vPokemon     models.ObjPokemonPost
		vginResponse gin.H
		err          error
	)

	strID := c.Params.ByName("id")
	err = c.BindJSON(&vPokemon)
	if err != nil {
		err = errors.New("*Error handlerUpdatePokemon: couldn't bind payload provided whit ObjPokemonPost: " + err.Error())
		vginResponse = gin.H{"message": "error reading payload provided", "response": nil, "error": "RE", "status": http.StatusBadRequest}
		c.JSON(http.StatusBadRequest, vginResponse)
		return
	}

	vboolExists, err := models.FunExistPokemon("name", vPokemon.Name)
	vmapResponse := make(map[string]bool)

	if err != nil {
		err = errors.New("*ERROR handlerPostPokemon: couldn't know if pokemon exists by " + " -> " + err.Error())
		vginResponse = gin.H{"message": "internal error", "response": vmapResponse, "error": "IE", "status": http.StatusInternalServerError}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}
	if vboolExists {
		vmapResponse["key"] = false
		vginResponse = gin.H{"message": "pokemon name exist", "response": vmapResponse, "error": nil, "status": http.StatusOK}
		c.JSON(http.StatusOK, vginResponse)
		return
	}

	err = models.FunUpdatePokemon(vPokemon, strID)
	if err != nil {
		vginResponse = gin.H{"message": "error quering a poke", "response": nil, "error": "RE"}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}

	vginResponse = gin.H{"message": "Pokemon successfully updated", "error": nil, "status": http.StatusOK}
	c.JSON(http.StatusOK, vginResponse)

}

func handlerDeletePokemon(c *gin.Context) {
	var (
		vginResponse gin.H
		err          error
	)

	strID := c.Params.ByName("id")

	vboolExists, err := models.FunExistPokemon("id", strID)
	vmapResponse := make(map[string]bool)

	if err != nil {
		err = errors.New("*ERROR handlerDeletePokemon: couldn't know if pokemon exists by " + strID + " -> " + err.Error())
		vginResponse = gin.H{"message": "internal error", "response": vmapResponse, "error": "IE", "status": http.StatusInternalServerError}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}
	if !vboolExists {
		vmapResponse["key"] = false

		vginResponse = gin.H{"message": "ID pokemon doesn't exist", "response": vmapResponse, "error": nil, "status": http.StatusOK}
		c.JSON(http.StatusOK, vginResponse)
		return
	}

	err = models.FunDeletePokemon(strID)
	if err != nil {
		err = errors.New("*Error handlerDeletePokemon: couldn't bind payload provided whit ObjPokemonPost: " + err.Error())
		vginResponse = gin.H{"message": "error delete pokemon", "response": nil, "error": "RE"}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}
	vginResponse = gin.H{"message": "Pokemon successfully deleted", "error": nil, "status": http.StatusOK}
	c.JSON(http.StatusOK, vginResponse)

}

func handlerGetPokemon(c *gin.Context) {
	var (
		vPokemon     models.ObjPokemonGet
		vginResponse gin.H
		err          error
	)

	strID := c.Params.ByName("id")

	vboolExists, err := models.FunExistPokemon("id", strID)
	vmapResponse := make(map[string]bool)

	if err != nil {
		err = errors.New("*ERROR handlerGetPokemon: couldn't know if pokemon exists by " + strID + " -> " + err.Error())
		vginResponse = gin.H{"message": "internal error", "response": vmapResponse, "error": "IE", "status": http.StatusInternalServerError}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}
	if !vboolExists {
		vmapResponse["key"] = false

		vginResponse = gin.H{"message": "ID pokemon doesn't exist", "response": vmapResponse, "error": nil, "status": http.StatusOK}
		c.JSON(http.StatusOK, vginResponse)
		return
	}

	vPokemon, err = models.FunGetPokemon(strID)
	if err != nil {
		vginResponse = gin.H{"message": "error quering a poke", "response": nil, "error": "RE"}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}

	vginResponse = gin.H{"message": "pokemon found", "error": nil, "status": http.StatusOK, "response": vPokemon}
	c.JSON(http.StatusOK, vginResponse)

}

func hanlderExistPokemon(c *gin.Context) {

	var vginResponse gin.H
	vmapResponse := make(map[string]bool)

	strURL := strings.TrimPrefix(c.Request.RequestURI, "/api/pokemon/exist/")
	strField := strings.Split(strURL, "/")[0]
	strID := c.Params.ByName("id")

	vboolFlagExists, err := models.FunExistPokemon(strField, strID)
	if err != nil {
		err = errors.New("*ERROR hanlderExistPokemon: couldn't know if pokemon exists by " + strField + ": " + strID + " -> " + err.Error())
		vginResponse = gin.H{"message": "internal error", "response": vmapResponse, "error": "IE", "status": http.StatusInternalServerError}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}
	if !vboolFlagExists {
		vmapResponse["key"] = false
		vginResponse = gin.H{"message": "pokemon does not exist", "response": vmapResponse, "error": nil, "status": http.StatusOK}
		c.JSON(http.StatusOK, vginResponse)
		return
	}

	vmapResponse["key"] = true
	vginResponse = gin.H{"message": "ID pokemon exist", "response": vmapResponse, "error": nil, "status": http.StatusOK}
	c.JSON(http.StatusOK, vginResponse)

}
func FunGetPokemonFromApi() error {

	var (
		responseObject models.ApiStruct
		vPokemon       models.ObjPokemonPost
	)

	UrlApi := "https://pokeapi.co/api/v2/pokemon-form/"

	response, err := http.Get(UrlApi)
	if err != nil {
		err = errors.New("*Error GetPokemonFromApi: couldn't get Api: " + UrlApi)
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &responseObject)

	PokemonCount := responseObject.Count
	PokemonResponse := responseObject.Results
	println(PokemonCount)

	for _, poke := range PokemonResponse {

		vPokemon.Name = poke.Name
		vPokemon.Url = poke.Url

		err = models.FunPostPokemon(vPokemon)
		if err != nil {
			err = errors.New("*ERROR GetPokemonFromApi: couldn't register pokemon -> " + err.Error())
			return err
		}
	}

	return nil

}
