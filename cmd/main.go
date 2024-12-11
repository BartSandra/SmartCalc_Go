package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/joho/godotenv"
	"sandra/APG2_SmartCalc_v3.0_Desktop_Go-1/src/internal/model"
	"sandra/APG2_SmartCalc_v3.0_Desktop_Go-1/src/internal/presenter"
	"sandra/APG2_SmartCalc_v3.0_Desktop_Go-1/src/internal/view"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка при загрузке .env файла: %v", err)
	}

	appInstance := app.New()

	modelPath := getModelPath()

	modelInstance, err := model.NewModel(modelPath)
	if err != nil {
		log.Fatalf("Ошибка при создании модели: %v", err)
	}

	viewCalc := view.NewCalculatorView(appInstance)

	presenter := presenter.NewPresenter(viewCalc, modelInstance)
	viewCalc.InitPresenter(presenter)

	appInstance.Run()
}

func getModelPath() string {
	modelPath := os.Getenv("MODEL_PATH")
	if modelPath == "" {
		modelPath = "./internal/model/model/model.so"
	}
	return modelPath
}
