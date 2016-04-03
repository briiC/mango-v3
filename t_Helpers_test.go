package mango

import (
	"fmt"
	"testing"
	"time"
)

// Parsing datetimes
func Test_ToTime(t *testing.T) {
	dtNow := time.Now() // to check current values

	cases := map[string]string{
		// y-m-d
		"07-02":               dtNow.Format("2006") + "-07-02 00:00:00 +0000 UTC",
		"07-02 07:01":         dtNow.Format("2006") + "-07-02 07:01:00 +0000 UTC",
		"2015-07-02":          "2015-07-02 00:00:00 +0000 UTC",
		"2015-11-19 23:47":    "2015-11-19 23:47:00 +0000 UTC",
		"2015-11-19 23:47:58": "2015-11-19 23:47:58 +0000 UTC",

		// LV
		"02.07":               dtNow.Format("2006") + "-07-02 00:00:00 +0000 UTC",
		"02.07 07:01":         dtNow.Format("2006") + "-07-02 07:01:00 +0000 UTC",
		"02.07.2015":          "2015-07-02 00:00:00 +0000 UTC",
		"19.11.2015 23:47":    "2015-11-19 23:47:00 +0000 UTC",
		"19.11.2015 23:47:58": "2015-11-19 23:47:58 +0000 UTC",

		// USA
		"07/02":               dtNow.Format("2006") + "-07-02 00:00:00 +0000 UTC",
		"07/02 07:01":         dtNow.Format("2006") + "-07-02 07:01:00 +0000 UTC",
		"07/02/2015":          "2015-07-02 00:00:00 +0000 UTC",
		"11/19/2015 23:47":    "2015-11-19 23:47:00 +0000 UTC",
		"11/19/2015 23:47:58": "2015-11-19 23:47:58 +0000 UTC",

		// Time only
		"23:47":    dtNow.Format("2006-01-02") + " 23:47:00 +0000 UTC",
		"23:47:58": dtNow.Format("2006-01-02") + " 23:47:58 +0000 UTC",
	}

	for in, expected := range cases {
		dt, err := toTime(in)
		if err != nil || dt.String() != expected {
			fmt.Println(":: ", in, " ---> ", dt)
			t.Fatal("\tEXPECTED: [", expected, "]", err)
		}
	}

}
