package main

import (
	"fmt"
	"github.com/tasks/GolangPythonMicroservice/Go/config"
	"github.com/tasks/GolangPythonMicroservice/Go/controllers"
	"log"
)

func init() {
	err := config.FunInitConfig()
	if err != nil {
		log.Fatal("*ERROR init: couldn't initialize configuration -> ", err.Error())
	}

	vsessionPokedex, err := config.FunOpenDatabaseConnection()
	defer vsessionPokedex.Close()
	if err != nil {
		fmt.Println("*ERROR init: couldn't connect database -> ", err.Error())
	}
}

func main() {
	controllers.PokesRouter()
	//controllers.FunGetPokemonFromApi()
}
