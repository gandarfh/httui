package sync

import (
	"encoding/json"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/internal/services"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/styles"
	"gorm.io/gorm"
)

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

			offline.NewDefault().Update(&offline.Default{LastWorkspaceSync: w.UpdatedAt.Add(1 * time.Minute)})

			return w, nil
		}

		SyncDatabases(locally, remotes, upsertLocal, upsertRemote)
		return nil
	}
}
