package handlers

import "fmt"

func validate_row(channel chan *MotorPolicy, trigger chan bool) {
	for row_data := range channel {
		fmt.Println(*row_data)
	}
	trigger <- true
}
