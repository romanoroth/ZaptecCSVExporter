package main

import (
	"testing"
	"time"
)

func TestCheckIfHighFare(t *testing.T) {
	type args struct {
		timeToCheck time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Check if High Fare when time is 18:59",
			args: args{
				timeToCheck: time.Date(2022, 12, 27, 18, 59, 00, 00, time.Local),
			},
			want: true,
		},
		{
			name: "Check if Low Fare when time is 19:00",
			args: args{
				timeToCheck: time.Date(2022, 12, 27, 19, 00, 00, 01, time.Local),
			},
			want: false,
		},
		{
			name: "Check if High Fare when time is 07:00",
			args: args{
				timeToCheck: time.Date(2022, 12, 27, 07, 00, 00, 01, time.Local),
			},
			want: true,
		},
		{
			name: "Check if High Fare when time is October 07:15",
			args: args{
				timeToCheck: time.Date(2022, 10, 24, 07, 15, 00, 01, time.Local),
			},
			want: true,
		},
		{
			name: "Check if Low Fare when time is 06:59",
			args: args{
				timeToCheck: time.Date(2022, 12, 27, 06, 59, 00, 01, time.Local),
			},
			want: false,
		},
		{
			name: "Check if Low Fare when it is Sunday 12:00",
			args: args{
				timeToCheck: time.Date(2022, 12, 18, 12, 00, 00, 01, time.Local),
			},
			want: false,
		},
		{
			name: "Check if High Fare when it is Wintertime 12:00",
			args: args{
				timeToCheck: time.Date(2022, 03, 26, 12, 00, 00, 01, time.Local),
			},
			want: true,
		},
		{
			name: "Check if High Fare when it is Summerime 12:00",
			args: args{
				timeToCheck: time.Date(2022, 04, 28, 12, 00, 00, 01, time.Local),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckIfHighFare(tt.args.timeToCheck); got != tt.want {
				t.Errorf("CheckIfHighFare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertAndFixTimeZone(t *testing.T) {
	type args struct {
		timestamp string
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
		name: "Check fix time zone 2023-02-26T12:44:10.282+00:00",
		args: args{
			timestamp: "2023-02-26T12:44:10.282+00:00",
		},
		want: time.Date(2023, 02, 26, 13, 44, 10, 282000000, time.Local),
		},
		
		{
			name: "Check fix time zone 2023-02-26T12:44:10.282+00:00",
			args: args{
				timestamp: "2022-10-26T12:44:10.282+00:00",
			},
			want: time.Date(2022, 10, 26, 14, 44, 10, 282000000, time.Local),
			},
	}
	for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := ConvertAndFixTimeZone(tt.args.timestamp)
            if !got.Equal(tt.want) {
                t.Errorf("ConvertAndFixTimeZone() = %v, want %v", got, tt.want)
            }
        })
	}
}
