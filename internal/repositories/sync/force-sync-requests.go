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

			offline.NewDefault().Update(&offline.Default{LastRequestSync: r.UpdatedAt.Add(1 * time.Minute)})

			return r, nil
		}

		SyncDatabases(locally, remotes, upsertLocal, upsertRemote)
		return nil
	}
}
