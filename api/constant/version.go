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
	PlatformCodeSDKFlutterEnum     = "fl-sdk"
	PlatformCodeSDKIOSEnum         = "ios-sdk"
	PlatformCodeSDKReactEnum       = "rt-sdk"
	PlatformCodeSDKReactNativeEnum = "rn-sdk"
	PlatformCodeSDKWebEnum         = "web-sdk"
)

const FeatureSDKSourcePlatformVersionData = `{
	"m2cm": {
		"chat": {
			"an-sdk": 1010,
			"fl-sdk": 9999,
			"ios-sdk":  1012,
			"rt-sdk": 15,
			"rn-sdk": 7,
			"web-sdk": 15
		}
	}
}`

type FeatureSDKSourcePlatformVersion struct {
	FeatureM2CMEnum SDKSource `json:"m2cm"`
}

type SDKSource struct {
	SDKSourceChat SDKVersions `json:"chat"`
	SDKSourcefeed SDKVersions `json:"feed"`
}

type SDKVersions struct {
	PlatformCodeSDKAndroidEnum     int `json:"an-sdk"`
	PlatformCodeSDKFlutterEnum     int `json:"fl-sdk"`
	PlatformCodeSDKIOSEnum         int `json:"ios-sdk"`
	PlatformCodeSDKReactEnum       int `json:"rt-sdk"`
	PlatformCodeSDKReactNativeEnum int `json:"rn-sdk"`
	PlatformCodeSDKWebEnum         int `json:"web-sdk"`
}

func ParseFeatureSDKSourcePlatformVersionData() *FeatureSDKSourcePlatformVersion {
	var featureSDKSourcePlatformVersion FeatureSDKSourcePlatformVersion
	err := json.Unmarshal([]byte(FeatureSDKSourcePlatformVersionData), &featureSDKSourcePlatformVersion)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed Struct: %+v\n", featureSDKSourcePlatformVersion)
	return &featureSDKSourcePlatformVersion
}

func VersionCheck(sdkSource string, platformCode string, versionCode int, apiVersion int, feature string) bool {
	featureSDKSourcePlatformVersionCodes := ParseFeatureSDKSourcePlatformVersionData()

	switch feature {
	case FeatureM2CMEnum:
		switch sdkSource {
		case SDKSourceChatEnum:
			switch platformCode {
			case PlatformCodeSDKAndroidEnum:
				return versionCode >= featureSDKSourcePlatformVersionCodes.FeatureM2CMEnum.SDKSourceChat.PlatformCodeSDKAndroidEnum
			case PlatformCodeSDKFlutterEnum:
				return versionCode >= featureSDKSourcePlatformVersionCodes.FeatureM2CMEnum.SDKSourceChat.PlatformCodeSDKFlutterEnum
			case PlatformCodeSDKIOSEnum:
				return versionCode >= featureSDKSourcePlatformVersionCodes.FeatureM2CMEnum.SDKSourceChat.PlatformCodeSDKIOSEnum
			case PlatformCodeSDKReactEnum:
				return versionCode >= featureSDKSourcePlatformVersionCodes.FeatureM2CMEnum.SDKSourceChat.PlatformCodeSDKReactEnum
			case PlatformCodeSDKReactNativeEnum:
				return versionCode >= featureSDKSourcePlatformVersionCodes.FeatureM2CMEnum.SDKSourceChat.PlatformCodeSDKReactEnum
			case PlatformCodeSDKWebEnum:
				return versionCode >= featureSDKSourcePlatformVersionCodes.FeatureM2CMEnum.SDKSourceChat.PlatformCodeSDKReactEnum
			default:
				log.Printf("err=unknown platform code, platform code=%s", platformCode)
				return false
			}
		default:
			log.Printf("err=unknown sdk source, sdk source=%s", sdkSource)
			return false
		}
	default:
		log.Printf("err=unknown feature, feature=%s", feature)
		return false
	}
}
