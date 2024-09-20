package main

import (
	"app-version-management-phonepe-lld/versionmanager"
	"fmt"
)

func main() {
	avms := versionmanager.NewAppVersionManagementSystem()

	// Sample device IDs
	deviceIDs := []string{"device1", "device2", "device3", "device4", "device5", "device6", "device7", "device8"}

	// Upload versions
	err := avms.UploadNewVersion("PhonePe", "v2.0", "Android-10", []byte("v2.0 content"), false)
	if err != nil {
		return
	}

	// Rollout version 2.0 to 50% of the devices using the percentage strategy
	err = avms.RolloutVersion("PhonePe", "v2.0", "percentage", 50, deviceIDs)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Percentage rollout completed")
	}

	// Rollout version 2.0 to specific beta devices
	betaDevices := []string{"device1", "device2"}
	err = avms.RolloutVersion("PhonePe", "v2.0", "beta", 0, betaDevices)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Beta rollout completed")
	}
}
