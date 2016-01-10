package plugins

import (
	"github.com/grafana/grafana/pkg/bus"
	m "github.com/grafana/grafana/pkg/models"
)

func GetOrgAppSettings(orgId int64) (map[string]*m.AppSettings, error) {
	query := m.GetAppSettingsQuery{OrgId: orgId}

	if err := bus.Dispatch(&query); err != nil {
		return nil, err
	}

	orgAppsMap := make(map[string]*m.AppSettings)
	for _, orgApp := range query.Result {
		orgAppsMap[orgApp.AppId] = orgApp
	}

	return orgAppsMap, nil
}

func GetEnabledPlugins(orgId int64) (*EnabledPlugins, error) {
	enabledPlugins := NewEnabledPlugins()
	orgApps, err := GetOrgAppSettings(orgId)
	if err != nil {
		return nil, err
	}

	seenPanels := make(map[string]bool)
	seenApi := make(map[string]bool)

	for appType, installedApp := range Apps {
		var app AppPlugin
		app = *installedApp

		// check if the app is stored in the DB for this org and if so, use the
		// state stored there.
		if b, ok := orgApps[appType]; ok {
			app.Enabled = b.Enabled
			app.Pinned = b.Pinned
		}

		// if app.Enabled {
		// 	for _, d := range app.DatasourcePlugins {
		// 		if ds, ok := DataSources[d]; ok {
		// 			enabledPlugins.DataSourcePlugins[d] = ds
		// 		}
		// 	}
		// 	for _, p := range app.PanelPlugins {
		// 		if panel, ok := Panels[p]; ok {
		// 			if _, ok := seenPanels[p]; !ok {
		// 				seenPanels[p] = true
		// 				enabledPlugins.PanelPlugins = append(enabledPlugins.PanelPlugins, panel)
		// 			}
		// 		}
		// 	}
		// 	for _, a := range app.ApiPlugins {
		// 		if api, ok := ApiPlugins[a]; ok {
		// 			if _, ok := seenApi[a]; !ok {
		// 				seenApi[a] = true
		// 				enabledPlugins.ApiPlugins = append(enabledPlugins.ApiPlugins, api)
		// 			}
		// 		}
		// 	}
		// 	enabledPlugins.AppPlugins = append(enabledPlugins.AppPlugins, &app)
		// }
	}

	// add all plugins that are not part of an App.
	for d, installedDs := range DataSources {
		if installedDs.App == "" {
			enabledPlugins.DataSources[d] = installedDs
		}
	}

	for p, panel := range Panels {
		if panel.App == "" {
			if _, ok := seenPanels[p]; !ok {
				seenPanels[p] = true
				enabledPlugins.Panels = append(enabledPlugins.Panels, panel)
			}
		}
	}

	for a, api := range ApiPlugins {
		if api.App == "" {
			if _, ok := seenApi[a]; !ok {
				seenApi[a] = true
				enabledPlugins.ApiList = append(enabledPlugins.ApiList, api)
			}
		}
	}

	return &enabledPlugins, nil
}
