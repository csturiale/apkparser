package apkparser

import "image"

type ApkInfo struct {
	Manifest Manifest
	Icon     image.Image
}

type Instrumentation struct {
	Name            string `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Target          string `xml:"http://schemas.android.com/apk/res/android targetPackage,attr"`
	HandleProfiling bool   `xml:"http://schemas.android.com/apk/res/android handleProfiling,attr"`
	FunctionalTest  bool   `xml:"http://schemas.android.com/apk/res/android functionalTest,attr"`
}

// ActivityAction is an action of an activity.
type ActivityAction struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// ActivityCategory is a category of an activity.
type ActivityCategory struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// ActivityIntentFilter is an int32ent filter of an activity.
type ActivityIntentFilter struct {
	Actions    []ActivityAction   `xml:"action"`
	Categories []ActivityCategory `xml:"category"`
}

// AppActivity is an activity in an application.
type AppActivity struct {
	Theme             string                 `xml:"http://schemas.android.com/apk/res/android theme,attr"`
	Name              string                 `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Label             string                 `xml:"http://schemas.android.com/apk/res/android label,attr"`
	ScreenOrientation string                 `xml:"http://schemas.android.com/apk/res/android screenOrientation,attr"`
	IntentFilters     []ActivityIntentFilter `xml:"intent-filter"`
}

// AppActivityAlias https://developer.android.com/guide/topics/manifest/activity-alias-element
type AppActivityAlias struct {
	Name           string                 `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Label          string                 `xml:"http://schemas.android.com/apk/res/android label,attr"`
	TargetActivity string                 `xml:"http://schemas.android.com/apk/res/android targetActivity,attr"`
	IntentFilters  []ActivityIntentFilter `xml:"intent-filter"`
}

// MetaData is a metadata in an application.
type MetaData struct {
	Name  string `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Value string `xml:"http://schemas.android.com/apk/res/android value,attr"`
}

// Application is an application in an APK.
type Application struct {
	AllowTaskReparenting  bool               `xml:"http://schemas.android.com/apk/res/android allowTaskReparenting,attr"`
	AllowBackup           bool               `xml:"http://schemas.android.com/apk/res/android allowBackup,attr"`
	BackupAgent           string             `xml:"http://schemas.android.com/apk/res/android backupAgent,attr"`
	Debuggable            bool               `xml:"http://schemas.android.com/apk/res/android debuggable,attr"`
	Description           string             `xml:"http://schemas.android.com/apk/res/android description,attr"`
	Enabled               bool               `xml:"http://schemas.android.com/apk/res/android enabled,attr"`
	HasCode               bool               `xml:"http://schemas.android.com/apk/res/android hasCode,attr"`
	HardwareAccelerated   bool               `xml:"http://schemas.android.com/apk/res/android hardwareAccelerated,attr"`
	Icon                  string             `xml:"http://schemas.android.com/apk/res/android icon,attr"`
	KillAfterRestore      bool               `xml:"http://schemas.android.com/apk/res/android killAfterRestore,attr"`
	LargeHeap             bool               `xml:"http://schemas.android.com/apk/res/android largeHeap,attr"`
	Label                 string             `xml:"http://schemas.android.com/apk/res/android label,attr"`
	Logo                  string             `xml:"http://schemas.android.com/apk/res/android logo,attr"`
	ManageSpaceActivity   string             `xml:"http://schemas.android.com/apk/res/android manageSpaceActivity,attr"`
	Name                  string             `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Permission            string             `xml:"http://schemas.android.com/apk/res/android permission,attr"`
	Persistent            bool               `xml:"http://schemas.android.com/apk/res/android persistent,attr"`
	Process               string             `xml:"http://schemas.android.com/apk/res/android process,attr"`
	RestoreAnyVersion     bool               `xml:"http://schemas.android.com/apk/res/android restoreAnyVersion,attr"`
	RequiredAccountType   string             `xml:"http://schemas.android.com/apk/res/android requiredAccountType,attr"`
	RestrictedAccountType string             `xml:"http://schemas.android.com/apk/res/android restrictedAccountType,attr"`
	SupportsRtl           bool               `xml:"http://schemas.android.com/apk/res/android supportsRtl,attr"`
	TaskAffinity          string             `xml:"http://schemas.android.com/apk/res/android taskAffinity,attr"`
	TestOnly              bool               `xml:"http://schemas.android.com/apk/res/android testOnly,attr"`
	Theme                 string             `xml:"http://schemas.android.com/apk/res/android theme,attr"`
	UIOptions             string             `xml:"http://schemas.android.com/apk/res/android uiOptions,attr"`
	VMSafeMode            bool               `xml:"http://schemas.android.com/apk/res/android vmSafeMode,attr"`
	Activities            []AppActivity      `xml:"activity"`
	ActivityAliases       []AppActivityAlias `xml:"activity-alias"`
	MetaData              []MetaData         `xml:"meta-data"`
}

// UsesSDK is target SDK version.
type UsesSDK struct {
	Min    int32 `xml:"http://schemas.android.com/apk/res/android minSdkVersion,attr"`
	Target int32 `xml:"http://schemas.android.com/apk/res/android targetSdkVersion,attr"`
	Max    int32 `xml:"http://schemas.android.com/apk/res/android maxSdkVersion,attr"`
}

// UsesPermission is user grant the system permission.
type UsesPermission struct {
	Name string `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Max  int32  `xml:"http://schemas.android.com/apk/res/android maxSdkVersion,attr"`
}

// Manifest is a manifest of an APK.
type Manifest struct {
	Package                   string           `xml:"package,attr"`
	CompileSDKVersion         int32            `xml:"http://schemas.android.com/apk/res/android compileSdkVersion,attr"`
	CompileSDKVersionCodename string           `xml:"http://schemas.android.com/apk/res/android compileSdkVersionCodename,attr"`
	VersionCode               int32            `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	VersionName               string           `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	App                       Application      `xml:"application"`
	Instrument                Instrumentation  `xml:"instrumentation"`
	SDK                       UsesSDK          `xml:"uses-sdk"`
	UsesPermissions           []UsesPermission `xml:"uses-permission"`
}
