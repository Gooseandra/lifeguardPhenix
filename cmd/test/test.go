// Code generated by go-swagger; DO NOT EDIT.

package main

import (
	"database/sql"
	"log"
	"os"
	"swagger/handlers"
	"swagger/services"
	"swagger/services/auth/sessions"
	"swagger/storages/callPostgres"
	"swagger/storages/crewPostgres"
	"swagger/storages/inventoryPostgres"
	"swagger/storages/userPostgres"
	"time"

	"github.com/go-openapi/loads"
	"github.com/go-yaml/yaml"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"

	"swagger/restapi"
	"swagger/restapi/operations"
)

// This file was generated by the swagger tool.
// Make sure not to overwrite this file after you generated it because all your edits would be lost!

func main() {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewMchsAPI(swaggerSpec)

	var settings Settings
	bytes, fail := os.ReadFile("db.yml")
	if fail != nil {
		log.Println(fail.Error())
		log.Panic(fail.Error())
	}
	fail = yaml.Unmarshal([]byte(bytes), &settings)
	if fail != nil {
		log.Panic(fail.Error())
	}
	log.Println(settings)
	db, err := sql.Open(settings.Database.Type, settings.Database.Arguments)
	if err != nil {
		log.Println(err.Error())
	}
	userStorage := userPostgres.NewStorage(db)
	if fail != nil {
		log.Panic(fail.Error())
	}

	//userStorage.New("root", "root", "123")

	crewStorage := crewPostgres.NewStorage(db)
	callStorage := callPostgres.NewStorage(db)
	inventoryStorage := inventoryPostgres.NewStorage(db)

	//crewStorage.New(time.Now(), 16, "none", []uint64{16,17})
	log.Println(crewStorage.List(0, 100))
	logService := services.NewLog()
	userService := services.NewUsers(userStorage)
	crewService := services.NewCrew(crewStorage)
	callService := services.NewCall(callStorage)
	inventoryService := services.NewInventory(inventoryStorage)

	sessionService := sessions.NewSessions(logService, userService, time.Hour)

	//api.CreateEventHandler = handlers.NewCreateEvent(logService, &sessionService, userService)
	api.DeleteInventoryItemHandler = handlers.NewDeleteInventoryItem(logService, &sessionService, inventoryService)
	api.UpdateInventoryHandler = handlers.NewUpdateInventoryItem(logService, &sessionService, inventoryService)
	api.GetInventoryTypesHandler = handlers.NewInventoryTypes(logService, &sessionService, inventoryService)
	api.GetInventoryItemHandler = handlers.NewByIDInventoryItem(logService, &sessionService, inventoryService)
	api.CreateInventoryItemHandler = handlers.NewCreateInventoryItem(logService, &sessionService, inventoryService)
	api.ListInventoryItemsHandler = handlers.NewListInventoryItem(logService, &sessionService, inventoryService)

	api.ListCallHandler = handlers.NewListCall(logService, &sessionService, callService)
	api.CreateCallHandler = handlers.NewCreateCall(logService, &sessionService, callService)
	api.UpdateCallHandler = handlers.NewUpdateCall(logService, &sessionService, callService)
	api.GetCallHandler = handlers.NewGetCall(logService, &sessionService, callService)

	api.UpdateCrewHandler = handlers.NewUpdateCrew(logService, &sessionService, crewService)
	api.CreateCrewHandler = handlers.NewCreateCrew(logService, &sessionService, crewService)
	api.ListCrewHandler = handlers.NewListCrew(logService, &sessionService, crewService)
	api.GetCrewHandler = handlers.NewGetCrew(logService, &sessionService, crewService)

	api.CreateUserHandler = handlers.NewCreateUser(logService, &sessionService, userService)
	api.GetUserHandler = handlers.NewGetUser(logService, &sessionService, userService)
	api.ListUsersHandler = handlers.NewListUser(logService, &sessionService, userService)
	api.UpdateUserHandler = handlers.NewUpdateUser(logService, &sessionService, userService)
	api.FiredUserHandler = handlers.NewFiredUser(logService, &sessionService, userService)

	api.LoginHandler = handlers.NewLogin(logService, &sessionService, userService)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "mchs"
	parser.LongDescription = "Приложение для МЧС"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}

}
