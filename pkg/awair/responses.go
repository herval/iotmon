package awair

import "time"

type DevicesResponse struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Name         string  `json:"name"`
	MacAddress   string  `json:"macAddress"`
	Latitude     float64 `json:"latitude"`
	Preference   string  `json:"preference"`
	Timezone     string  `json:"timezone"`
	RoomType     string  `json:"roomType"`
	DeviceType   string  `json:"deviceType"`
	Longitude    float64 `json:"longitude"`
	SpaceType    string  `json:"spaceType"`
	DeviceUUID   string  `json:"deviceUUID"`
	DeviceId     int     `json:"deviceId"`
	LocationName string  `json:"locationName"`
}

type RawDataResponse struct {
	Data []*RawDataPoints `json:"data"`
}

type RawDataPoints struct {
	DeviceId  int       // not available in the original json
	Timestamp time.Time `json:"timestamp"`
	Score     float64   `json:"score"`
	Sensors   []struct {
		Comp  string  `json:"comp"`
		Value float64 `json:"value"`
	} `json:"sensors"`
	Indices []struct {
		Comp  string  `json:"comp"`
		Value float64 `json:"value"`
	} `json:"indices"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type LocalDeviceConfig struct {
	DeviceId   int
	DeviceUuid string `json:"device_uuid"`
	WifiMac    string `json:"wifi_mac"`
	Ssid       string `json:"ssid"`
	Ip         string `json:"ip"`
	Netmask    string `json:"netmask"`
	Gateway    string `json:"gateway"`
	FwVersion  string `json:"fw_version"`
	Timezone   string `json:"timezone"`
	Display    string `json:"display"`
	Led        struct {
		Mode       string `json:"mode"`
		Brightness int    `json:"brightness"`
	} `json:"led"`
	VocFeatureSet int `json:"voc_feature_set"`
}
