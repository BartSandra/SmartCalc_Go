package view

import (
	"fmt"
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (v *View) openHistory() {
	historyWindow := fyne.CurrentApp().NewWindow("History")
	v.showHistory(historyWindow)
}

func (v *View) showHistory(mainWindow fyne.Window) {
	content, err := os.ReadFile(v.historyFilePath)
	if err != nil {
		dialog.ShowError(fmt.Errorf("Failed to read history file: %v", err), mainWindow)
		return
	}

	historyLines := strings.Fields(string(content))
	if len(historyLines) == 0 {
		historyLines = []string{"No history available."}
	}

	historyList := widget.NewList(
		func() int {
			return len(historyLines)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("", nil)
		},
		func(index widget.ListItemID, obj fyne.CanvasObject) {
			button := obj.(*widget.Button)
			button.SetText(historyLines[index])
			button.Importance = widget.LowImportance

			if historyLines[index] != "No history available." {
				button.OnTapped = func() {
					v.updateDisplayLabelWithHistory(historyLines[index])
					mainWindow.Close()
				}
			} else {
				button.Disable()
			}
		},
	)

	scrollContainer := container.NewScroll(historyList)
	scrollContainer.SetMinSize(fyne.NewSize(400, 380))

	buttonBackgroundColor := color.NRGBA{R: 220, G: 185, B: 240, A: 128}

	clearButtonText := canvas.NewText("Clear History", color.Black)
	clearButtonText.Alignment = fyne.TextAlignCenter
	clearButtonText.TextStyle = fyne.TextStyle{Bold: true}

	clearButtonBackground := canvas.NewRectangle(buttonBackgroundColor)
	clearButtonBackground.SetMinSize(fyne.NewSize(120, 40))

	clearButton := widget.NewButton("", func() {
		dialog.ShowConfirm("Confirm", "Are you sure you want to clear the history?", func(confirmed bool) {
			if confirmed {
				err := os.WriteFile(v.historyFilePath, []byte{}, 0644)
				if err != nil {
					dialog.ShowError(fmt.Errorf("Failed to clear history: %v", err), mainWindow)
				} else {
					v.showHistory(mainWindow)
				}
			}
		}, mainWindow)
	})

	clearButtonWithBackground := container.NewStack(
		clearButton,
		clearButtonBackground,
		container.NewCenter(clearButtonText),
	)

	contentContainer := container.NewVBox(
		scrollContainer,
		clearButtonWithBackground,
	)

	mainWindow.SetContent(contentContainer)
	mainWindow.Resize(fyne.NewSize(400, 400))
	mainWindow.CenterOnScreen()
	mainWindow.SetFixedSize(true)
	mainWindow.Show()
}

func (v *View) updateDisplayLabelWithHistory(historyItem string) {
	v.displayLabel.SetText(historyItem)
}
