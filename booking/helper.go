package main

import (
	"fmt"
	"strings"
)

func GetUserInput() UserData {
	var city string
	var userName string
	var userTicket int

	fmt.Print("Choose a city: ")
	fmt.Scan(&city)
	switch city {
	case "Berlin":
		fmt.Println("you select Berlin")
	case "Singapore":
		fmt.Println("you select Singapore")
	case "London":
		fmt.Println("you select London")
	default:
		fmt.Println("No valid city")
	}

	fmt.Println("Please enter a user name: ")
	fmt.Scan(&userName)
	userName = strings.Replace(userName, "\n", "", -1)
	fmt.Println("Please enter a ticket number: ")
	fmt.Scan(&userTicket)

	var userData = UserData{
		userName:   userName,
		userTicket: userTicket,
		city:       city,
	}
	return userData
}
