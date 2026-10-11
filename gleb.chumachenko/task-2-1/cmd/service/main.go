package main

import "fmt"

func getTemperature(operator string, temperature int, minTemperature *int, maxTemperature *int) int {
	switch operator {
	case ">=":
		if temperature > *minTemperature {
			*minTemperature = temperature
		}
	case "<=":
		if temperature < *maxTemperature {
			*maxTemperature = temperature
		}
	}

	if *minTemperature > *maxTemperature {
		return -1
	}

	return *minTemperature
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil || departments < 1 || departments > 1000 {
		fmt.Println("Invalid number of departments")

		return
	}

	for range departments {
		var employees int

		_, err = fmt.Scan(&employees)
		if err != nil || employees < 1 || employees > 1000 {
			fmt.Println("Invalid number of employees")

			return
		}

		minTemperature := 15
		maxTemperature := 30

		for range employees {
			var operator string

			var temperature int

			_, err = fmt.Scan(&operator, &temperature)
			if err != nil {
				fmt.Println("Invalid temperature requirement")

				return
			}

			if operator != ">=" && operator != "<=" {
				fmt.Println("Invalid operator")

				return
			}

			result := getTemperature(operator, temperature, &minTemperature, &maxTemperature)
			fmt.Println(result)
		}
	}
}
