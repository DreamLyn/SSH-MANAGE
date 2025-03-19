package migrations

import (
	"os"
	"strings"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// add up queries...
		superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}

		records, err := app.FindAllRecords(superusers)
		if err != nil {
			return err
		}

		if len(records) == 0 {
			envUsername := strings.TrimSpace(os.Getenv("SSH_MANAGE_ADMIN_USERNAME"))
			if envUsername == "" {
				envUsername = "lyn@dreamlyn.cn"
			}

			envPassword := strings.TrimSpace(os.Getenv("SSH_MANAGE_ADMIN_PASSWORD"))
			if envPassword == "" {
				envPassword = "1234567890"
			}

			record := core.NewRecord(superusers)
			record.Set("email", envUsername)
			record.Set("password", envPassword)
			return app.Save(record)
		}

		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}
