package main

import (
	"app-version-management-phonepe-lld/versionmanager"
	"fmt"
	"time"
)

func main() {

	start := time.Now()
	avms := versionmanager.NewAppVersionManagementSystem()

	// Sample device IDs
	deviceIDs := []string{"device1", "device2", "device3", "device4", "device5", "device6", "device7", "device8"}

	// Upload versions
	err := avms.UploadNewVersion("PhonePe", "v1.0", "Android-9", []byte("v1.0 content"), false)
	if err != nil {
		return
	}
	err = avms.UploadNewVersion("PhonePe", "v2.0", "Android-10", []byte("v2.0 content"), false)
	if err != nil {
		return
	}

	// Check for updates
	currentVersion := "v1.0"
	deviceOSVersion := "Android-10"
	updateVersion, available := avms.CheckForUpdates("PhonePe", currentVersion, deviceOSVersion)
	if available {
		fmt.Printf("An update is available! New version: %s\n", updateVersion)
	} else {
		fmt.Println("No update available.")
	}

	// Rollout version 2.0 to 50% of the devices using the percentage strategy
	err = avms.ReleaseVersion("PhonePe", "v2.0", "percentage", 50, deviceIDs)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Percentage rollout completed")
	}

	// Rollout version 2.0 to specific beta devices
	betaDevices := []string{"device1", "device2"}
	err = avms.ReleaseVersion("PhonePe", "v2.0", "beta", 0, betaDevices)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Beta rollout completed")
	}

	// End the timer just before the program exits
	elapsed := time.Since(start)
	fmt.Printf("\n⏱️ Total execution time: %s\n", elapsed)
}
