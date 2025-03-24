package dateparser

import (
	"testing"
)

func TestDateParsing(t *testing.T) {
	testCases := []struct {
		name    string
		dateStr string
		wantErr bool
	}{
		{"EETE_R AM Format", "TueAMEETE_RMarchC822", false},
		{"Spanish Date with Timezone", "Vie, 30 Sep 2022 21:27:13 -0500", false},
		{"English with Single Digit Hour", "Tue, 18 Mar 2025 5:58:24 PM", false},
		{"Spanish Date with Slash", "Vie, 03/21/2025 - 00:00", false},
		{"EETE_R PM Format", "SunPMEETE_RJanuaryC822", false},
		{"Complex Date with Repetition", "Dom, 23 Mar 2025 00:05:36 +0000 2025-03-23 00:05:36", false},
		{"Weekday Month Day Format", "Sunday Mar 23 2025 04:37:41", false},
		{"Greek Month", "23 Μαρ 2025 13:11:00 +0000", false},
		{"Dotted Date with Pipe", "23.03.2025 | 08:28", false},
		{"Named Timezone Region", "Sun, 23 Mar 2025 08:14:46 Europe/Dublin", false},
		{"Month Day Year with AM/PM", "Mar 22, 2025, 12:00pm", false},
		{"Time First Format", "12:06 23.03.2025", false},
		{"Time with Named Timezone", "Sat, 22 Mar 2025 08:36 PM EDT", false},
		{"Italian Date Format", "Domenica, 23 Marzo, 2025 - 10:33", false},
		{"ISO with Extra TZ", "2025-03-23T11:02:13Z +0300", false},
		{"French Weekday with TZ", "ven, 21 mar 2025 13:49:00 CDT", false},
		{"Extra Spaces in Date", " Mon, 5 Oct 2020 19:30:00 GMT ", false},
		{"Arabic Date", "22 مارس 2025", false},
		{"French Month Only", "22 Mars 2025", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseDateWithDefaultTZ(tc.dateStr)

			if tc.wantErr {
				if err == nil {
					t.Errorf("ParseDateWithDefaultTZ(%q) succeeded, want error", tc.dateStr)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseDateWithDefaultTZ(%q) failed: %v", tc.dateStr, err)
				return
			}

			// Just verify we got a non-zero time value
			if result.IsZero() {
				t.Errorf("ParseDateWithDefaultTZ(%q) returned zero time", tc.dateStr)
			}

			t.Logf("Successfully parsed %q as %v", tc.dateStr, result)
		})
	}
}

// TestAllDateStrings tests all the date strings from the document
func TestAllDateStrings(t *testing.T) {
	dateStrings := []string{
		"TueAMEETE_RMarchC822",
		"Vie, 30 Sep 2022 21:27:13 -0500",
		"Tue, 18 Mar 2025 5:58:24 PM",
		"Vie, 03/21/2025 - 00:00",
		"SunPMEETE_RJanuaryC822",
		"FriPMEETE_RMarchC822",
		"TuePMEESTE_RJuneC822",
		"ThuAMEETE_RNovemberC822",
		"Dom, 23 Mar 2025 00:05:36 +0000 2025-03-23 00:05:36",
		"Sunday Mar 23 2025 04:37:41",
		"23 Μαρ 2025 13:11:00 +0000",
		"Sat, 22 Mar 2025 19:57:06 TIME_ZONE",
		"Sunday Mar 23 2025 04:50:22",
		"Sunday, March 23, 2025, 16:20 GMT +5:30",
		"23.03.2025 | 08:28",
		"Friday 05 Jul 2024 08:00:00 -0600",
		"Mar 22, 2025, 12:00pm",
		"Sunday Mar 23 2025 13:12:05",
		"March 23, 2025, 4:50 pm",
		"SunPMEETE_RMarchC822",
		"23 Μαρ 2025 13:31:00 +0000",
		"Domenica, 23 Marzo, 2025 - 10:33",
		"Jue, 29 Jun 23 15:34:11 +0200",
		"TueAMEETE_RJanuaryC822",
		"Sunday Mar 23 2025 13:12:05",
		"12:06 23.03.2025",
		"FriPMEETE_RFebruaryC822",
		"Sun, 23 Mar 2025 08:14:46 Europe/Dublin",
		"Mar, 31 Ago 2021 22:46:32 -0500",
		"dom, 23 mar 2025 12:40:59 +0100",
		"Sun, 23 Mar 2025 09:40:10 Europe/Dublin",
		"Sun,23 Mar 2025 18:37:00 +07",
		"23 March 2025 - 12:10",
		"Dom, 23 Mar 2025 00:05:36 +0000 2025-03-23 00:05:36",
		"FriAMEESTE_RSeptemberC822",
		"Sun, 23 Mar 2025 11:27:51 Europe/Dublin",
		"Sáb, 22 Mar 2025 18:22:10 +0000 2025-03-22 18:22:10",
		"Sat, 22 Mar 2025  CST",
		"TuePMEETE_RJanuaryC822",
		"23 March 2025 - 12:10",
		"TuePMEESTE_ROctoberC822",
		"ven, 21 mar 2025 13:49:00 CDT",
		"Sunday Mar 23 2025 13:44:53",
		"SunPMEETE_RMarchC822",
		"23 March 2025 - 12:10",
		"Sun, 23 Mar 2025 11:24:58 Europe/Dublin",
		"Sat, 22 Mar 2025 11:07:41 PM",
		"Sun, 23 Mar 2025 11:55:28 Europe/Dublin",
		"SunPMEETE_RJanuaryC822",
		"Mar 22, 2025, 8:01am",
		"22 Mars 2025",
		"Sun, 23 Mar 2025 11:41:30 Europe/Dublin",
		"23-03-2025 11:15",
		"2025-03-23T11:02:13Z +0300",
		"SunAMEETE_RMarchC822",
		"Sun, 23 March 2025, 05:06:27 PM +0530",
		"22 مارس 2025",
		"Mar 22, 2025, 12:00pm",
		"Sun, 23 Mar 2025 08:14:46 Europe/Dublin",
		"MonAMEETE_RMarchC822",
		"Sun, 23 Mar 2025 08:14:46 Europe/Dublin",
		"Sun, 23 Mar 2025 05:49:21 z",
		"SunPMEETE_RMarchC822",
		"Sun, 23 Mar 2025 11:41:30 Europe/Dublin",
		"Sat, 22 Mar 2025  CST",
		"Sun, 23 Mar 2025 09:25:30 Europe/Dublin",
		"Sunday Mar 23 2025 13:54:16",
		"WedAMEETE_RMarchC822",
		"Sunday Mar 23 2025 13:54:16",
		"12:58 23.03.2025",
		"Sat, 22 Mar 2025 08:36 PM EDT",
		"Sun, 23 Mar 2025 08:11 AM EDT",
		"Sun, 23 Mar 2025 08:19 AM EDT",
		"Sun, 23 Mar 2025 08:16 AM EDT",
		"Thu, 13 Mar 2025 11:23 AM EDT",
		"Fri, 21 Mar 2025 07:52 PM EDT",
		"Fri, 21 Mar 2025 01:36 PM EDT",
		"Sunday Mar 23 2025 04:37:41",
		"Sunday Mar 23 2025 13:54:16",
		"Fri, 21 Mar 2025 12:03 PM EDT",
		"Sat, 22 Mar 2025 09:17 PM EDT",
		"Fri, 21 Mar 2025 07:14 PM EDT",
		"Sun, 23 Mar 2025 08:16 AM EDT",
		"Thu, 20 Mar 2025 01:20 PM EDT",
		"Sat, 22 Mar 2025 08:36 PM EDT",
		"Sat, 22 Mar 2025 08:36 PM EDT",
		" Mon, 5 Oct 2020 19:30:00 GMT ",
		"Sun, 23 Mar 2025 08:11 AM EDT",
		"Sat, 22 Mar 2025 11:07:41 PM",
		"Sat, 22 Mar 2025 10:26 AM EDT",
		"23 March 2025 - 13:00",
		"Sun, 23 Mar 2025 08:00 AM EDT",
		"Sun, 23 Mar 2025 06:00 AM CDT",
		"Sat, 22 Mar 2025 11:17 PM EDT",
		"Sat, 22 Mar 2025 05:41 PM EDT",
		"Sun, 23 Mar 2025 07:35 AM EDT",
		"Mon, 17 Mar 2025 24:15:59 +0530",
	}

	var failCount int
	for i, dateStr := range dateStrings {
		result, err := ParseDateWithDefaultTZ(dateStr)
		if err != nil {
			t.Logf("❌ Failed to parse date %d: %q - Error: %v", i, dateStr, err)
			failCount++
		} else {
			t.Logf("✅ Successfully parsed date %d: %q - Result: %v", i, dateStr, result)
		}
	}

	t.Logf("Parsing summary: %d succeeded, %d failed out of %d total",
		len(dateStrings)-failCount, failCount, len(dateStrings))

	// Test is successful if at least 95% of dates can be parsed
	maxAllowedFailures := len(dateStrings) * 5 / 100
	if failCount > maxAllowedFailures {
		t.Errorf("Too many parsing failures: %d (more than 5%% of total)", failCount)
	}
}
