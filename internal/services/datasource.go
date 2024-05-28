package services

import (
	"github.com/gandarfh/httui/internal/config"
	"github.com/gandarfh/httui/pkg/client/v2"
)

var HttuiApiDatasource client.Client

func DatasourceStart() {
	if config.Config.Settings.Local {

		HttuiApiDatasource = client.New().URL("http://localhost:5000")
	} else {
		HttuiApiDatasource = client.New().URL("https://api.httui.com")
	}

	if config.Config.Settings.Token != "" {
		HttuiApiDatasource = HttuiApiDatasource.
			Auth("Authorization", "Bearer "+config.Config.Settings.Token)
	}
}
