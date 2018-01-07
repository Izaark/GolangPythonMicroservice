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
	//"strings"
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
		pokemon.GET("/pokemon/:id", handlerGetPokemon)
		pokemon.POST("/register/pokemon", handlerPostPokemon)
		pokemon.PUT("/update/pokemon/:id", handlerUpdatePokemon)
		pokemon.DELETE("/delete/pokemon/:id", handlerDeletePokemon)

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

	pokemon, err := models.FunGetAllPokemon()
	if err != nil {
		vginResponse = gin.H{"message": "error quering pokes", "response": nil, "error": "RE", "status": http.StatusBadRequest}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}
	for _, poke := range pokemon {

		if vPokemon.Name == poke.Name {
			vginResponse = gin.H{"message": "The pokemon exists", "response": nil, "error": "RE", "status": http.StatusConflict}
			c.JSON(http.StatusConflict, vginResponse)
			return
		}
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
	err = models.FunDeletePokemon(strID)
	if err != nil {
		vginResponse = gin.H{"message": "error delete pokemon", "response": nil, "error": "RE"}
		c.JSON(http.StatusInternalServerError, vginResponse)
		return
	}
	vginResponse = gin.H{"message": "Pokemon successfully deleted", "error": nil, "status": http.StatusOK}
	c.JSON(http.StatusOK, vginResponse)

}

func handlerGetPokemon(c *gin.Context) {
	var (
		vPokemon models.ObjPokemonGet
		response gin.H
		err      error
	)

	strID := c.Params.ByName("id")
	vPokemon, err = models.FunGetPokemon(strID)
	if err != nil {
		response = gin.H{"message": "error quering a poke", "response": nil, "error": "RE"}
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	//Todo: create a funcion if exists pokemon
	response = gin.H{"message": "pokemon found", "error": nil, "status": http.StatusOK, "response": vPokemon}
	c.JSON(http.StatusOK, response)

}

func FunGetPokemonFromApi() error {

	var (
		responseObject models.ApiStruct
		vPokemon       models.ObjPokemonPost
	)

	UrlApi := "https://pokeapi.co/api/v2/pokemon-form/"
	println(UrlApi)

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
		println(vPokemon.Name)

		err = models.FunPostPokemon(vPokemon)
		if err != nil {
			err = errors.New("*ERROR GetPokemonFromApi: couldn't register pokemon -> " + err.Error())
			return err
		}
	}

	return nil

}
