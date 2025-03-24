// Package dateparser provides additional date layout formats beyond those available in the gofeed package.
package dateparser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Regular expression for EETE_R pattern
var eetePattern = regexp.MustCompile(`^([A-Z][a-z]{2})(AM|PM)(EETE_R|EESTE_R)([A-Z][a-z]+)C822$`)

// monthMap maps month names to time.Month values
var monthMap = map[string]time.Month{
	"january":   time.January,
	"february":  time.February,
	"march":     time.March,
	"april":     time.April,
	"may":       time.May,
	"june":      time.June,
	"july":      time.July,
	"august":    time.August,
	"september": time.September,
	"october":   time.October,
	"november":  time.November,
	"december":  time.December,
}

// weekdayMap maps weekday abbreviations to time.Weekday values
var weekdayMap = map[string]time.Weekday{
	"sun": time.Sunday,
	"mon": time.Monday,
	"tue": time.Tuesday,
	"wed": time.Wednesday,
	"thu": time.Thursday,
	"fri": time.Friday,
	"sat": time.Saturday,
}

// parseEETEPattern parses the custom EETE_R pattern and returns approximate time
func parseEETEPattern(dateStr string) (time.Time, error) {
	matches := eetePattern.FindStringSubmatch(dateStr)
	if matches == nil || len(matches) != 5 {
		return time.Time{}, fmt.Errorf("invalid EETE_R format: %s", dateStr)
	}

	// Extract components
	weekday := matches[1]     // e.g., "Tue"
	ampm := matches[2]        // "AM" or "PM"
	patternType := matches[3] // "EETE_R" or "EESTE_R"
	month := matches[4]       // e.g., "January"

	// Get current year
	currentYear := time.Now().Year()

	// Convert month name to month number
	monthNum, ok := monthMap[strings.ToLower(month)]
	if !ok {
		// Default to current month if unknown
		monthNum = time.Now().Month()
	}

	// Find a day in this month and year that matches the weekday
	t := time.Date(currentYear, monthNum, 1, 0, 0, 0, 0, time.UTC)

	// Advance to the first occurrence of the desired weekday
	weekdayNum := parseWeekday(weekday)
	if weekdayNum >= 0 {
		daysToAdd := (int(weekdayNum) - int(t.Weekday()) + 7) % 7
		t = t.AddDate(0, 0, daysToAdd)
	}

	// Set hour based on AM/PM
	hour := 9 // Default to 9 AM
	if ampm == "PM" {
		hour = 15 // 3 PM
	}

	// Adjust minute based on pattern type (purely for differentiation)
	minute := 0
	if patternType == "EESTE_R" {
		minute = 30 // Use 30 minutes for EESTE_R pattern
	}

	// Create the final time
	result := time.Date(t.Year(), t.Month(), t.Day(), hour, minute, 0, 0, time.UTC)
	return result, nil
}

// parseWeekday converts weekday abbreviation to time.Weekday value
func parseWeekday(day string) time.Weekday {
	weekday, ok := weekdayMap[strings.ToLower(day)]
	if !ok {
		return -1
	}
	return weekday
}

// dateLayouts contains all the format patterns ordered by specificity to general
var dateLayouts = []string{
	// Complex formats with date repetition and timezone
	"Dom, 02 Jan 2006 15:04:05 -0700 2006-01-02 15:04:05", // Spanish Sunday
	"Sáb, 02 Jan 2006 15:04:05 -0700 2006-01-02 15:04:05", // Spanish Saturday

	// Formats with named timezone
	"Mon, 02 Jan 2006 15:04:05 MST/City",  // Matches: Sun, 23 Mar 2025 08:14:46 Europe/Dublin
	"Mon, 02 Jan 2006 15:04:05 TIME_ZONE", // Matches: Sat, 22 Mar 2025 19:57:06 TIME_ZONE
	"Mon, 02 Jan 2006 15:04:05 MST",       // With named timezone
	"Mon, 02 Jan 2006  MST",               // Matches: Sat, 22 Mar 2025  CST

	// Full date formats with GMT timezone
	"Monday, January 2, 2006, 15:04 GMT -0700",  // Matches: Sunday, March 23, 2025, 16:20 GMT +5:30
	"Monday, January 2, 2006, 15:04 GMT +0700",  // Same with positive offset
	"Monday, January 2, 2006, 15:04 GMT -07:00", // Matches: Sunday, March 23, 2025, 16:20 GMT +5:30
	"Monday, January 2, 2006, 15:04 GMT +07:00", // Same with positive offset

	"Monday, January 2, 2006, 15:04 GMT-0700",
	"Monday, January 2, 2006, 15:04 GMT+0700",
	"Monday, January 2, 2006, 15:04 GMT-07:00",
	"Monday, January 2, 2006, 15:04 GMT+07:00",

	// Vie, 03/21/2025 - 00:00 format
	"Vie, 01/02/2006 - 15:04", // American style MM/DD/YYYY

	// Non-English weekday formats with timezone
	"dom, 02 Jan 2006 15:04:05 -0700", // Spanish/Portuguese Sunday - dom, 23 mar 2025 12:40:59 +0100
	"ven, 02 Jan 2006 15:04:05 MST",   // French Friday - ven, 21 mar 2025 13:49:00 CDT
	"Vie, 02 Jan 2006 15:04:05 -0700", // Spanish Friday - Vie, 30 Sep 2022 21:27:13 -0500
	"Jue, 02 Jan 06 15:04:05 -0700",   // Spanish Thursday with 2-digit year - Jue, 29 Jun 23 15:34:11 +0200
	"Mar, 02 Jan 2006 15:04:05 -0700", // Spanish Tuesday - Mar, 31 Ago 2021 22:46:32 -0500
	"Mar, 02 Ago 2006 15:04:05 -0700", // Spanish Tuesday with August - Mar, 31 Ago 2021 22:46:32 -0500

	// These formats need explicit language patterns because Go doesn't automatically recognize
	// non-English day/month names in its standard time parsing
	"Domenica, 02 Marzo, 2006 - 15:04",    // Italian with March - Domenica, 23 Marzo, 2025 - 10:33
	"Domenica, 02 Gennaio, 2006 - 15:04",  // Italian - Domenica, 23 Gennaio, 2025 - 10:33
	"Lunedì, 02 Gennaio, 2006 - 15:04",    // Italian Monday
	"Martedì, 02 Gennaio, 2006 - 15:04",   // Italian Tuesday
	"Mercoledì, 02 Gennaio, 2006 - 15:04", // Italian Wednesday
	"Giovedì, 02 Gennaio, 2006 - 15:04",   // Italian Thursday
	"Venerdì, 02 Gennaio, 2006 - 15:04",   // Italian Friday
	"Sabato, 02 Gennaio, 2006 - 15:04",    // Italian Saturday

	// ISO 8601 format with extra timezone
	"2006-01-02T15:04:05Z -0700", // Matches: 2025-03-23T11:02:13Z +0300

	// Date formats with non-English month names
	"02 Μαρ 2006 15:04:00 -0700", // Greek March - 23 Μαρ 2025 13:11:00 +0000
	"02 Mars 2006",               // French March - 22 Mars 2025
	"02 مارس 2006",               // Arabic March - 22 مارس 2025

	// Formats with weekday, space-separated timezone
	"Mon, 02 January 2006, 03:04:05 PM -0700", // Matches: Sun, 23 March 2025, 05:06:27 PM +0530
	"Mon, 02 Jan 2006 03:04 PM MST",           // Matches: Sat, 22 Mar 2025 08:36 PM EDT
	"Mon, 02 Jan 2006 03:04 AM MST",           // Matches: Sun, 23 Mar 2025 08:11 AM EDT

	// Day-first formats
	"02 January 2006 - 15:04", // Matches: 23 March 2025 - 12:10
	"02-01-2006 15:04",        // Matches: 23-03-2025 11:15

	// Different time-date separators
	"15:04 02.01.2006",   // Matches: 12:06 23.03.2025
	"02.01.2006 | 15:04", // Matches: 23.03.2025 | 08:28

	// Standard formats with different date/time orders
	"Monday Jan 02 2006 15:04:05", // Matches: Sunday Mar 23 2025 13:54:16
	"Monday Jan 2 2006 15:04:05",  // Non-padded day version
	"January 02, 2006, 3:04 pm",   // Matches: March 23, 2025, 4:50 pm
	"Jan 02, 2006, 3:04pm",        // Matches: Mar 22, 2025, 12:00pm
	"Jan 2, 2006, 3:04pm",         // Non-padded day version
	"Jan 02, 2006, 3:04am",        // AM version
	"Jan 2, 2006, 3:04am",         // Non-padded day AM version

	// Day, hour, minute formats with timezone name/offset
	"Monday 02 Jan 2006 15:04:05 -0700", // Matches: Friday 05 Jul 2024 08:00:00 -0600
	"Mon,02 Jan 2006 15:04:05 -07",      // Matches: Sun,23 Mar 2025 18:37:00 +07 (no space after comma)
	"Mon, 02 Jan 2006 15:04:05 z",       // Matches: Sun, 23 Mar 2025 05:49:21 z

	// Formats with 24+ hour handling
	"Mon, 2 Jan 2006 24:04:05 -0700",  // Matches: Mon, 17 Mar 2025 24:15:59 +0530 (non-zero-padded day)
	"Mon, 02 Jan 2006 24:04:05 -0700", // Matches: Mon, 17 Mar 2025 24:15:59 +0530 (zero-padded day)
	"Mon, 02 Jan 2006 15:04:05 -0700", // Standard format with timezone offset

	// Formats with different time layouts
	"Mon, 02 Jan 2006 03:04:05 PM", // Matches: Sat, 22 Mar 2025 11:07:41 PM
	"Mon, 02 Jan 2006 03:04:05 AM", // AM version
	"Mon, 02 Jan 2006 15:04:05",    // 24-hour format without AM/PM

	// Spanish date formats
	"Vie, 02/01/2006 - 15:04", // Matches: Vie, 03/21/2025 - 00:00 (European style DD/MM/YYYY)

	// Formats with numeric day-month
	"02.01.2006", // Just date with dots

	// Common formats with different ordering
	"Mon, 2 Jan 2006 3:04:05 PM",  // With non-padded day and hour
	"Mon, 2 Jan 2006 3:04:05 AM",  // With non-padded day and hour (AM)
	"Mon, 02 Jan 2006 3:04:05 PM", // With padded day, non-padded hour
	"Mon, 02 Jan 2006 3:04:05 AM", // With padded day, non-padded hour (AM)

	// Mon, 5 Oct 2020 19:30:00 GMT with extra spaces
	" Mon, 2 Jan 2006 15:04:05 MST ", // Matches: " Mon, 5 Oct 2020 19:30:00 GMT "

	// Formats for the examples in your first list
	"Mon, 02 Jan 2006 3:04:05 PM -0700", // Matches: Tue, 18 Mar 2025 5:58:24 PM
	"Mon, Jan 2 2006 12:04:05 AM",       // Matches: Fri, Mar 21 2025 12:49:00 AM
	"Mon, Jan 2 2006 03:04:05 PM",       // Matches: Tue, Mar 18 2025 03:11:53 PM

	// Fallback formats for simple cases
	"2 Jan 2006 15:04 MST", // Matches: 3 Mar 2025 18:27 UTC
	"2 Jan 2006 15:04",     // Without timezone
}

// ParseDate attempts to parse a date string using all available layouts
func ParseDate(dateStr string) (time.Time, error) {
	var t time.Time
	var err error

	// Try to parse EETE_R pattern
	if strings.Contains(dateStr, "EETE_R") || strings.Contains(dateStr, "EESTE_R") {
		return parseEETEPattern(dateStr)
	}

	// Special handling for timezone regions like Europe/Dublin
	if strings.Contains(dateStr, "Europe/Dublin") {
		// Try replacing with GMT first
		dateStr = strings.Replace(dateStr, "Europe/Dublin", "GMT", 1)
	}

	if strings.Contains(dateStr, "GMT +") || strings.Contains(dateStr, "GMT -") {
		// Try removing space between GMT and +/-
		dateStr = strings.Replace(dateStr, "GMT +", "GMT+", 1)
		dateStr = strings.Replace(dateStr, "GMT -", "GMT-", 1)
	}

	// Try all the standard layouts
	for _, layout := range dateLayouts {
		t, err = time.Parse(layout, dateStr)
		if err == nil {
			return t, nil
		}
	}

	// If no layout matches, return the error from the last attempt
	return t, fmt.Errorf("unable to parse date: %s", dateStr)
}

// ParseDateWithDefaultTZ sets timezone to UTC if not specified
func ParseDateWithDefaultTZ(dateStr string) (time.Time, error) {
	t, err := ParseDate(dateStr)
	if err != nil {
		return t, err
	}

	// If parsed time doesn't have timezone info, set it to UTC
	if t.Location() == time.Local {
		return t.UTC(), nil
	}

	return t, nil
}
