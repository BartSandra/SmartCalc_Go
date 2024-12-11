package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
)

func (v *View) openHelp() {
	helpWindow := fyne.CurrentApp().NewWindow("Help")
	v.showHelp(helpWindow)
}

func (v *View) showHelp(mainWindow fyne.Window) {

	err := godotenv.Load(".env")
	if err != nil {
		showErrorHelp(fyne.CurrentApp().NewWindow("Error"), "Error loading .env file")
		return
	}

	helpFilePath := os.Getenv("HELP_FILE_PATH")
	if helpFilePath == "" {
		showErrorHelp(fyne.CurrentApp().NewWindow("Error"), "HELP_FILE_PATH is not set in .env file")
		return
	}

	content, err := ioutil.ReadFile(helpFilePath)
	if err != nil {
		showErrorHelp(fyne.CurrentApp().NewWindow("Error"), fmt.Sprintf("Error reading help file: %v", err))
		return
	}

	helpText := string(content)
	markdownText := widget.NewRichTextFromMarkdown(helpText)

	scrollContainer := container.NewScroll(markdownText)

	mainWindow.SetContent(scrollContainer)

	mainWindow.Resize(fyne.NewSize(800, 700))
	mainWindow.CenterOnScreen()
	mainWindow.Show()
}

func showErrorHelp(mainWindow fyne.Window, message string) {
	content := container.NewVBox(widget.NewLabel(message))

	dialog := widget.NewPopUp(content, mainWindow.Canvas())
	dialog.Show()
}
