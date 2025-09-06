package main

import (
	"testing"
	"time"
)

func TestAkkMensaApiImpl_GetAvailableDates(t *testing.T) {
	type fields struct {
		baseUrl string
	}
	tests := []struct {
		name      string
		fields    fields
		checkFunc func([]time.Time) bool
		wantErr   bool
	}{
		{
			name:   "Test not empty with real URL",
			fields: fields{baseUrl: "https://mensa.akk.org/json"},
			checkFunc: func(times []time.Time) bool {
				return len(times) > 0
			},
			wantErr: false,
		},
		{
			name:    "Test if times plausible with real URL",
			fields:  fields{baseUrl: "https://mensa.akk.org/json"},
			wantErr: false,
			checkFunc: func(times []time.Time) bool {
				for _, t := range times {
					if t.Year() < 2020 || t.Year() > 2030 {
						return false
					}
				}
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &AkkMensaApiImpl{
				baseUrl: tt.fields.baseUrl,
			}
			got, err := api.GetAvailableDates()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvailableDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc != nil && !tt.checkFunc(got) {
				t.Errorf("GetAvailableDates() got = %v, checkFunc failed", got)
			}
		})
	}
}

func TestAkkMensaApiImpl_GetMenuForDate(t *testing.T) {
	type fields struct {
		baseUrl string
	}
	type args struct {
		date time.Time
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		checkFunc func(map[string]interface{}) bool
		wantErr   bool
	}{
		{
			name:   "Test valid date with real URL if response not empty",
			fields: fields{baseUrl: "https://mensa.akk.org/json"},
			args:   args{date: time.Date(2020, 07, 10, 0, 0, 0, 0, time.UTC)},
			checkFunc: func(menu map[string]interface{}) bool {
				return len(menu) > 0
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := &AkkMensaApiImpl{
				baseUrl: tt.fields.baseUrl,
			}
			got, err := api.GetMenuForDate(tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMenuForDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFunc != nil && !tt.checkFunc(got) {
				t.Errorf("GetMenuForDate() got = %v, checkFunc failed", got)
			}
		})
	}
}
