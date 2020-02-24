package config

import (
	"reflect"
	"testing"
)

func TestFromFile(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    *Config
		wantErr bool
	}{
		{"case 01", "./config.json", &Config{Server{"127.0.0.1", "8080"},
			DB{"postgres", "127.0.0.1", "5432", "postgres", "", "articledb"}}, false},
		{"case 02", "./config.yml", &Config{}, true},
		{"case 03", "./config_.json", &Config{}, true},
		{"case 03", "./confi.json", &Config{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromFile(tt.args)

			if err != nil && !tt.wantErr {
				t.Errorf("unable to process config file :%v", err)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromFile() = %s, want %s", got, tt.want)
			}
		})
	}
}
