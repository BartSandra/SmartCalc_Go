package view

import (
	"sandra/APG2_SmartCalc_v3.0_Desktop_Go-1/src/internal/helpers"
)

func (v *View) handleEInput() {
	v.presenter.HandleEInput()
}

func (v *View) UpdatedisplayLabelWithText(inputText string) {
	v.counter = len(inputText)
	v.displayLabel.SetText(inputText)
}

func (v *View) UpdateXLabelWithText(inputText string) {
	if !helpers.IsValidInput(inputText) {
		v.UpdatedisplayLabelWithText(inputText)
	} else {
		v.variableXLabel.SetText(inputText)
		v.UpdatedisplayLabelWithText("0")
	}
}

func (v *View) appendOperator(operator string) {
	v.presenter.AppendOperator(operator)
}

func (v *View) appendButtonText(inputText string) {
	v.presenter.AppendButtonText(inputText)
}

func (v *View) resetButton() {
	v.presenter.ResetButton()
}

func (v *View) addDecimalPoint() {
	v.presenter.AddDecimalPoint()
}

func (v *View) deleteButton() {
	v.presenter.DeleteButton()
}

func (v *View) appendX() {
	v.presenter.AppendX()
}

func (v *View) initializeXButton() {
	v.presenter.InitializeXButton()
}

func (v *View) evaluateExpression() {
	v.presenter.EvaluateAndProcessExpression()
}

func (v *View) inverseSign() {
	v.presenter.InverseSign()
}

func (v *View) GetCounter() int {
	return v.counter
}

func (v *View) GetVariableXLabel() string {
	return v.variableXLabel.Text
}

func (v *View) GetDisplayLabel() string {
	return v.displayLabel.Text
}
