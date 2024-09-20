package main

import (
	"errors"
	"fmt"
	"sync"
)

// AppVersion Basic structure for AppVersion to store metadata
type AppVersion struct {
	VersionID       string
	MinOSVersion    string
	FileContent     []byte
	IsBeta          bool
	ReleasedDevices []string // Device IDs for Beta releases
}

// App Structure for App, managing versions
type App struct {
	Name     string
	Versions map[string]*AppVersion // Maps VersionID to AppVersion
}

// AppVersionManagementSystem manages app installs, updates, and rollouts
type AppVersionManagementSystem struct {
	Apps       map[string]*App
	RolloutMap map[string]string // Map from deviceID to App VersionID
	mu         sync.Mutex        // Ensure concurrency safety
}

// NewAppVersionManagementSystem Initialize a new system
func NewAppVersionManagementSystem() *AppVersionManagementSystem {
	return &AppVersionManagementSystem{
		Apps:       make(map[string]*App),
		RolloutMap: make(map[string]string),
	}
}

// Upload a new version for an app
func (avms *AppVersionManagementSystem) uploadNewVersion(appName string, versionID string, minOSVersion string, fileContent []byte, isBeta bool) error {
	avms.mu.Lock()
	defer avms.mu.Unlock()

	// Check if the app exists
	app, exists := avms.Apps[appName]
	if !exists {
		// Create a new app if it doesn't exist
		app = &App{
			Name:     appName,
			Versions: make(map[string]*AppVersion),
		}
		avms.Apps[appName] = app
	}

	// Add the new version
	app.Versions[versionID] = &AppVersion{
		VersionID:    versionID,
		MinOSVersion: minOSVersion,
		FileContent:  fileContent,
		IsBeta:       isBeta,
	}

	return nil
}

// Create an update patch between two versions
func (avms *AppVersionManagementSystem) createUpdatePatch(appName string, fromVersion string, toVersion string) ([]byte, error) {
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

	// Use createDiffPack capability (assuming it exists)
	return createDiffPack(fromVer.FileContent, toVer.FileContent), nil
}

// Rollout a version to devices based on a strategy
func (avms *AppVersionManagementSystem) releaseVersion(appName string, versionID string, strategy string, deviceIDs []string) error {
	avms.mu.Lock()
	defer avms.mu.Unlock()

	app, exists := avms.Apps[appName]
	if !exists {
		return errors.New("app does not exist")
	}

	version, versionExists := app.Versions[versionID]
	if !versionExists {
		return errors.New("version does not exist")
	}

	switch strategy {
	case "beta":
		version.ReleasedDevices = deviceIDs
	case "percentage":
		percentage := len(deviceIDs) * 10 / 100
		version.ReleasedDevices = deviceIDs[:percentage]
	default:
		return errors.New("unknown rollout strategy")
	}

	return nil
}

// Check if a specific app version supports a given device
func (avms *AppVersionManagementSystem) isAppVersionSupported(appName string, versionID string, deviceOSVersion string) bool {
	avms.mu.Lock()
	defer avms.mu.Unlock()

	app, exists := avms.Apps[appName]
	if !exists {
		return false
	}

	version, versionExists := app.Versions[versionID]
	if !versionExists {
		return false
	}

	return version.MinOSVersion <= deviceOSVersion
}

// Check if a fresh installation is supported on a device, returns the latest versionID and a bool
func (avms *AppVersionManagementSystem) checkForInstall(appName string, deviceOSVersion string) (string, bool) {
	// Remove the mutex lock here as it's already locked in executeTask
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

// Check if an update is available for a device
func (avms *AppVersionManagementSystem) checkForUpdates(appName string, currentVersionID string, deviceOSVersion string) (string, bool) {
	avms.mu.Lock()
	defer avms.mu.Unlock()

	app, exists := avms.Apps[appName]
	if !exists {
		return "", false
	}

	for versionID, version := range app.Versions {
		if versionID > currentVersionID && version.MinOSVersion <= deviceOSVersion {
			return versionID, true
		}
	}

	return "", false
}

// Execute a task - install or update
func (avms *AppVersionManagementSystem) executeTask(taskType string, appName string, deviceID string, currentVersionID string) error {
	avms.mu.Lock()

	if _, exists := avms.Apps[appName]; !exists {
		avms.mu.Unlock() // Unlock here to prevent deadlock
		return errors.New("app does not exist")
	}

	avms.mu.Unlock() // Unlock before making the function calls

	switch taskType {
	case "install":
		versionID, available := avms.checkForInstall(appName, deviceID)
		if available {
			installApp(versionID)
		}
	case "update":
		newVersionID, updateAvailable := avms.checkForUpdates(appName, currentVersionID, deviceID)
		if updateAvailable {
			patch, err := avms.createUpdatePatch(appName, currentVersionID, newVersionID)
			if err != nil {
				return err
			}
			updateApp(patch)
		}
	default:
		return errors.New("invalid task type")
	}

	return nil
}

// Dummy function to represent diff pack creation
func createDiffPack(fromVersion []byte, toVersion []byte) []byte {
	// Implement real diff logic here
	return append(fromVersion, toVersion...)
}

// Dummy function to represent installing an app
func installApp(versionID string) {
	fmt.Printf("Installing app version: %s\n", versionID)
}

// Dummy function to represent updating an app
func updateApp(diffPack []byte) {
	fmt.Println("Updating app with diff pack")
}

func main() {
	avms := NewAppVersionManagementSystem()

	// Upload versions
	err := avms.uploadNewVersion("PhonePe", "v1.0", "Android-9", []byte("v1.0 content"), false)
	if err != nil {
		return
	}
	err = avms.uploadNewVersion("PhonePe", "v2.0", "Android-10", []byte("v2.0 content"), false)
	if err != nil {
		return
	}

	// Create a diff patch
	patch, _ := avms.createUpdatePatch("PhonePe", "v1.0", "v2.0")
	fmt.Printf("Created patch: %v\n", patch)

	// Rollout a version
	err = avms.releaseVersion("PhonePe", "v2.0", "beta", []string{"device1", "device2"})
	if err != nil {
		return
	}

	// Execute an install task
	err = avms.executeTask("install", "PhonePe", "device1", "")
	if err != nil {
		return
	}
}
