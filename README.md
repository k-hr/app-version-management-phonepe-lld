App Version Management System

This project is a simplified App Version Management System for handling app versions, installs, updates, and roll-out strategies. The system is designed to work with multiple mobile platforms (Android, iOS) and implements core features like managing versions, creating update patches, and checking compatibility for installs and updates.
Table of Contents

    Overview
    Features
    Directory Structure
    Getting Started
        Prerequisites
        Running the Application
        Running Tests
    Modules
    Tests
    Contributing

Overview

The App Version Management System allows an app owner to directly manage app installations and updates on mobile devices, without needing an external app store (like Google Play Store or Apple App Store). It handles the following scenarios:

    Install: For first-time app installations, ensuring compatibility with the device's OS.
    Update: When a new version is rolled out, it creates a patch between the installed and the latest version.
    Rollout Strategies: Beta releases or percentage rollouts to a limited number of devices.

Features

    Upload New Version: Store a new app version with metadata such as minimum OS version.
    Create Update Patch: Generate a diff (patch) between two app versions for updating existing installations.
    Rollout Strategies: Roll out updates using beta or percentage-based strategies.
    Check Compatibility: Verify whether a device is compatible for installation or update.
    Concurrency Safe: Thread-safe operations using sync.Mutex to handle concurrent version uploads and queries.

Directory Structure

bash

app-version-management/
├── main.go                    # Entry point
├── app_version_manager.go      # Core logic for managing app versions
├── models.go                  # Struct definitions (App, AppVersion)
├── helpers.go                 # Helper methods (createDiffPack, installApp, updateApp)
└── app_version_manager_test.go # Unit tests for the system

Getting Started
Prerequisites

    Go 1.22 or later installed on your system.
    A terminal or IDE for running the project.

Running the Application

    Clone the repository to your local system.

bash

git clone https://github.com/your-username/app-version-management.git

cd app-version-management

    Run the Go application.

bash

go run main.go

This will simulate the app management flow, uploading versions, creating patches, and handling installs.
Running Tests

Run the following command to execute the unit tests:

bash

go test -v

This will run the tests located in app_version_manager_test.go and display the test results.
Modules
1. models.go

Defines the data structures for managing apps and versions:

    App: Represents an app and its versions.
    AppVersion: Stores metadata like version ID, minimum OS version, and file content.

2. app_version_manager.go

Contains the core logic for managing app versions, including:

    UploadNewVersion: Stores a new app version with metadata.
    CreateUpdatePatch: Generates a diff between two versions.
    CheckForInstall: Checks if an installation is possible for a device based on its OS.

3. helpers.go

Contains helper methods to:

    Create a diff pack between two app versions (createDiffPack).
    Simulate app installation (installApp) and updates (updateApp).

Tests

The app_version_manager_test.go file includes unit tests for the core functionalities:

    UploadNewVersion: Tests uploading a new version of an app.
    CreateUpdatePatch: Tests creating a patch between two app versions.
    CheckForInstall: Tests whether the system can determine if an install is possible for a given device.

Run the tests using go test -v to verify that the system behaves as expected.
Contributing

Feel free to open issues or submit pull requests for enhancements, bug fixes, or new features.