package sync

import (
	"encoding/json"
	"log"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gandarfh/httui/internal/repositories/offline"
	"github.com/gandarfh/httui/internal/services"
	"github.com/gandarfh/httui/pkg/common"
	"github.com/gandarfh/httui/pkg/convert"
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
		}()

		locally, _ := offline.NewRequest().ListForSync()

		httuiClient := services.HttuiApiDatasource

		if !d.LastRequestSync.IsZero() {
			httuiClient = httuiClient.Params("lastUpdate", d.LastRequestSync.Format(time.RFC3339))
		}

		response, err := httuiClient.Get("/requests")
		if err != nil {
			log.Println("Jaum", err.Error())
			return nil
		}

		data := []map[string]any{}
		err = response.Decode(&data)

		if err != nil {
			log.Println("Jaum", err.Error())
			return nil
		}

		log.Println("On", len(data))
		remotes := ProcessRemoteRequests(data)

		upsertLocal := func(r offline.Request, exists bool) (offline.Request, error) {
			sync := true
			r.Sync = &sync
			if offline.NewRequest().Sql.Model(&r).Session(&gorm.Session{FullSaveAssociations: true}).Where("external_id = ?", r.ExternalId).Updates(&r).RowsAffected == 0 {
				offline.NewRequest().Sql.Model(&r).Create(&r)
			}

			return r, nil
		}

		upsertRemote := func(r offline.Request, exists bool) (offline.Request, error) {
			if r.Type == "" {
				r.Type = offline.REQUEST
			}

			payload := map[string]any{}
			convert.ToSource(r, &payload)
			log.Println("Start UpsertRemote", r.Name)

			var parentId *string
			if r.ParentID != nil {
				parentRequest, err := offline.NewRequest().FindOne(*r.ParentID)
				if err == nil {
					if parentRequest.ExternalId != "" {
						parentId = &parentRequest.ExternalId
						payload["parentId"] = parentId
					}
				} else {
					payload["parentId"] = nil
				}
			}

			body, _ := json.Marshal(payload)
			client, err := services.HttuiApiDatasource.Body(body).Post("/requests")
			if err != nil {
				return r, nil
			}

			res := map[string]any{}
			if err := client.Decode(&res); err != nil {
				return r, err
			}

			if res["parentId"] != nil {
				p := res["parentId"]

				parentRequest := new(offline.Request)
				err := offline.NewRequest().Sql.
					Model(parentRequest).
					Where("external_id = ?", p).
					First(parentRequest).Error

				if err != nil {
					res["parentId"] = nil
				} else {
					res["parentId"] = &parentRequest.ID
				}
			}

			convert.ToSource(res, &r)

			sync := true
			r.Sync = &sync

			offline.NewRequest().
				Sql.Model(&r).
				Session(&gorm.Session{FullSaveAssociations: true}).
				Where("id = ?", r.ID).
				Updates(&r)

			offline.NewDefault().Update(offline.Default{LastRequestSync: r.UpdatedAt.Add(1 * time.Minute)})

			return r, nil
		}

		// Sync locally to remote
		ExecuteOfflineRequests(locally, remotes, upsertLocal, upsertRemote)

		// Sync remote to locally
		remotes = ProcessRemoteRequests(data)
		SyncDatabases([]offline.Request{}, remotes, upsertLocal, upsertRemote)
		return nil
	}
}

func ProcessRemoteRequests(remotes []map[string]any) []offline.Request {
	requests := []offline.Request{}
	for _, remote := range remotes {
		if remote["parentId"] != nil {
			p := remote["parentId"]

			parentRequest := new(offline.Request)
			err := offline.NewRequest().Sql.
				Model(parentRequest).
				Where("external_id = ?", p).
				First(parentRequest).Error

			if err != nil {
				remote["parentId"] = nil
			} else {
				remote["parentId"] = &parentRequest.ID
			}
		}

		request := new(offline.Request)
		convert.ToSource(remote, request)

		requests = append(requests, *request)
	}

	log.Println(len(remotes), len(requests))

	return requests
}

func GenerateRequestTree(requests []offline.Request) map[*uint][]offline.Request {
	tree := map[*uint][]offline.Request{}

	for _, req := range requests {
		tree[req.ParentID] = append(tree[req.ParentID], req)
	}

	return tree
}

func ExecuteOfflineRequests(locally, remotes []offline.Request, upsertLocal, upsertRemote func(data offline.Request, exists bool) (offline.Request, error)) {
	locallyTree := GenerateRequestTree(locally)

	keys := []*uint{}
	for key := range locallyTree {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i] == nil {
			return true
		}
		if keys[j] == nil {
			return false
		}
		return *keys[i] < *keys[j]
	})

	// Sync locally to remote
	for _, parent := range keys {
		if parent == nil {
			SyncDatabases(locallyTree[parent], remotes, upsertLocal, upsertRemote)
		} else {
			SyncDatabases(locallyTree[parent], []offline.Request{}, upsertLocal, upsertRemote)
		}
	}
}
