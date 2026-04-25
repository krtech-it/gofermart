package service

import "strconv"

func isValidLuhn(orderNumber string) bool {
	sumLuhn := 0
	chet := len(orderNumber)%2 == 0
	for i := len(orderNumber) - 1; i >= 0; i-- {
		number, err := strconv.Atoi(string(orderNumber[i]))
		if err != nil {
			return false
		}
		if !chet && i%2 == 1 || chet && i%2 == 0 {
			number = number * 2
			if number > 9 {
				number = number - 9
			}
		}
		sumLuhn += number
	}
	return sumLuhn%10 == 0
}
