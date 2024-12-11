#ifndef MODEL_WRAPPER_H
#define MODEL_WRAPPER_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

int calculate(const char *expression, char *resultBuffer, size_t bufferSize);

void creditAnnuity(double sum_of_credit, double duration_of_credit,
                   double annual_interest_rate, double *month_pay,
                   double *over_pay, double *all_sum_of_pay);

void creditDifferentiated(double sum_of_credit, double duration_of_credit,
                          double annual_interest_rate, double *month_pay_first,
                          double *month_pay_last, double *over_pay,
                          double *all_sum_of_pay);

#ifdef __cplusplus
}
#endif

#endif  // MODEL_WRAPPER_H
