package versionmanager

import (
	"app-version-management-phonepe-lld/models"
	"errors"
	"sync"
)

// AppVersionManagementSystem manages app installs, updates, and rollouts.
type AppVersionManagementSystem struct {
	Apps       map[string]*models.App
	RolloutMap map[string]string // Maps deviceID to App VersionID
	mu         sync.Mutex        // Ensures concurrency safety
}

// NewAppVersionManagementSystem initializes a new AppVersionManagementSystem.
func NewAppVersionManagementSystem() *AppVersionManagementSystem {
	return &AppVersionManagementSystem{
		Apps:       make(map[string]*models.App),
		RolloutMap: make(map[string]string),
	}
}

// Upload a new version for an app.
func (avms *AppVersionManagementSystem) UploadNewVersion(appName string, versionID string, minOSVersion string, fileContent []byte, isBeta bool) error {
	avms.mu.Lock()
	defer avms.mu.Unlock()

	// Check if the app exists
	app, exists := avms.Apps[appName]
	if !exists {
		app = &models.App{
			Name:     appName,
			Versions: make(map[string]*models.AppVersion),
		}
		avms.Apps[appName] = app
	}

	app.Versions[versionID] = &models.AppVersion{
		VersionID:    versionID,
		MinOSVersion: minOSVersion,
		FileContent:  fileContent,
		IsBeta:       isBeta,
	}
	return nil
}

// Create an update patch between two versions.
func (avms *AppVersionManagementSystem) CreateUpdatePatch(appName string, fromVersion string, toVersion string) ([]byte, error) {
	avms.mu.Lock()
	defer avms.mu.Unlock()

	app, exists := avms.Apps[appName]
	if !exists {
		return nil, errors.New("app does not exist")
	}

	fromVer, fromExists := app.Versions[fromVersion]
	toVer, toExists := app.Versions[toVersion]

	if !fromExists || !toExists {
		return nil, errors.New("one of the versions does not exist")
	}

	return createDiffPack(fromVer.FileContent, toVer.FileContent), nil
}

// Check if an install is possible for a device.
func (avms *AppVersionManagementSystem) CheckForInstall(appName string, deviceOSVersion string) (string, bool) {
	app, exists := avms.Apps[appName]
	if !exists {
		return "", false
	}

	var latestVersionID string
	for versionID, version := range app.Versions {
		if version.MinOSVersion <= deviceOSVersion {
			if latestVersionID == "" || versionID > latestVersionID {
				latestVersionID = versionID
			}
		}
	}

	if latestVersionID != "" {
		return latestVersionID, true
	}
	return "", false
}
