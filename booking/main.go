package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var conferenceName = "Happy Concert"

const conferenceTicket = 50

var remainingTickets int = 50

var wg = sync.WaitGroup{}

type UserData struct {
	userName   string
	city       string
	userTicket int
}

func main() {
	var bookings = make([]UserData, 0)

	greetUser()

	for remainingTickets > 0 && len(bookings) < 50 {
		userData := GetUserInput()

		userTicket := userData.userTicket

		if userTicket > remainingTickets {
			fmt.Printf("We only have %v tickets remaining, you can't book %v tickets\n", remainingTickets, userTicket)
			continue
		}

		remainingTickets = remainingTickets - userTicket
		bookings = append(bookings, userData)

		fmt.Printf("Booking %v tickets\n", remainingTickets)
		wg.Add(1)
		go sendTicket(userData.userName, userTicket)

		firsNames := getFirstNames(bookings)
		fmt.Printf("First names %v\n", firsNames)
		fmt.Printf("All bookings %v\n", bookings)

		if remainingTickets == 0 {
			fmt.Println("Our conference is out of tickets. Please come back next year")
			break
		}
	}
	wg.Wait()
}

func greetUser() {
	fmt.Printf("Welcome to the %v.\n", conferenceName)
	fmt.Printf("We have total of %v tickets and %v are still available.\n", conferenceTicket, remainingTickets)
	fmt.Println("Get your tickets here to attend")
}

func getFirstNames(bookings []UserData) []string {
	var firsNames []string
	for _, booking := range bookings {
		var names = strings.Fields(booking.userName)
		firsNames = append(firsNames, names[0])
	}
	return firsNames
}

func sendTicket(userName string, userTicket int) {
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("Sending ticket %v to user %v\n", userTicket, userName)
	wg.Done()
}
