package examples

import (
	"../models"
	"fmt"
	"log"
	"time"
)

func create() {
	u := models.User{
		Name:     "Juan Carlos Perez Romirez",
		Password: "1234",
		Email:    "juancarlos.perez.romirez@gmail.com",
		Active:   false,
		Age:      time.Now(),
	}
	id, err := models.CreateUser(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Creado exitosamente: ", id)
}
func consultar() {
	u, err := models.GetUsers()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(u)
}
func update() {
	u := models.User{
		Id:     3,
		Active: true,
		Email:  "carlos.update@gmail.com",
		Name:   "Carlos Perez",
	}
	err := models.UpdateUser(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modificado exitosamente")
}
func delete() {
	id := 3
	err := models.DeleteUser(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Se ha eliminado exitosamente")
}
func example() {
	//create()
	//consultar()
	//update()
	//delete()
}
