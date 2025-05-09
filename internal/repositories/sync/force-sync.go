package sync

import (
	"log"

	"github.com/gandarfh/httui/internal/config"
	"github.com/gandarfh/httui/internal/repositories/offline"
)

func SyncDatabases[T offline.Entity](locally []T, remotes []T, UpsertLocal, UpsertRemote func(data T, exists bool) (T, error)) error {
	localIndex := BuildIndex(locally)
	remoteIndex := BuildIndex(remotes)

	log.Println("localIndex", len(localIndex))

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
          log.Println(err.Error())
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
		key := data[i].GetExternalID()

		if key == "" {
			key = data[i].GetID()
		}

		index[key] = data[i]
	}
	return index
}
