package versionmanager

import (
	"app-version-management-phonepe-lld/models"
	"crypto/sha256"
	"errors"
	"fmt"
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

// UploadNewVersion Upload a new version for an app
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

// CreateUpdatePatch Create an update patch between two versions
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

	return CreateDiffPack(fromVer.FileContent, toVer.FileContent), nil
}

// RolloutVersion rolls out a new version of an app using different strategies, including percentage rollout
func (avms *AppVersionManagementSystem) RolloutVersion(appName string, versionID string, strategy string, percentage int, deviceIDs []string) error {
	avms.mu.Lock()
	defer avms.mu.Unlock()

	// Fetch the app and version
	app, exists := avms.Apps[appName]
	if !exists {
		return fmt.Errorf("app does not exist")
	}

	version, versionExists := app.Versions[versionID]
	if !versionExists {
		return fmt.Errorf("version does not exist")
	}

	switch strategy {
	case "percentage":
		rolloutCount := len(deviceIDs) * percentage / 100
		for _, deviceID := range avms.selectDevicesByPercentage(deviceIDs, rolloutCount) {
			// Update RolloutMap with the new version for the selected devices
			avms.RolloutMap[deviceID] = versionID
			version.ReleasedDevices = append(version.ReleasedDevices, deviceID)
		}
	case "beta":
		// Add a fixed list of devices to receive the update as part of beta testing
		for _, deviceID := range deviceIDs {
			avms.RolloutMap[deviceID] = versionID
			version.ReleasedDevices = append(version.ReleasedDevices, deviceID)
		}
	default:
		return fmt.Errorf("unknown rollout strategy")
	}

	return nil
}

// selectDevicesByPercentage selects a deterministic subset of deviceIDs for percentage rollout.
func (avms *AppVersionManagementSystem) selectDevicesByPercentage(deviceIDs []string, rolloutCount int) []string {
	selectedDevices := []string{}
	for _, deviceID := range deviceIDs {
		if avms.shouldSelectDevice(deviceID) && len(selectedDevices) < rolloutCount {
			selectedDevices = append(selectedDevices, deviceID)
		}
	}
	return selectedDevices
}

// shouldSelectDevice uses a hash function to deterministically decide if a device is selected for the rollout.
func (avms *AppVersionManagementSystem) shouldSelectDevice(deviceID string) bool {
	hash := sha256.Sum256([]byte(deviceID))
	// Use the first byte of the hash to get a value between 0-255, then select based on a range.
	// This allows deterministic selection of a device
	return hash[0] < 128 // Choose a device if its hash falls in the lower half (or tweak based on percentage)
}

// CheckForInstall Check if an installation is possible for a device
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
