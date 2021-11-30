package configs

import (
	_ "encoding/json"
	_ "io"
	_ "io/ioutil"
	"os"
	"proxy_checker/types"
)

func New() *types.Configs {
	return &types.Configs {
		SiteAddress: getEnv("SITE_ADDRESS", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}