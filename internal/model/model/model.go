package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L.
#include "model_wrapper.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

func Calculate(expression *string, x string) (float64, error) {
	err := parser(expression, x)
	if err != nil {
		return 0.0, err
	}

	// Преобразуем строку Go в C-строку
	cStr := C.CString(*expression)
	defer C.free(unsafe.Pointer(cStr))

	// Создаём буфер для результата
	bufferSize := C.size_t(256)
	resultBuffer := (*C.char)(C.malloc(bufferSize))
	defer C.free(unsafe.Pointer(resultBuffer))

	// Вызываем C-функцию для вычисления
	status := C.calculate(cStr, resultBuffer, bufferSize)
	if status == 1 {
		// Преобразуем результат из C-строки в Go-строку
		result := C.GoString(resultBuffer)

		// Преобразуем строку в число с плавающей точкой
		value, err := strconv.ParseFloat(result, 64)
		if err != nil {
			return 0.0, err
		}

		return value, nil
	}

	return 0.0, errors.New("calculation error")
}

func parser(expression *string, x string) error {
	replacements := map[string]string{
		"sqrt": "q",
		"acos": "C",
		"asin": "S",
		"atan": "T",
		"cos":  "c",
		"sin":  "s",
		"tan":  "t",
		"log":  "L",
		"ln":   "l",
	}

	for old, new := range replacements {
		*expression = strings.ReplaceAll(*expression, old, new)
	}

	*expression = strings.ReplaceAll(*expression, " ", "")
	*expression = strings.ReplaceAll(*expression, "x", x)

	*expression = strings.ReplaceAll(*expression, "e-", "/10^")
	*expression = strings.ReplaceAll(*expression, "e+", "*10^")

	invalidOperators := regexp.MustCompile(`[+\-*/^%]{2,}`)
	if invalidOperators.MatchString(*expression) {
		return errors.New("invalid expression: contains consecutive operators")
	}

	invalidExpression := regexp.MustCompile(`\d+[a-zA-Z]+|\*\*|[^\d\+\-\*/^%.()a-zA-Z0-9]`)
	if invalidExpression.MatchString(*expression) {
		return errors.New("invalid expression: contains invalid characters or operators")
	}

	match, _ := regexp.MatchString(`^[+\-*/%^cstCSTqLle.()0123456789]+$`, *expression)
	if match {
		return nil
	}

	return errors.New("invalid expression")
}

func CreditAnnuity(creditSum, creditDuration, annualInterestRate float64) (float64, float64, float64, error) {
	if creditSum <= 0 || creditDuration <= 0 || annualInterestRate <= 0 {
		return 0.0, 0.0, 0.0, errors.New("invalid input values for CreditAnnuity")
	}

	var monthPay, overPay, totalPay C.double

	C.creditAnnuity(C.double(creditSum), C.double(creditDuration), C.double(annualInterestRate),
		&monthPay, &overPay, &totalPay)

	if monthPay == 0 || overPay == 0 || totalPay == 0 {
		return 0.0, 0.0, 0.0, errors.New("calculation error in CreditAnnuity: invalid result")
	}

	return float64(monthPay), float64(overPay), float64(totalPay), nil
}

func CreditDifferentiated(creditSum, creditDuration, annualInterestRate float64) (float64, float64, float64, float64, error) {
	if creditSum <= 0 || creditDuration <= 0 || annualInterestRate <= 0 {
		return 0.0, 0.0, 0.0, 0.0, errors.New("invalid input values for CreditDifferentiated")
	}

	var monthPayFirst, monthPayLast, overPay, totalPay C.double

	C.creditDifferentiated(C.double(creditSum), C.double(creditDuration), C.double(annualInterestRate),
		&monthPayFirst, &monthPayLast, &overPay, &totalPay)

	if monthPayFirst == 0 || monthPayLast == 0 || overPay == 0 || totalPay == 0 {
		return 0.0, 0.0, 0.0, 0.0, errors.New("calculation error in CreditDifferentiated: invalid result")
	}

	return float64(monthPayFirst), float64(monthPayLast), float64(overPay), float64(totalPay), nil
}

func main() {}
