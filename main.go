package main

import (
	"net/http"

	"github.com/Deikioveca/AuthenticationSystem/v2/database"
	"github.com/Deikioveca/AuthenticationSystem/v2/server"
	"github.com/Deikioveca/AuthenticationSystem/v2/utils"
)

func main() {
	database.InitRedis()
	utils.LoadTemplate()
	router := server.Server{}.InitServer()
	http.ListenAndServe(":8080", router)
}
