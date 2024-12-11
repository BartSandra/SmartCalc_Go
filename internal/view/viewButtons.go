package view

import (
	"fyne.io/fyne/v2/canvas"
	"image/color"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/joho/godotenv"
	"sandra/APG2_SmartCalc_v3.0_Desktop_Go-1/src/internal/presenter"
)

const (
	WindowTitle   = "SmartCalc_Go"
	DefaultNumber = "0"
)

type ButtonConfig struct {
	Label    string
	Callback func()
}

type View struct {
	mainWindow      fyne.Window
	displayLabel    *widget.Label
	variableXLabel  *widget.Label
	variableLabel   *widget.Label
	presenter       *presenter.Presenter
	historyFilePath string
	counter         int
	useScientific   bool
}

func (v *View) GetUseScientific() bool {
	return v.useScientific
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func (v *View) GetHistoryFilePath() string {
	return os.Getenv("HISTORY_FILE_PATH")
}

func NewCalculatorView(myApp fyne.App) *View {

	historyFilePath := os.Getenv("HISTORY_FILE_PATH")

	view := &View{
		mainWindow:      myApp.NewWindow(WindowTitle),
		displayLabel:    widget.NewLabel(DefaultNumber),
		variableXLabel:  widget.NewLabel(DefaultNumber),
		variableLabel:   widget.NewLabel("x:"),
		historyFilePath: historyFilePath,
		counter:         1,
	}

	view.variableLabel.TextStyle = fyne.TextStyle{Bold: true}
	view.variableXLabel.TextStyle = fyne.TextStyle{Italic: true}
	view.displayLabel.TextStyle = fyne.TextStyle{Bold: false, Italic: true}
	view.displayLabel.Alignment = fyne.TextAlignTrailing

	view.mainWindow.SetContent(view.createCalculatorLayout())
	view.mainWindow.Resize(fyne.NewSize(445, 250))
	view.mainWindow.SetFixedSize(true)
	view.mainWindow.Show()

	return view
}

func (v *View) InitPresenter(p *presenter.Presenter) {
	v.presenter = p
}

func (v *View) createCalculatorLayout() *fyne.Container {
	variableBox := container.NewHBox(v.variableLabel, v.variableXLabel)

	buttonBox := container.NewHBox(
		v.createButtonColumn(v.getButtonColumnConfig0(), color.NRGBA{R: 220, G: 185, B: 240, A: 128}),
		v.createButtonColumn(v.getButtonColumnConfig1(), color.NRGBA{R: 220, G: 185, B: 240, A: 128}),
		v.createButtonColumn(v.getButtonColumnConfig2(), color.NRGBA{R: 220, G: 185, B: 240, A: 128}),
		v.createButtonColumn(v.getButtonColumnConfig3(), color.NRGBA{R: 220, G: 185, B: 240, A: 128}),
		v.createButtonColumn(v.getButtonColumnConfig4(), color.NRGBA{R: 220, G: 185, B: 240, A: 128}),
		v.createButtonColumn(v.getButtonColumnConfig5(), color.NRGBA{R: 220, G: 185, B: 240, A: 128}),
		v.createButtonColumn(v.getButtonColumnConfig6(), color.NRGBA{R: 220, G: 185, B: 240, A: 128}),
		v.createButtonColumn(v.getButtonColumnConfig7(), color.NRGBA{R: 220, G: 185, B: 240, A: 128}),
	)

	scrollVariableBox := container.NewHScroll(variableBox)
	scrollVariableBox.SetMinSize(fyne.NewSize(350, 40))
	scrollDisplayLabel := container.NewHScroll(v.displayLabel)
	scrollDisplayLabel.SetMinSize(fyne.NewSize(350, 40))

	return container.NewVBox(scrollDisplayLabel, scrollVariableBox, buttonBox)
}

func (v *View) createButtonColumn(configs []ButtonConfig, bgColor color.Color) *fyne.Container {
	var buttons []fyne.CanvasObject
	for _, config := range configs {
		background := canvas.NewRectangle(bgColor)
		background.SetMinSize(fyne.NewSize(60, 40))

		clickableButton := widget.NewButton("", config.Callback)

		buttonText := canvas.NewText(config.Label, color.Black)
		buttonText.Alignment = fyne.TextAlignCenter
		buttonText.TextStyle = fyne.TextStyle{Bold: true}

		buttonWithBackground := container.NewStack(
			clickableButton,
			background,
			container.NewCenter(buttonText),
		)

		buttonWithBackground.Resize(fyne.NewSize(60, 40))
		buttons = append(buttons, buttonWithBackground)
	}

	return container.NewVBox(buttons...)
}

func (v *View) getButtonColumnConfig0() []ButtonConfig {
	return []ButtonConfig{
		{"Plot", v.openPlot},
		{" Help ", v.openHelp},
		{"History", v.openHistory},
		{"Credit", v.openCreditCalculator},
		{"e", func() { v.appendOperator("e") }},
	}
}

func (v *View) getButtonColumnConfig1() []ButtonConfig {
	return []ButtonConfig{
		{"(", func() { v.appendButtonText("(") }},
		{"ln", func() { v.appendButtonText("ln(") }},
		{"atan", func() { v.appendButtonText("atan(") }},
		{"acos", func() { v.appendButtonText("acos(") }},
		{"asin", func() { v.appendButtonText("asin(") }},
	}
}

func (v *View) getButtonColumnConfig2() []ButtonConfig {
	return []ButtonConfig{
		{")", func() { v.appendButtonText(")") }},
		{" log ", func() { v.appendButtonText("log(") }},
		{" tan ", func() { v.appendButtonText("tan(") }},
		{" cos ", func() { v.appendButtonText("cos(") }},
		{" sin ", func() { v.appendButtonText("sin(") }},
	}
}

func (v *View) getButtonColumnConfig3() []ButtonConfig {
	return []ButtonConfig{
		{"<-", v.deleteButton},
		{"   ^    ", func() { v.appendOperator("^") }},
		{"sqrt", func() { v.appendButtonText("sqrt(") }},
		{"pi", func() { v.appendButtonText("pi") }},
		{"x<-", v.initializeXButton},
	}
}

func (v *View) getButtonColumnConfig4() []ButtonConfig {
	return []ButtonConfig{
		{"  AC  ", v.resetButton},
		{"7", func() { v.appendButtonText("7") }},
		{"4", func() { v.appendButtonText("4") }},
		{"1", func() { v.appendButtonText("1") }},
		{"x", v.appendX},
	}
}

func (v *View) getButtonColumnConfig5() []ButtonConfig {
	return []ButtonConfig{
		{"  +/-  ", v.inverseSign},
		{"8", func() { v.appendButtonText("8") }},
		{"5", func() { v.appendButtonText("5") }},
		{"2", func() { v.appendButtonText("2") }},
		{"0", func() { v.appendButtonText("0") }},
	}
}

func (v *View) getButtonColumnConfig6() []ButtonConfig {
	return []ButtonConfig{
		{"   %   ", func() { v.appendOperator("%") }},
		{"9", func() { v.appendButtonText("9") }},
		{"6", func() { v.appendButtonText("6") }},
		{"3", func() { v.appendButtonText("3") }},
		{"   .   ", v.addDecimalPoint},
	}
}

func (v *View) getButtonColumnConfig7() []ButtonConfig {
	return []ButtonConfig{
		{"   /    ", func() { v.appendOperator("/") }},
		{"*", func() { v.appendOperator("*") }},
		{"-", func() { v.appendButtonText("-") }},
		{"+", func() { v.appendOperator("+") }},
		{"=", v.evaluateExpression},
	}
}
