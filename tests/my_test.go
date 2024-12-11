package test

import (
	"fmt"
	"log"
	"math"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"sandra/APG2_SmartCalc_v3.0_Desktop_Go-1/src/internal/model"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	code := m.Run()
	os.Exit(code)
}

func getModelPath() string {
	modelPath := os.Getenv("MODEL_PATH")
	if modelPath == "" {
		modelPath = "./internal/model/model/model.so"
	}
	return modelPath
}

func TestSimpleNumericExpressions(t *testing.T) {
	expressions := []string{
		"0/7",
		"42421135678 - 0.3423433132323",
		"2+3^3*2",
		"-(2+2)",
		"(1-2)+4",
		"2*(-2-2)",
		"-(9-2)",
		"-1-3",
		"1-3",
		"1000+333-123.23",
		"5*5+10/2",
		"3*3^3-10",
		"(2+2)*(3-1)",
		"sqrt(144)",
		"1.2345*1000-1234",
		"1000/2*4",
		"-4+6/2-3",
		"1e+10",
		"1e-10",
	}
	expectedResults := []float64{
		0 / 7,
		42421135678 - 0.3423433132323,
		2 + math.Pow(3, 3)*2,
		-(2 + 2),
		(1 - 2) + 4,
		2 * (-2 - 2),
		-(9 - 2),
		-1 - 3,
		1 - 3,
		1000 + 333 - 123.23,
		5*5 + 10/2,
		3*math.Pow(3, 3) - 10,
		(2 + 2) * (3 - 1),
		math.Sqrt(144),
		1.2345*1000 - 1234,
		1000 / 2 * 4,
		-4 + 6/2 - 3,
		1e+10,
		1e-10,
	}

	calc, err := model.NewModel(getModelPath())
	if err != nil {
		t.Fatalf("Error creating the model: %v", err)
	}

	for i, expr := range expressions {
		t.Run(fmt.Sprintf("SimpleExpr%d", i), func(t *testing.T) {
			got, err := calc.Calculate(&expr, "")
			if err != nil {
				t.Errorf("Error calculating expression: %s. Error: %v", expr, err)
			} else if math.Abs(got-expectedResults[i]) > 1e-9 {
				t.Errorf("Calc(%s) = %.8f, expected %.8f", expr, got, expectedResults[i])
			}
		})
	}
}

func TestExpressionsWithFunctions(t *testing.T) {
	expressions := []string{
		"-(7+(4+c(1)))",
		"-(c(c(2)))",
		"s(1)+c(1)-3*(-2)*(30%2)+s(1)",
		"1+c(2)+(5-10)",
		"c(4)+s(5)",
		"s(1)+c(1)-3+2*(-2)+3^3-l(10)",
		"t(7)",
		"S(1)",
		"T(2)",
		"L(2)",
		"q(q(79))",
		"t(S(0.5))",
		"l(10)",
		"s(4)",
		"t(45)",
		"S(0.5)+C(0.5)",
	}
	expectedResults := []float64{
		-(7 + (4 + math.Cos(1))),
		-(math.Cos(math.Cos(2))),
		math.Sin(1) + math.Cos(1) - 3*-2*math.Mod(30, 2) + math.Sin(1),
		1 + math.Cos(2) + (5 - 10),
		math.Cos(4) + math.Sin(5),
		math.Sin(1) + math.Cos(1) - 3 + 2*-2 + math.Pow(3, 3) - math.Log(10),
		math.Tan(7),
		math.Asin(1),
		math.Atan(2),
		math.Log10(2),
		math.Sqrt(math.Sqrt(79)),
		math.Tan(math.Asin(0.5)),
		math.Log(10),
		math.Sin(4),
		math.Tan(45),
		math.Asin(0.5) + math.Acos(0.5),
	}

	calc, err := model.NewModel(getModelPath())
	if err != nil {
		t.Fatalf("Error creating the model: %v", err)
	}

	for i, expr := range expressions {
		t.Run(fmt.Sprintf("FunctionExpr%d", i), func(t *testing.T) {
			got, err := calc.Calculate(&expr, "")
			if err != nil {
				t.Errorf("Error calculating function: %s. Error: %v", expr, err)
			} else if math.Abs(got-expectedResults[i]) > 1e-9 {
				t.Errorf("Calc(%s) = %.8f, expected %.8f", expr, got, expectedResults[i])
			}
		})
	}
}

func TestComplexExpressions(t *testing.T) {
	expressions := []string{
		"2+3*(4-1)^2/3",
		"q(16) + (5^2 - 3*3)/3",
		"l(20) * s(1) + c(0.5)^2",
		"(2+2)*L(100)/q(25)",
		"-4 + t(45) - 5^3 + l(3)",
		"c(C(0.5)) + s(S(0.5))",
		"T(1)*4",
		"l(7) * L(1000)",
		"4 * (3 + 5 * (2 - 7)^2) / 2 - 10",
		"(1 + 2*3^2) * (sqrt(49) - ln(1))",
		"3^3^2 / (2^5)",
		"1e+10 + 1e-10 - (2^30)",
		"q(100) + c(0) + t(0)",
		"2^(1/3)",
	}
	expectedResults := []float64{
		2 + 3*math.Pow(4-1, 2)/3,
		math.Sqrt(16) + (math.Pow(5, 2)-3*3)/3,
		math.Log(20)*math.Sin(1) + math.Pow(math.Cos(0.5), 2),
		(2 + 2) * math.Log10(100) / math.Sqrt(25),
		-4 + math.Tan(45) - math.Pow(5, 3) + math.Log(3),
		math.Cos(math.Acos(0.5)) + math.Sin(math.Asin(0.5)),
		math.Atan(1) * 4,
		math.Log(7) * math.Log10(1000),
		4*(3+5*math.Pow(2-7, 2))/2 - 10,
		(1 + 2*math.Pow(3, 2)) * (math.Sqrt(49) - math.Log(1)),
		math.Pow(3, math.Pow(3, 2)) / math.Pow(2, 5),
		1e+10 + 1e-10 - math.Pow(2, 30),
		math.Sqrt(100) + math.Cos(0) + math.Tan(0),
		math.Pow(2, 1.0/3.0),
	}

	calc, err := model.NewModel(getModelPath())
	if err != nil {
		t.Fatalf("Error creating the model: %v", err)
	}

	for i, expr := range expressions {
		t.Run(fmt.Sprintf("ComplexExpr%d", i), func(t *testing.T) {
			got, err := calc.Calculate(&expr, "")
			if err != nil {
				t.Errorf("Error calculating complex expression: %s. Error: %v", expr, err)
			} else if math.Abs(got-expectedResults[i]) > 1e-9 {
				t.Errorf("Calc(%s) = %.8f, expected %.8f", expr, got, expectedResults[i])
			}
		})
	}
}

func TestCalcError(t *testing.T) {
	invalidExpressions := []string{
		"5c",
		"(((1+2))",
		"7/0",
		"^Cc",
		"1e*10",
		"^99",
		"--4-5",
		"8++2",
		"-(-(-(-(-(5+9*)))))",
		"",
		"5++2",
		"1/0",
		"^++2",
		"2(3^2",
		"q(-16)",
	}

	calc, err := model.NewModel(getModelPath())
	if err != nil {
		t.Fatalf("Error creating the model: %v", err)
	}

	for _, expr := range invalidExpressions {
		t.Run(fmt.Sprintf("InvalidExpr_%s", expr), func(t *testing.T) {
			got, err := calc.Calculate(&expr, "")
			if err == nil {
				t.Errorf("An error was expected for the expression: %s, but none occurred", expr)
			} else if got != 0.0 {
				t.Errorf("Result for invalid expression %s should be 0.0, got: %.8f", expr, got)
			}
		})
	}
}
