package versionmanager

import (
	"reflect"
	"testing"
)

func TestInstallApp(t *testing.T) {
	type args struct {
		versionID string
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InstallApp(tt.args.versionID)
		})
	}
}

func Test_createDiffPack(t *testing.T) {
	type args struct {
		fromVersion []byte
		toVersion   []byte
	}
	var tests []struct {
		name string
		args args
		want []byte
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createDiffPack(tt.args.fromVersion, tt.args.toVersion); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createDiffPack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateApp(t *testing.T) {
	type args struct {
		diffPack []byte
	}
	var tests []struct {
		name string
		args args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateApp(tt.args.diffPack)
		})
	}
}
