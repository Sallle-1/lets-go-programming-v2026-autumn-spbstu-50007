package main

import "fmt"

func main() {
	var departments int
	_, err := fmt.Scan(&departments)
	if err != nil || departments < 1 || departments > 1000 {
		fmt.Println("Invalid number of departments")
		return
	}

	for department := 0; department < departments; department++ {
		var employees int
		_, err = fmt.Scan(&employees)
		if err != nil || employees < 1 || employees > 1000 {
			fmt.Println("Invalid number of employees")
			return
		}

		minTemperature := 15
		maxTemperature := 30

		for employee := 0; employee < employees; employee++ {
			var operator string
			var temperature int
			_, err = fmt.Scan(&operator, &temperature)
			if err != nil {
				fmt.Println("Invalid temperature requirement")
				return
			}

			switch operator {
			case ">=":
				if temperature > minTemperature {
					minTemperature = temperature
				}
			case "<=":
				if temperature < maxTemperature {
					maxTemperature = temperature
				}
			default:
				fmt.Println("Invalid operator")
				return
			}

			if minTemperature > maxTemperature {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemperature)
			}
		}
	}
}
