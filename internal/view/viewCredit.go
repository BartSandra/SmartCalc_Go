package view

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
)

func (v *View) openCreditCalculator() {
	calcWindow := fyne.CurrentApp().NewWindow("Credit Calculator")
	v.showCreditCalculator(calcWindow)
}

func (v *View) showCreditCalculator(mainWindow fyne.Window) {

	buttonBackgroundColor := color.NRGBA{R: 220, G: 185, B: 240, A: 128}

	amountEntry := widget.NewEntry()
	amountEntry.SetPlaceHolder("Total Loan Amount")

	termEntry := widget.NewEntry()
	termEntry.SetPlaceHolder("Term (months)")

	rateEntry := widget.NewEntry()
	rateEntry.SetPlaceHolder("Interest Rate (%)")

	creditType := widget.NewSelect([]string{"Annuity", "Differentiated"}, func(selected string) {
	})

	monthlyPaymentLabel := widget.NewLabel("Monthly Payment: ")
	overpayLabel := widget.NewLabel("Overpayment: ")
	totalRepaymentLabel := widget.NewLabel("Total Repayment: ")

	// Текст и оформление кнопки "Calculate"
	calculateButtonText := canvas.NewText("Calculate", color.Black)
	calculateButtonText.Alignment = fyne.TextAlignCenter
	calculateButtonText.TextStyle = fyne.TextStyle{Bold: true}

	calculateButtonBackground := canvas.NewRectangle(buttonBackgroundColor)
	calculateButtonBackground.SetMinSize(fyne.NewSize(120, 40))

	calculateButton := widget.NewButton("", func() {
		loanAmount, err := parseInput(amountEntry.Text, "loan amount")
		if err != nil {
			showError(mainWindow, fmt.Sprintf("Invalid loan amount.\nPlease enter a positive number (e.g., 10000, 150000.75)."))
			return
		}

		term, err := parseInput(termEntry.Text, "loan term")
		if err != nil {
			showError(mainWindow, fmt.Sprintf("Invalid loan term.\nPlease enter a positive integer (e.g., 12, 24, 36)."))
			return
		}

		rate, err := parseInput(rateEntry.Text, "interest rate")
		if err != nil {
			showError(mainWindow, fmt.Sprintf("Invalid interest rate.\nPlease enter a positive number (e.g., 5, 12.5, 28.1)."))
			return
		}

		results, err := v.presenter.CalculateCredit(creditType.Selected, loanAmount, term, rate)
		if err != nil {
			showError(mainWindow, err.Error())
			return
		}

		if creditType.Selected == "Annuity" {
			monthlyPaymentLabel.SetText(fmt.Sprintf("Monthly Payment: %.2f", results["MonthlyPayment"]))
		} else {
			monthlyPaymentLabel.SetText(fmt.Sprintf("First Month: %.2f, Last Month: %.2f", results["FirstMonth"], results["LastMonth"]))
		}
		overpayLabel.SetText(fmt.Sprintf("Overpayment: %.2f", results["Overpay"]))
		totalRepaymentLabel.SetText(fmt.Sprintf("Total Repayment: %.2f", results["TotalRepayment"]))
	})

	calculateButtonWithBackground := container.NewStack(
		calculateButton,
		calculateButtonBackground,
		container.NewCenter(calculateButtonText),
	)

	// Текст и оформление кнопки "Clear"
	clearButtonText := canvas.NewText("Clear", color.Black)
	clearButtonText.Alignment = fyne.TextAlignCenter
	clearButtonText.TextStyle = fyne.TextStyle{Bold: true}

	clearButtonBackground := canvas.NewRectangle(buttonBackgroundColor)
	clearButtonBackground.SetMinSize(fyne.NewSize(120, 40))

	clearButton := widget.NewButton("", func() {
		amountEntry.SetText("")
		termEntry.SetText("")
		rateEntry.SetText("")
		monthlyPaymentLabel.SetText("Monthly Payment: ")
		overpayLabel.SetText("Overpayment: ")
		totalRepaymentLabel.SetText("Total Repayment: ")
	})

	clearButtonWithBackground := container.NewStack(
		clearButton,
		clearButtonBackground,
		container.NewCenter(clearButtonText),
	)

	form := container.NewVBox(
		amountEntry,
		termEntry,
		rateEntry,
		creditType,
		calculateButtonWithBackground,
		monthlyPaymentLabel,
		overpayLabel,
		totalRepaymentLabel,
		clearButtonWithBackground,
	)

	mainWindow.SetContent(form)
	mainWindow.Resize(fyne.NewSize(400, 400))
	mainWindow.CenterOnScreen()
	mainWindow.SetFixedSize(true)
	mainWindow.Show()
}

func showError(mainWindow fyne.Window, message string) {
	content := container.NewVBox(widget.NewLabel(message))

	closeButton := widget.NewButton("Close", func() {
		mainWindow.Canvas().Overlays().Remove(mainWindow.Canvas().Overlays().Top())
	})

	content.Add(closeButton)

	dialog := widget.NewPopUp(content, mainWindow.Canvas())
	dialog.Show()
}

func parseInput(input string, fieldName string) (float64, error) {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil || value <= 0 {
		return 0.0, fmt.Errorf("invalid %s", fieldName)
	}
	return value, nil
}
