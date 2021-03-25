package config

import (
	"reflect"
	"testing"
)

func TestGetConfiguration(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    Configuration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetConfiguration(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfiguration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getContentFromFileName(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getContentFromFileName(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("getContentFromFileName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getContentFromFileName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConfigurationFromContent(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Configuration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getConfigurationFromContent(tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfigurationFromContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getConfigurationFromContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
