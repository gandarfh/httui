package sync

import (
	"encoding/json"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/config"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/internal/services"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/styles"
	"gorm.io/gorm"
)

func SyncDatabases[T offline.Entity](locally []T, remotes []T, UpsertLocal, UpsertRemote func(data T, exists bool) (T, error)) error {
	localIndex := BuildIndex(locally)
	remoteIndex := BuildIndex(remotes)

	if config.Config.Settings.AutoSync.BeforeOpen.Locally {
		for _, remote := range remotes {
			local, exists := localIndex[remote.GetExternalID()]

			if !exists || remote.GetUpdatedAt().After(local.GetUpdatedAt()) {
				localUpdated, err := UpsertLocal(remote, exists)
				if err != nil {
					return err
				}
				localIndex[localUpdated.GetID()] = localUpdated
			}

		}
	}

	if config.Config.Settings.AutoSync.BeforeOpen.Remote {
		for _, local := range localIndex {
			remote, exists := remoteIndex[local.GetExternalID()]
			if !exists || local.GetUpdatedAt().After(remote.GetUpdatedAt()) {
				_, err := UpsertRemote(local, exists)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func BuildIndex[T offline.Entity](data []T) map[string]T {
	index := make(map[string]T)
	for i := range data {
		index[data[i].GetExternalID()] = data[i]
	}
	return index
}

func SyncWorkspaces(program *tea.Program) tea.Cmd {
	return func() tea.Msg {
		d, err := offline.NewDefault().First()

		if err != nil {
			log.Println("SyncWorkspaces:", err.Error())
			return nil
		}

		program.Send(common.SetLoading(true, lipgloss.NewStyle().Bold(true).Foreground(styles.DefaultTheme.PrimaryText).Render(" sync"))())

		defer func() {
			program.Send(common.SetLoading(false)())
			program.Send(common.ResetWorkspace()())
		}()

		locally, _ := offline.NewWorkspace().ListForSync()

		httuiClient := services.HttuiApiDatasource
		if !d.LastWorkspaceSync.IsZero() {
			httuiClient = httuiClient.Params("lastUpdate", d.LastWorkspaceSync.Format(time.RFC3339))
		}

		response, err := httuiClient.Get("/workspaces")
		if err != nil {

			return nil
		}

		remotes := []offline.Workspace{}
		err = response.Decode(&remotes)

		if err != nil {
			return nil
		}

		upsertLocal := func(w offline.Workspace, exists bool) (offline.Workspace, error) {
			sync := true
			w.Sync = &sync
			if offline.NewWorkspace().Sql.Model(&w).Session(&gorm.Session{FullSaveAssociations: true}).Where("external_id = ?", w.ExternalId).Updates(&w).RowsAffected == 0 {
				offline.NewWorkspace().Sql.Model(&w).Create(&w)
			}

			return w, nil
		}

		upsertRemote := func(w offline.Workspace, exists bool) (offline.Workspace, error) {
			body, _ := json.Marshal(w)
			client, err := services.HttuiApiDatasource.Body(body).Post("/workspaces")
			if err != nil {
				return w, nil
			}

			if err := client.Decode(&w); err != nil {
				return w, err
			}

			sync := true
			w.Sync = &sync

			offline.NewWorkspace().
				Sql.Model(&w).
				Session(&gorm.Session{FullSaveAssociations: true}).
				Where("id = ?", w.ID).
				Updates(&w)

			offline.NewDefault().Update(&offline.Default{LastWorkspaceSync: time.Now()})

			return w, nil
		}

		SyncDatabases(locally, remotes, upsertLocal, upsertRemote)
		return nil
	}
}

func SyncRequests(program *tea.Program) tea.Cmd {
	return func() tea.Msg {
		d, err := offline.NewDefault().First()
		if err != nil {
			log.Println("SyncRequests:", err.Error())
			return nil
		}

		program.Send(common.SetLoading(true, lipgloss.NewStyle().Bold(true).Foreground(styles.DefaultTheme.PrimaryText).Render(" sync"))())

		defer func() {
			program.Send(common.SetLoading(false)())
			program.Send(common.ResetWorkspace()())
		}()

		locally, _ := offline.NewRequest().ListForSync()

		httuiClient := services.HttuiApiDatasource
		if !d.LastRequestSync.IsZero() {
			httuiClient = httuiClient.Params("lastUpdate", d.LastRequestSync.Format(time.RFC3339))
		}

		response, err := httuiClient.Get("/requests")
		if err != nil {

			return nil
		}

		remotes := []offline.Request{}
		err = response.Decode(&remotes)

		if err != nil {
			return nil
		}

		upsertLocal := func(r offline.Request, exists bool) (offline.Request, error) {
			sync := true
			r.Sync = &sync
			if offline.NewRequest().Sql.Model(&r).Session(&gorm.Session{FullSaveAssociations: true}).Where("external_id = ?", r.ExternalId).Updates(&r).RowsAffected == 0 {
				offline.NewRequest().Sql.Model(&r).Create(&r)
			}

			return r, nil
		}

		upsertRemote := func(r offline.Request, exists bool) (offline.Request, error) {
			body, _ := json.Marshal(r)
			client, err := services.HttuiApiDatasource.Body(body).Post("/requests")
			if err != nil {
				return r, nil
			}

			if err := client.Decode(&r); err != nil {
				return r, err
			}

			sync := true
			r.Sync = &sync

			offline.NewRequest().
				Sql.Model(&r).
				Session(&gorm.Session{FullSaveAssociations: true}).
				Where("id = ?", r.ID).
				Updates(&r)

			offline.NewDefault().Update(&offline.Default{LastRequestSync: time.Now()})

			return r, nil
		}

		SyncDatabases(locally, remotes, upsertLocal, upsertRemote)
		return nil
	}
}
