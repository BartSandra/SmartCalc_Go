#include "model_wrapper.h"

#include <iostream>
#include <string>

#include "model.h"

int calculate(const char *expression, char *resultBuffer, size_t bufferSize) {
  s21::Model model;
  std::string expr(expression);

  if (model.Calculate(expr, 0.0)) {
    try {
      std::string formatted = expr;
      strncpy(resultBuffer, formatted.c_str(), bufferSize - 1);
      resultBuffer[bufferSize - 1] = '\0';
    } catch (const std::exception &e) {
      std::cerr << "Conversion error: " << e.what() << std::endl;
      return 0;
    }
    return 1;
  }
  return 0;
}

void creditAnnuity(double sum_of_credit, double duration_of_credit,
                   double annual_interest_rate, double *month_pay,
                   double *over_pay, double *all_sum_of_pay) {
  if (sum_of_credit <= 0 || duration_of_credit <= 0 ||
      annual_interest_rate <= 0) {
    std::cerr << "Invalid input values for creditAnnuity." << std::endl;
    return;
  }
  s21::Model model;
  try {
    model.CreditAnnuity(sum_of_credit, duration_of_credit, annual_interest_rate,
                        month_pay, over_pay, all_sum_of_pay);
  } catch (const std::exception &e) {
    std::cerr << "Credit Annuity Calculation Error: " << e.what() << std::endl;
  }
}

void creditDifferentiated(double sum_of_credit, double duration_of_credit,
                          double annual_interest_rate, double *month_pay_first,
                          double *month_pay_last, double *over_pay,
                          double *all_sum_of_pay) {
  if (sum_of_credit <= 0 || duration_of_credit <= 0 ||
      annual_interest_rate <= 0) {
    std::cerr << "Invalid input values for creditDifferentiated." << std::endl;
    return;
  }
  s21::Model model;
  try {
    model.CreditDifferentiated(sum_of_credit, duration_of_credit,
                               annual_interest_rate, month_pay_first,
                               month_pay_last, over_pay, all_sum_of_pay);
  } catch (const std::exception &e) {
    std::cerr << "Credit Differentiated Calculation Error: " << e.what()
              << std::endl;
  }
}
