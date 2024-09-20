package main

import (
	"app-version-management-phonepe-lld/versionmanager"
	"fmt"
)

func main() {
	avms := versionmanager.NewAppVersionManagementSystem()

	// Upload versions
	err := avms.UploadNewVersion("PhonePe", "v1.0", "Android-9", []byte("v1.0 content"), false)
	if err != nil {
		return
	}
	err = avms.UploadNewVersion("PhonePe", "v2.0", "Android-10", []byte("v2.0 content"), false)
	if err != nil {
		return
	}

	// Create a diff patch
	patch, _ := avms.CreateUpdatePatch("PhonePe", "v1.0", "v2.0")
	fmt.Printf("Created patch: %v\n", patch)

	// Check install availability
	versionID, available := avms.CheckForInstall("PhonePe", "Android-9")
	if available {
		versionmanager.InstallApp(versionID)
	}
}
