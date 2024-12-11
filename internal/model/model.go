package model

import (
	"fmt"
	"log"
	"plugin"
)

type Model struct {
	library              string
	creditAnnuity        func(float64, float64, float64) (float64, float64, float64, error)
	creditDifferentiated func(float64, float64, float64) (float64, float64, float64, float64, error)
	calculate            func(*string, string) (float64, error)
	pluginObj            *plugin.Plugin
}

func NewModel(libraryPath string) (*Model, error) {
	if libraryPath == "" {
		return nil, fmt.Errorf("library path cannot be empty")
	}

	plug, err := plugin.Open(libraryPath)
	if err != nil {
		log.Printf("Error opening plugin: %v", err)
		return nil, err
	}

	symCreditAnnuity, err := plug.Lookup("CreditAnnuity")
	if err != nil {
		log.Printf("Error finding creditAnnuity function: %v", err)
		return nil, err
	}

	//log.Printf("Type of symCreditAnnuity: %T", symCreditAnnuity)

	creditAnnuityFunc, ok := symCreditAnnuity.(func(float64, float64, float64) (float64, float64, float64, error))
	if !ok {
		err := fmt.Errorf("unexpected type for creditAnnuity: %T", symCreditAnnuity)
		log.Println(err)
		return nil, err
	}

	symCreditDifferentiated, err := plug.Lookup("CreditDifferentiated")
	if err != nil {
		log.Printf("Error finding creditDifferentiated function: %v", err)
		return nil, err
	}

	//log.Printf("Type of symCreditDifferentiated: %T", symCreditDifferentiated)

	creditDifferentiatedFunc, ok := symCreditDifferentiated.(func(float64, float64, float64) (float64, float64, float64, float64, error))
	if !ok {
		err := fmt.Errorf("unexpected type for creditDifferentiated: %T", symCreditDifferentiated)
		log.Println(err)
		return nil, err
	}

	symCalculate, err := plug.Lookup("Calculate")
	if err != nil {
		log.Printf("Error finding Calculate function: %v", err)
		return nil, err
	}

	calculateFunc, ok := symCalculate.(func(*string, string) (float64, error))
	if !ok {
		err := fmt.Errorf("failed to cast Calculate function")
		log.Println(err)
		return nil, err
	}

	return &Model{
		library:              libraryPath,
		creditAnnuity:        creditAnnuityFunc,
		creditDifferentiated: creditDifferentiatedFunc,
		calculate:            calculateFunc,
		pluginObj:            plug,
	}, nil
}

func (m *Model) Calculate(s *string, x string) (float64, error) {

	if m.pluginObj == nil || m.calculate == nil {
		err := fmt.Errorf("plugin not loaded or Calculate function not set")
		log.Println(err)
		return 0, err
	}
	return m.calculate(s, x)
}

func (m *Model) CreditAnnuity(sum_of_credit, duration_of_credit, annual_interest_rate float64) (float64, float64, float64, error) {

	if m.pluginObj == nil || m.creditAnnuity == nil {
		err := fmt.Errorf("plugin not loaded or creditAnnuity function not set")
		log.Println(err)
		return 0, 0, 0, err
	}

	month_pay, over_pay, all_sum_of_pay, err := m.creditAnnuity(sum_of_credit, duration_of_credit, annual_interest_rate)

	if err != nil {
		log.Printf("Error in CreditAnnuity: %v", err)
		return 0, 0, 0, err
	}

	return month_pay, over_pay, all_sum_of_pay, nil
}

func (m *Model) CreditDifferentiated(sum_of_credit, duration_of_credit, annual_interest_rate float64) (float64, float64, float64, float64, error) {

	if m.pluginObj == nil || m.creditDifferentiated == nil {
		err := fmt.Errorf("plugin not loaded or creditDifferentiated function not set")
		log.Println(err)
		return 0, 0, 0, 0, err
	}

	month_pay_first, month_pay_last, over_pay, all_sum_of_pay, err := m.creditDifferentiated(sum_of_credit, duration_of_credit, annual_interest_rate)

	if err != nil {
		log.Printf("Error in CreditDifferentiated: %v", err)
		return 0, 0, 0, 0, err
	}

	return month_pay_first, month_pay_last, over_pay, all_sum_of_pay, nil
}
