package main

import (
	"fmt"
	"os"

	"github.com/thomasbeukema/somtoday_api/Go"
)

func main() {
	som := new(somtoday.Somtoday)
	var credentials [4]string
	credentials[0] = "271488"      //"YOUR USERNAME"
	credentials[1] = "Thomas16tb@" //"YOUR PASSWORD"
	credentials[2] = "fiotey"      //"YOUR SCHOOL"
	credentials[3] = "02KB"        //"YOUR BRIN"
	som.SetCredentials(credentials)
	err := som.Login()

	if err != nil {
		fmt.Printf("%s", err)
	}

	timetable, err := som.GetTimetable(0) // Retrieve timetable; Parameter is how many days ahead: 0 for today, 1 for tommorow etc.
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	homework, err := som.GetHW(31) // Retrieves homework; Parameter: for how many days do you want to get hw?
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	grades, err := som.GetGrades(0) // Retrieves grades; Parameter: '0' for 10 most recent grades, '1' for all grades
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	if som.ChangeHwDone("YOUR APPOINTMENTID", "YOUR HWID", true) {
		fmt.Println("Changed hw status succesfully")
	} else {
		fmt.Println("Changed hw status not so succesfully :(")
	}

	fmt.Printf("Timetable:\n%s\n\n", timetable)
	fmt.Printf("Homework:\n%s\n\n", homework)
	fmt.Printf("Grades:\n%s\n\n", grades)
}
