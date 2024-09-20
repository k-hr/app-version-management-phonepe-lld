package versionmanager

import (
	"app-version-management-phonepe-lld/models"
	"reflect"
	"sync"
	"testing"
)

func TestAppVersionManagementSystem_CheckForInstall(t *testing.T) {
	type fields struct {
		Apps       map[string]*models.App
		RolloutMap map[string]string
		mu         sync.Mutex
	}
	type args struct {
		appName         string
		deviceOSVersion string
	}
	var tests []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avms := &AppVersionManagementSystem{
				Apps:       tt.fields.Apps,
				RolloutMap: tt.fields.RolloutMap,
				mu:         tt.fields.mu,
			}
			got, got1 := avms.CheckForInstall(tt.args.appName, tt.args.deviceOSVersion)
			if got != tt.want {
				t.Errorf("CheckForInstall() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CheckForInstall() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestAppVersionManagementSystem_CreateUpdatePatch(t *testing.T) {
	type fields struct {
		Apps       map[string]*models.App
		RolloutMap map[string]string
		mu         sync.Mutex
	}
	type args struct {
		appName     string
		fromVersion string
		toVersion   string
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avms := &AppVersionManagementSystem{
				Apps:       tt.fields.Apps,
				RolloutMap: tt.fields.RolloutMap,
				mu:         tt.fields.mu,
			}
			got, err := avms.CreateUpdatePatch(tt.args.appName, tt.args.fromVersion, tt.args.toVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUpdatePatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUpdatePatch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppVersionManagementSystem_UploadNewVersion(t *testing.T) {
	type fields struct {
		Apps       map[string]*models.App
		RolloutMap map[string]string
		mu         sync.Mutex
	}
	type args struct {
		appName      string
		versionID    string
		minOSVersion string
		fileContent  []byte
		isBeta       bool
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			avms := &AppVersionManagementSystem{
				Apps:       tt.fields.Apps,
				RolloutMap: tt.fields.RolloutMap,
				mu:         tt.fields.mu,
			}
			if err := avms.UploadNewVersion(tt.args.appName, tt.args.versionID, tt.args.minOSVersion, tt.args.fileContent, tt.args.isBeta); (err != nil) != tt.wantErr {
				t.Errorf("UploadNewVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAppVersionManagementSystem(t *testing.T) {
	var tests []struct {
		name string
		want *AppVersionManagementSystem
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAppVersionManagementSystem(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppVersionManagementSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}
