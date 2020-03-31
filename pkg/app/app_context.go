package app

import "github.com/DWethmar/go-api/pkg/contententry"

type AppContext struct {
	Entries contententry.Service
}

func createAppContext(entryServive contententry.Service) *AppContext {
	return &AppContext{
		Entries: entryServive,
	}
}
