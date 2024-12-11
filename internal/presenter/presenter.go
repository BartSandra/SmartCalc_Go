package presenter

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"unicode"

	"sandra/APG2_SmartCalc_v3.0_Desktop_Go-1/src/internal/helpers"
	"sandra/APG2_SmartCalc_v3.0_Desktop_Go-1/src/internal/model"
)

const (
	eLiteral = "e"
	piValue  = "3.14159265359"
	eValue   = "2.71828182846"
)

type ViewInterface interface {
	UpdatedisplayLabelWithText(inputText string)
	UpdateXLabelWithText(inputText string)
	GetUseScientific() bool
	GetHistoryFilePath() string
	GetCounter() int
	GetVariableXLabel() string
	GetDisplayLabel() string
}

type Presenter struct {
	view  ViewInterface
	model *model.Model
}

func NewPresenter(v ViewInterface, m *model.Model) *Presenter {
	return &Presenter{
		view:  v,
		model: m,
	}
}

func (p *Presenter) EvaluateExpression(expression *string, xValue string) {
	p.view.UpdatedisplayLabelWithText(p.calculateResult(expression, xValue))
}

func (p *Presenter) EvaluateWithX(expression *string, xValue string) {
	p.view.UpdateXLabelWithText(fmt.Sprint(p.calculateResult(expression, xValue)))
}

func (p *Presenter) calculateResult(expression *string, xValue string) string {
	var result string
	if res, err := p.model.Calculate(expression, xValue); err != nil {
		result = err.Error()
	} else {
		if p.view.GetUseScientific() {
			result = strconv.FormatFloat(res, 'e', 8, 64)
		} else {
			result = strconv.FormatFloat(res, 'f', -1, 64)
		}
	}
	return result
}

func (p *Presenter) CalculatePlotResult(expression *string, xValue string) (string, error) {
	if expression == nil || *expression == "" {
		return "", fmt.Errorf("expression is empty")
	}

	res, err := p.model.Calculate(expression, xValue)
	if err != nil {
		return "", fmt.Errorf("calculation error at x=%s: %v", xValue, err)
	}

	if !helpers.IsValidInput(fmt.Sprint(res)) {
		return "", fmt.Errorf("invalid result (NaN, Inf, or invalid input) at x=%s", xValue)
	}

	bigRes := new(big.Float).SetFloat64(res)
	bigRes.SetPrec(200)

	if bigRes.String() == "NaN" || bigRes.String() == "Inf" || bigRes.String() == "-Inf" {
		return "", fmt.Errorf("invalid result after big.Float conversion at x=%s", xValue)
	}

	return bigRes.Text('f', 15), nil
}

func (p *Presenter) CalculateCredit(creditType string, sum, duration, rate float64) (map[string]float64, error) {
	switch creditType {
	case "Annuity":
		monthly, overpay, total, err := p.model.CreditAnnuity(sum, duration, rate)
		if err != nil {
			return nil, err
		}
		return map[string]float64{
			"MonthlyPayment": monthly,
			"Overpay":        overpay,
			"TotalRepayment": total,
		}, nil
	case "Differentiated":
		first, last, overpay, total, err := p.model.CreditDifferentiated(sum, duration, rate)
		if err != nil {
			return nil, err
		}
		return map[string]float64{
			"FirstMonth":     first,
			"LastMonth":      last,
			"Overpay":        overpay,
			"TotalRepayment": total,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported credit type")
	}
}

func (p *Presenter) HandleEInput() {
	currentDisplay := p.view.GetDisplayLabel()
	if p.view.GetCounter() > 0 && unicode.IsDigit(rune(currentDisplay[p.view.GetCounter()-1])) {
		p.view.UpdatedisplayLabelWithText(currentDisplay + eLiteral)
		return
	}

	if currentDisplay[p.view.GetCounter()-1] == '+' || currentDisplay[p.view.GetCounter()-1] == '-' ||
		currentDisplay[p.view.GetCounter()-1] == '*' || currentDisplay[p.view.GetCounter()-1] == '/' ||
		currentDisplay[p.view.GetCounter()-1] == '(' {
		p.view.UpdatedisplayLabelWithText(currentDisplay + eLiteral)
		return
	}

	if unicode.IsLetter(rune(currentDisplay[p.view.GetCounter()-1])) {
		p.view.UpdatedisplayLabelWithText("error")
		return
	}

	p.view.UpdatedisplayLabelWithText("error")
}

func (p *Presenter) AppendOperator(operator string) {
	if p.view.GetCounter() < 256 {
		currentDisplay := p.view.GetDisplayLabel()

		if operator == "pi" {
			if p.view.GetCounter() > 0 && unicode.IsLetter(rune(currentDisplay[p.view.GetCounter()-1])) {
				return
			}

			if p.view.GetCounter() > 1 && p.view.GetCounter() <= len(currentDisplay) && unicode.IsDigit(rune(currentDisplay[p.view.GetCounter()-1])) {
				p.view.UpdatedisplayLabelWithText("error")
				return
			}

			operator = piValue
		}

		if operator == "e" {
			if currentDisplay == "0" || !helpers.IsValidInput(currentDisplay) {
				currentDisplay = ""
			}
			p.view.UpdatedisplayLabelWithText(eLiteral)
		}

		if !helpers.IsValidInput(currentDisplay) {
			currentDisplay = ""
		}
		p.view.UpdatedisplayLabelWithText(currentDisplay + operator)
	}
}

func (p *Presenter) AppendButtonText(inputText string) {
	if p.view.GetCounter() < 256 {
		currentDisplay := p.view.GetDisplayLabel()

		if p.view.GetCounter() >= 2 && currentDisplay[p.view.GetCounter()-2:p.view.GetCounter()] == "pi" && unicode.IsDigit(rune(inputText[0])) {
			p.view.UpdatedisplayLabelWithText("error")
			return
		}

		if (inputText != "-" && inputText != "(" && inputText != ")") && currentDisplay[p.view.GetCounter()-1] == 'x' {
			return
		}

		if currentDisplay == "0" || !helpers.IsValidInput(currentDisplay) {
			currentDisplay = ""
		}
		p.view.UpdatedisplayLabelWithText(currentDisplay + inputText)
	}
}

func (p *Presenter) ResetButton() {
	p.view.UpdatedisplayLabelWithText("0")
}

func (p *Presenter) AddDecimalPoint() {
	currentDisplay := p.view.GetDisplayLabel()

	if len(currentDisplay) == 0 {
		p.view.UpdatedisplayLabelWithText("0.")
		return
	}

	hasDecimal := false
	for i := len(currentDisplay) - 1; i >= 0; i-- {
		if currentDisplay[i] == '.' {
			hasDecimal = true
			break
		}
		if !unicode.IsDigit(rune(currentDisplay[i])) {
			break
		}
	}

	if !hasDecimal {
		p.view.UpdatedisplayLabelWithText(currentDisplay + ".")
	}
}

func (p *Presenter) DeleteButton() {
	currentDisplay := p.view.GetDisplayLabel()
	if len(currentDisplay) == 0 || currentDisplay == "0" {
		p.view.UpdatedisplayLabelWithText("0")
		return
	}

	if len(currentDisplay) >= len(piValue) && currentDisplay[len(currentDisplay)-len(piValue):] == piValue {
		currentDisplay = currentDisplay[:len(currentDisplay)-len(piValue)]
	} else {
		currentDisplay = currentDisplay[:len(currentDisplay)-1]
	}

	if len(currentDisplay) == 0 {
		currentDisplay = "0"
	}
	p.view.UpdatedisplayLabelWithText(currentDisplay)
}

func (p *Presenter) AppendX() {
	if p.view.GetCounter() < 256 {
		currentDisplay := p.view.GetDisplayLabel()
		if currentDisplay[p.view.GetCounter()-1] != 'x' && !unicode.IsDigit(rune(currentDisplay[p.view.GetCounter()-1])) {
			p.view.UpdatedisplayLabelWithText(currentDisplay + "x")
		}
		if currentDisplay == "0" {
			p.view.UpdatedisplayLabelWithText("x")
		}
	}
}

func (p *Presenter) InitializeXButton() {
	currentDisplay := p.view.GetDisplayLabel()
	if helpers.IsValidInput(currentDisplay) {
		p.EvaluateWithX(&currentDisplay, p.view.GetVariableXLabel())
	} else {
		p.view.UpdatedisplayLabelWithText("0")
	}
}

func (p *Presenter) InverseSign() {
	currentDisplay := p.view.GetDisplayLabel()
	if !helpers.IsValidInput(currentDisplay) {
		currentDisplay = "0"
	}

	if currentDisplay != "0" {
		currentDisplay = "(" + currentDisplay + ")*(-1)"
	}
	if currentDisplay != "0" {
		p.SaveHistory()
	}
	p.EvaluateExpression(&currentDisplay, p.view.GetVariableXLabel())
}

func (p *Presenter) SaveHistory() {
	historyFilePath := p.view.GetHistoryFilePath()
	if historyFilePath == "" {
		log.Println("History file path is not set. Skipping save.")
		return
	}

	file, err := os.OpenFile(historyFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open history file '%s': %v", historyFilePath, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", p.view.GetDisplayLabel()))
	if err != nil {
		log.Printf("Failed to write to history file '%s': %v", historyFilePath, err)
	}
}

func (p *Presenter) EvaluateAndProcessExpression() {
	currentDisplay := p.view.GetDisplayLabel()

	if strings.Contains(currentDisplay, ")") {
		for i := 0; i < len(currentDisplay)-1; i++ {
			if currentDisplay[i] == ')' && unicode.IsLetter(rune(currentDisplay[i+1])) {
				p.view.UpdatedisplayLabelWithText("error")
				return
			}
		}
	}

	if strings.Contains(currentDisplay, "pi") {
		index := strings.Index(currentDisplay, "pi")
		if (index > 0 && unicode.IsDigit(rune(currentDisplay[index-1]))) ||
			(index+2 < len(currentDisplay) && unicode.IsDigit(rune(currentDisplay[index+2]))) {
			p.view.UpdatedisplayLabelWithText("error")
			return
		}
	}

	currentDisplay = strings.ReplaceAll(currentDisplay, "pi", piValue)
	currentDisplay = helpers.ReplaceEConstant(currentDisplay)

	if !helpers.IsValidInput(currentDisplay) {
		currentDisplay = "0"
	}

	if currentDisplay != "0" {
		p.SaveHistory()
	}

	p.view.UpdatedisplayLabelWithText(p.calculateResult(&currentDisplay, p.view.GetVariableXLabel()))
}
