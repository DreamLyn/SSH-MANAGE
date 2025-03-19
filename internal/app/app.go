package app

import (
	"sync"

	"github.com/pocketbase/pocketbase"
)

var instance *pocketbase.PocketBase

var intanceOnce sync.Once

//确保
func GetApp() *pocketbase.PocketBase {
	intanceOnce.Do(func() {
		instance = pocketbase.NewWithConfig(pocketbase.Config{
			HideStartBanner: true,
		})
	})

	return instance
}