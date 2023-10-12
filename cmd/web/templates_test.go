package main

import (
	"testing"
	"time"

	"github.com/wenkanglu/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		inp  time.Time
		exp  string
	}{
		{
			name: "UTC",
			inp:  time.Date(2023, 3, 17, 10, 15, 0, 0, time.UTC),
			exp:  "17 Mar 2023 at 10:15",
		},
		{
			name: "Empty",
			inp:  time.Time{},
			exp:  "",
		},
		{
			name: "CET",
			inp:  time.Date(2023, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			exp:  "17 Mar 2023 at 09:15",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.inp)

			assert.Equal(t, hd, tt.exp)
		})
	}
}
