#include "model.h"

namespace s21 {
void Model::CreditAnnuity(double sum_of_credit, double duration_of_credit,
                          double annual_interest_rate, double *month_pay,
                          double *over_pay, double *all_sum_of_pay) {
  double monthly_interest_rate = (annual_interest_rate / 12.0) / 100;
  *month_pay = monthly_interest_rate *
               pow(1 + monthly_interest_rate, duration_of_credit) /
               (pow(1 + monthly_interest_rate, duration_of_credit) - 1) *
               sum_of_credit;
  *over_pay = *month_pay * duration_of_credit - sum_of_credit;
  *all_sum_of_pay = *over_pay + sum_of_credit;
}

void Model::CreditDifferentiated(double sum_of_credit,
                                 double duration_of_credit,
                                 double annual_interest_rate,
                                 double *month_pay_first,
                                 double *month_pay_last, double *over_pay,
                                 double *all_sum_of_pay) {
  annual_interest_rate = annual_interest_rate / 100;
  *all_sum_of_pay = 0.0;
  double sum_in_month = sum_of_credit / duration_of_credit;
  double sum = sum_of_credit;
  *month_pay_last = 0.0;

  for (int i = 0; i < duration_of_credit; i++) {
    *month_pay_last =
        sum * annual_interest_rate * 30.4166666666667 / 365 + sum_in_month;
    sum = sum - sum_in_month;
    *all_sum_of_pay = *all_sum_of_pay + *month_pay_last;
    if (!i) {
      *month_pay_first = *month_pay_last;
    }
  }
  *over_pay = *all_sum_of_pay - sum_of_credit;
}

}  // namespace s21
