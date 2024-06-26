package plantillas

import "time"

// ================================================================ //
//
// Allows for easy calculation of the age of an entity,
// provided with the date of birth of that entity.
// https://github.com/bearbin/go-age/blob/master/age.go
//
// ================================================================ //

// EdadEn gets the age of an entity at a certain time.
func EdadEn(birthDate time.Time, now time.Time) int {
	// Get the year number change since the player's birth.
	years := now.Year() - birthDate.Year()

	// If the date is before the date of birth, then not that many years have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years -= 1
	}

	return years
}

// ================================================================ //

// Edad is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
func Edad(birthDate time.Time) int {
	return EdadEn(birthDate, time.Now())
}

// ================================================================ //

// EdadPtr is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
func EdadPtr(birthDate *time.Time) int {
	return EdadEn(*birthDate, time.Now())
}

// ================================================================ //

// Gets the adjusted date of birth to work around leap year differences.
func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

// Works out if a time.Time is in a leap year.
func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}
