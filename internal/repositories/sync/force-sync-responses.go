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

func SyncResponses(program *tea.Program) tea.Cmd {
	return func() tea.Msg {
		d, err := offline.NewDefault().First()
		if err != nil {
			log.Println("SyncResponses:", err.Error())
			return nil
		}

		program.Send(common.SetLoading(true, lipgloss.NewStyle().Bold(true).Foreground(styles.DefaultTheme.PrimaryText).Render(" sync"))())

		defer func() {
			program.Send(common.SetLoading(false)())
		}()

		locally, _ := offline.NewResponse().ListForSync()

		httuiClient := services.HttuiApiDatasource
		if !d.LastResponseSync.IsZero() {
			httuiClient = httuiClient.Params("lastUpdate", d.LastResponseSync.Format(time.RFC3339))
		}

		response, err := httuiClient.Get("/responses")
		if err != nil {
			return nil
		}

		remotes := []offline.Response{}
		err = response.Decode(&remotes)

		if err != nil {
			return nil
		}

		upsertLocal := func(r offline.Response, exists bool) (offline.Response, error) {
			sync := true
			r.Sync = &sync
			if offline.NewResponse().Sql.Model(&r).Session(&gorm.Session{FullSaveAssociations: true}).Where("external_id = ?", r.RequestExternalId).Updates(&r).RowsAffected == 0 {
				offline.NewResponse().Sql.Model(&r).Create(&r)
			}

			return r, nil
		}

		upsertRemote := func(r offline.Response, exists bool) (offline.Response, error) {
			body, _ := json.Marshal(r)
			client, err := services.HttuiApiDatasource.Body(body).Post("/responses")
			if err != nil {
				return r, nil
			}

			if err := client.Decode(&r); err != nil {
				return r, err
			}

			sync := true
			r.Sync = &sync

			offline.NewResponse().
				Sql.Model(&r).
				Session(&gorm.Session{FullSaveAssociations: true}).
				Where("id = ?", r.ID).
				Updates(&r)

			offline.NewDefault().Update(&offline.Default{LastResponseSync: r.UpdatedAt.Add(1 * time.Minute)})

			return r, nil
		}

		SyncDatabases(locally, remotes, upsertLocal, upsertRemote)
		return nil
	}
}
