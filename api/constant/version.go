package constant

import (
	"encoding/json"
	"fmt"
	"log"
)

const (
	FeatureM2CMEnum = "m2cm"
)

const (
	SDKSourceChatEnum = "chat"
	SDKSourceFeedEnum = "feed"
)

const (
	PlatformCodeSDKAndroidEnum     = "an-sdk"
	PlatformCodeSDKIOSEnum         = "ios-sdk"
	PlatformCodeSDKReactEnum       = "rt-sdk"
	PlatformCodeSDKReactNativeEnum = "rn-sdk"
	PlatformCodeSDKwebEnum         = "web-sdk"
)

const FeatureSDKSourcePlatformVersionData = `{
	"m2cm": {
		"chat": {
			"an-sdk": 1,
			"ios-sdk": 2
		},
		"feed": {
			"an-sdk": 1,
			"ios-sdk": 2
		}
	}
}`

type FeatureSDKSourcePlatformVersion struct {
	M2CM SDKSource `json:"m2cm"`
}

type SDKSource struct {
	Chat SDKVersions `json:"chat"`
	Feed SDKVersions `json:"feed"`
}

type SDKVersions struct {
	PlatformCodeSDKAndroidEnum int `json:"an-sdk"`
	PlatformCodeSDKIOSEnum     int `json:"ios-sdk"`
}

func ParseFeatureSDKSourcePlatformVersionData() FeatureSDKSourcePlatformVersion {
	var featureSDKSourcePlatformVersion FeatureSDKSourcePlatformVersion
	err := json.Unmarshal([]byte(FeatureSDKSourcePlatformVersionData), &featureSDKSourcePlatformVersion)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed Struct: %+v\n", featureSDKSourcePlatformVersion)
	return featureSDKSourcePlatformVersion
}

var FeatureSDKSourcePlatformVersionCodes = ParseFeatureSDKSourcePlatformVersionData()
