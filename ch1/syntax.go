package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type user struct {
	Name            string
	Birthday        time.Time
	Password        string
	FavoriteWeekday Weekday
}

func main() {
	first := user{
		Name:            "Gigel",
		Password:        "you'll never guess it",
		Birthday:        time.Date(1980, time.July, 4, 0, 0, 0, 0, time.UTC),
		FavoriteWeekday: Friday,
	}

	fmt.Println(first)

	for i := 0; i < 5; i++ {
		awesomeness, err := beAwesome()
		if err == nil {
			fmt.Println(awesomeness)
		} else {
			fmt.Println(err)
		}
	}

	var durations [3]time.Duration
	durations[0], _ = time.ParseDuration("300ms")
	durations[1], _ = time.ParseDuration("1024s")
	durations[2], _ = time.ParseDuration("1024000ms")

	for i, d := range durations {
		fmt.Printf("%d: %s\n", i, d)
	}
}

func beAwesome() (string, error) {
	chance := rand.Intn(3)
	if chance/2 == 0 {
		return "Yeah, I'm awesome", nil
	} else {
		return "", errors.New("Not now, I'm busy")
	}
}
