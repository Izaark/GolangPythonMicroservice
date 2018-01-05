package config

import (
	"errors"
	r "gopkg.in/gorethink/gorethink.v3"
	"os"
	"strconv"
	"time"
)

func FunOpenDatabaseConnection() (*r.Session, error) {

	var vSession *r.Session
	var vstrDatabase string
	var vintTimeoutSecs int
	var err error

	vstrDatabase = os.Getenv("POK_ENV_DATABASE")

	// Preventing error while parsing timeout from environment variables
	vintTimeoutSecs, err = strconv.Atoi(os.Getenv("POK_DATABASE_TIMEOUT_SECS"))
	if err != nil {
		vintTimeoutSecs = 1
	}

	// Creating connection with defined parameters
	vSession, err = r.Connect(r.ConnectOpts{
		Address:  os.Getenv("POK_ENV_DATABASE_ADDRESS"),
		Database: vstrDatabase,
		Timeout:  time.Duration(vintTimeoutSecs) * time.Second,
	})
	if err != nil {
		err = errors.New("*ERROR FunOpenDatabaseConnection: couldn't connect to rethinkdb -> " + err.Error())
		return vSession, err
	}

	result, err := r.Expr("-> Database successfully connected").Run(vSession)
	defer result.Close()
	if err != nil {
		err = errors.New("*ERROR FunOpenDatabaseConnection: couldn't connect " + vstrDatabase + " database -> " + err.Error())
		return vSession, err
	}

	response := ""
	err = result.One(&response)
	if err != nil {
		err = errors.New("*ERROR FunOpenDatabaseConnection: " + vstrDatabase + " database is not responding -> " + err.Error())
		return vSession, err
	}

	//fmt.Println(response)
	return vSession, nil
}
