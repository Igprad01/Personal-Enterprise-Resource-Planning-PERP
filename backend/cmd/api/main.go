package main

import (
	"log"

	"personal-erp-backend/internal/application"
	"personal-erp-backend/internal/config"
	"personal-erp-backend/internal/database"
	"personal-erp-backend/internal/handler"
	"personal-erp-backend/internal/repository"
	"personal-erp-backend/internal/router"
)

func main() {
	config := config.Load()

	db := must(database.Connect(config))
	defer db.Close()
	mustNoErr(database.Migrate(db))
	
	categoryRepo := repository.NewCategoryRepository(db)
	activityRepo := repository.NewActivityRepository(db)

	categoryHandler := handler.NewCategoryHandler(application.NewCategoryService(categoryRepo))
	activityHandler := handler.NewActivityHandler(application.NewActivityService(activityRepo, categoryRepo))

	engine := router.New(
		handler.NewHealthHandler(db),
		categoryHandler,
		activityHandler,
	)

	log.Printf("server listening on :%s", config.ServerPort)
	if err := engine.Run(":" + config.ServerPort); err != nil {
		log.Fatal(err)
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		log.Fatalf("%v", err)
	}
	return v
}

func mustNoErr(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}