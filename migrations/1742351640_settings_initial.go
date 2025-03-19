package migrations

import (
	"github.com/dreamlyn/ssh-manage/internal/domain"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		settings, err := app.FindCollectionByNameOrId(domain.CollectionNameSettings)
		if err != nil {
			println(err.Error())
			return err
		}
		//添加默认主题
		records, _ := app.FindRecordsByFilter(
			domain.CollectionNameSettings, // collection
			"name = 'theme'",              // filter
			"",                            // sort
			10,                            // limit
			0,                             // offset
			dbx.Params{},                  // optional filter params
		)
		if len(records) == 0 {
			record := core.NewRecord(settings)
			record.Set("name", "theme")
			record.Set("content", "dark")
			app.Save(record)
		}
		//添加默认语言
		records, _ = app.FindRecordsByFilter(
			domain.CollectionNameSettings, // collection
			"name = 'language'",           // filter
			"",                            // sort
			10,                            // limit
			0,                             // offset
			dbx.Params{},                  // optional filter params
		)
		if len(records) == 0 {
			record := core.NewRecord(settings)
			record.Set("name", "language")
			record.Set("content", "zh")
			app.Save(record)
		}
		return nil
	}, func(app core.App) error {
		// add down queries...

		return nil
	})
}
