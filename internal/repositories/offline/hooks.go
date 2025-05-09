package offline

import (
	"encoding/json"
	"log"

	"github.com/gandarfh/httui/internal/config"
	"github.com/gandarfh/httui/internal/services"
	"github.com/gandarfh/httui/pkg/convert"
	"gorm.io/gorm"
)

func (w *Workspace) AfterSave(tx *gorm.DB) (err error) {
	if config.Config.Settings.AutoSync.BeforeCreate.Remote {
		if w.Sync == nil || !*w.Sync {
			go func() {
				body, _ := json.Marshal(w)
				res, err := services.HttuiApiDatasource.Body(body).Post("workspaces")
				if err != nil {
					log.Println("Error to Sync request", err.Error())
					return
				}

				if res.StatusCode < 300 {
					res.Decode(w)
					sync := true
					w.Sync = &sync

					NewWorkspace().Sql.
						Session(&gorm.Session{FullSaveAssociations: true}).
						Where("id = ?", w.ID).
						Updates(w)
				}
			}()
		}
	}

	return nil
}

func (r *Request) AfterSave(tx *gorm.DB) (err error) {
	if config.Config.Settings.AutoSync.BeforeCreate.Remote {
		if r.Sync == nil || !*r.Sync {
			go func() {
				var parentId *string
				if r.ParentID != nil {
					parentRequest, _ := NewRequest().FindOne(*r.ParentID)
					parentId = &parentRequest.ExternalId
				}

				payload := map[string]any{}
				convert.ToSource(r, &payload)
				payload["parentId"] = parentId

				body, _ := json.Marshal(payload)
				res, err := services.HttuiApiDatasource.Body(body).Post("requests")
				if err != nil {
					log.Println("Error to Sync request", err.Error())
					return
				}

				if res.StatusCode < 300 {
					res.Decode(r)
					sync := true
					r.Sync = &sync

					NewRequest().Sql.
						Session(&gorm.Session{FullSaveAssociations: true}).
						Where("id = ?", r.ID).
						Updates(r)
				}
			}()
		}
	}

	return nil
}

func (r *Request) AfterDelete(tx *gorm.DB) (err error) {
	if config.Config.Settings.AutoSync.BeforeDelete.Remote {
		go func() {
			body, _ := json.Marshal(r)
			services.HttuiApiDatasource.Body(body).Delete("/requests/" + r.ExternalId)
		}()
	}

	return nil
}

func (r *Response) AfterSave(tx *gorm.DB) (err error) {
	if config.Config.Settings.AutoSync.BeforeCreate.Remote {
		if r.Sync == nil || !*r.Sync {
			go func() {
				body, _ := json.Marshal(r)
				res, err := services.HttuiApiDatasource.Body(body).Post("responses")
				if err != nil {
					log.Println("Error to Sync response", err.Error())
					return
				}

				if res.StatusCode < 300 {
					res.Decode(r)
					sync := true
					r.Sync = &sync

					NewResponse().Sql.
						Session(&gorm.Session{FullSaveAssociations: true}).
						Where("id = ?", r.ID).
						Updates(r)
				}
			}()
		}
	}

	return nil
}
