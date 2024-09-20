package models

// AppVersion stores metadata about an app version.
type AppVersion struct {
	VersionID       string
	MinOSVersion    string
	FileContent     []byte
	IsBeta          bool
	ReleasedDevices []string
}

// App represents an app and its versions.
type App struct {
	Name     string
	Versions map[string]*AppVersion // Maps VersionID to AppVersion
}
