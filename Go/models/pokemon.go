package models

import (
	"errors"
	"github.com/tasks/GolangPythonMicroservice/Go/config"
	r "gopkg.in/gorethink/gorethink.v3"
)

const (
	CstTrainerTable = "trainer"
	CstPokemonTable = "pokemon"
)

type ApiStruct struct {
	Results []ApiResults
	Count   int `json:"count, omitempty"`
}
type ApiResults struct {
	Name string `json:"name" binding:"required"`
	Url  string `json:"url, omitempty"`
}

type ObjPokemonPost struct {
	Name string `json:"name" binding:"required" gorethink:"name"`
	Url  string `json:"url, omitempty" binding:"required"  gorethink:"url"`
}
type ObjPokemonGet struct {
	Name string `json:"name" binding:"required" gorethink:"name"`
	Url  string `json:"url, omitempty" gorethink:"url"`
}

func FunGetAllPokemon() ([]ObjPokemonGet, error) {

	var vPokes []ObjPokemonGet
	var cursor *r.Cursor

	vSessiondb, err := config.FunOpenDatabaseConnection()
	defer vSessiondb.Close()
	if err != nil {
		err = errors.New("*Error FunGetAllPokemon: couldn't connect database -> " + err.Error())
		return nil, err
	}

	cursor, err = r.Table(CstPokemonTable).Run(vSessiondb)
	defer cursor.Close()
	if err != nil {
		err = errors.New("*Error FunGetAllPokemon: couldn't retrieve pokemons" + " -> " + err.Error())
		return vPokes, err
	}

	err = cursor.All(&vPokes)
	if err != nil {
		err = errors.New("*Error FunGetAllPokemon:  couldn't use cursor to retrieve pokemons" + " --> " + err.Error())
		return vPokes, err
	}
	return vPokes, nil

}

func FunPostPokemon(vPoke ObjPokemonPost) error {

	vSessionDb, err := config.FunOpenDatabaseConnection()
	defer vSessionDb.Close()
	if err != nil {
		err = errors.New("*Error FunPostPokemon: + couldn't connect database -> " + err.Error())
		return err
	}
	_, err = r.Table(CstPokemonTable).Insert(vPoke).RunWrite(vSessionDb)
	if err != nil {
		err = errors.New("*Error FunPostPokemon couldn't insert new Pokemon" + vPoke.Name + "->" + err.Error())
		return err
	}
	return nil

}

//todo: update whit name or id !
func FunUpdatePokemon(vPoke ObjPokemonPost, strID string) error {

	vSessionDb, err := config.FunOpenDatabaseConnection()
	defer vSessionDb.Close()
	if err != nil {
		err = errors.New("*Error FunUpdatePokemon: + couldn't connect database -> " + err.Error())
		return err
	}
	_, err = r.Table(CstPokemonTable).Get(strID).Update(vPoke).RunWrite(vSessionDb)
	if err != nil {
		err = errors.New("*Error FunUpdatePokemon couldn't update Pokemon" + vPoke.Name + "->" + err.Error())
		return err
	}
	return nil

}

func FunDeletePokemon(strID string) error {

	vSessionDb, err := config.FunOpenDatabaseConnection()
	defer vSessionDb.Close()
	if err != nil {
		err = errors.New("*Error FunDeletePokemon: + couldn't connect database -> " + err.Error())
		return err
	}
	_, err = r.Table(CstPokemonTable).Get(strID).Delete().Run(vSessionDb)
	if err != nil {
		err = errors.New("*Error FunDeletePokemon couldn't Delete Pokemon" + "->" + err.Error())
		return err
	}
	return nil

}

func FunGetPokemon(strID string) (ObjPokemonGet, error) {

	var vPokes ObjPokemonGet
	var cursor *r.Cursor

	vSessiondb, err := config.FunOpenDatabaseConnection()
	defer vSessiondb.Close()
	if err != nil {
		err = errors.New("*Error FunGetAllPokemon: couldn't connect database -> " + err.Error())
		return vPokes, err
	}

	cursor, err = r.Table(CstPokemonTable).Get(strID).Run(vSessiondb)
	defer cursor.Close()
	if err != nil {
		err = errors.New("*Error FunGetAllPokemon: couldn't retrieve pokemon" + " -> " + err.Error())
		return vPokes, err
	}

	err = cursor.One(&vPokes)
	if err != nil {
		err = errors.New("*Error FunGetAllPokemon:  couldn't use cursor to retrieve pokemon" + " --> " + err.Error())
		return vPokes, err
	}
	return vPokes, nil

}

func FunExistPokemon(strField, strID string) (bool, error) {
	var vCursor *r.Cursor
	var vboolExist bool
	var vintCounter int

	vSessionDb, err := config.FunOpenDatabaseConnection()
	defer vSessionDb.Close()
	if err != nil {
		err = errors.New("*ERROR FunExistPokemon: couldn't connect database -> " + err.Error())
		return false, err
	}

	vCursor, err = r.Table(CstPokemonTable).GetAllByIndex(strField, strID).Count(r.Row.Field("id")).Run(vSessionDb)
	defer vCursor.Close()

	if err != nil {
		err = errors.New("*ERROR FunExistPokemon: couldn't verify pokemon information " + "-> " + err.Error())
		return false, err
	}
	err = vCursor.One(&vintCounter)

	if err != nil {
		err = errors.New("*ERROR FunExistPokemon: couldn't use cursor to verify information" + err.Error())
		return false, err
	}

	switch vintCounter {
	case 0:
		err = nil
		vboolExist = false
	case 1:
		err = nil
		vboolExist = true
	default:
		err = errors.New("*ERROR FunExistPokemon: more than one pokenon")
		vboolExist = false
	}

	return vboolExist, err

}
