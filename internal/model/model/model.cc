#include "model.h"

#include <iostream>

std::string DoubleToString(double value) {
  std::ostringstream out;
  out << std::fixed << std::setprecision(16) << value;
  return out.str();
}

namespace s21 {

int s21::Model::Calculate(std::string &str, double x) {
  int flag = 0;
  if (str.length() < 255) {
    Replace(str, "e-", "/10^");
    Replace(str, "e+", "*10^");
    ReplaceDot(str, ".", "0.", ".0");
  } else {
    str = "ERROR";
  }
  if (CheckFunctions(str)) {
    try {
      auto tmp = Parser(str);
      ValueX(x, tmp);
      PolisNotation(tmp);
      double result_polis_notation = Counter(tmp);
      str = DoubleToString(result_polis_notation);
      Check(str);
      CheckE(str);
      flag = 1;
    } catch (const std::runtime_error &e) {
      str = e.what();
    }
  } else {
    str = "ERROR";
  }
  return flag;
}

std::list<s21::Model::Lexeme> s21::Model::Parser(const std::string &str) {
  std::list<Lexeme> List;
  char first = str[0];
  for (size_t i = 0; i < str.length(); i++) {
    if (isdigit(str[i]) || str[i] == 'e') {
      if (str[i] == 'e') {
        List.push_back(Lexeme(EPS, 20, 0));
      } else {
        double num = atof(&str[i]);
        List.push_back(Lexeme(num, 20, 0));
        while (isdigit(str[i]) || str[i] == '.') {
          i++;
        }
        i--;
      }
    } else if (str[i] == '+' || str[i] == '-') {
      if (first == '+') {
        List.push_back(Lexeme(0, 16, 3));
        first = ';';
      } else if (first == '-') {
        List.push_back(Lexeme(0, 17, 3));
        first = ';';
      } else if ((!isdigit(str[i - 1]) && str[i - 1] != 'x' &&
                  str[i - 1] != ')') ||
                 (str[i - 1] == '(')) {
        if (str[i] == '+') List.push_back(Lexeme(0, 16, 3));
        if (str[i] == '-') List.push_back(Lexeme(0, 17, 3));
      } else {
        int operation = GetOperators(&str[i]);
        List.push_back(Lexeme(0, operation, GetPriorities(operation)));
      }
    } else if (str[i] == '*' || str[i] == '/' || str[i] == '%' ||
               str[i] == '^' || str[i] == ')' || str[i] == '(') {
      int operation = GetOperators(&str[i]);
      List.push_back(Lexeme(0, operation, GetPriorities(operation)));
    } else if (str[i] == 's' || str[i] == 'q' || str[i] == 'C' ||
               str[i] == 'S' || str[i] == 'T' || str[i] == 'c' ||
               str[i] == 't' || str[i] == 'l' || str[i] == 'L') {
      int operation = GetOperators(&str[i]);
      List.push_back(Lexeme(0, operation, GetPriorities(operation)));
      while (!isdigit(str[i]) && str[i] != '(') {
        i++;
      }
      i = i - 1;
    } else if (str[i] == 'x') {
      List.push_back(Lexeme(0, 21, 0));
    }
  }
  return List;
}

void s21::Model::PolisNotation(std::list<s21::Model::Lexeme> &List) {
  std::list<Lexeme> list_lexeme, support;
  std::list<Lexeme> tmp;
  tmp = List;
  for (auto i = tmp.begin(); i != tmp.end(); i++) {
    if ((*i).GetOperation() == 20) {
      list_lexeme.push_back(Lexeme((*i).GetValue(), 20, 0));
      List.pop_front();
    } else if ((*i).GetOperation() == 19) {
      while (support.back().GetOperation() != 18) {
        list_lexeme.push_back(Lexeme(0, support.back().GetOperation(),
                                     support.back().GetPriority()));
        support.pop_back();
      }
      support.pop_back();
    } else if ((*i).GetOperation() < 19) {
      if (support.empty()) {
        support.push_back(Lexeme(0, (*i).GetOperation(), (*i).GetPriority()));
        List.pop_front();
      } else {
        if (CheckPriority(support, (*i).GetPriority()) &&
            support.back().GetOperation() != 18) {
          while (support.size()) {
            if (CheckPriority(support, (*i).GetPriority()) &&
                support.back().GetOperation() != 18) {
              list_lexeme.push_back(Lexeme(0, support.back().GetOperation(),
                                           support.back().GetPriority()));
              support.pop_back();
            } else {
              break;
            }
          }
          support.push_back(Lexeme(0, (*i).GetOperation(), (*i).GetPriority()));
          List.pop_front();
        } else {
          support.push_back(Lexeme(0, (*i).GetOperation(), (*i).GetPriority()));
          List.pop_front();
        }
      }
    }
  }
  while (!support.empty()) {
    list_lexeme.push_back(
        Lexeme(0, support.back().GetOperation(), support.back().GetPriority()));
    support.pop_back();
  }
  List = list_lexeme;
}

int s21::Model::GetOperators(const char *str) {
  int res = 0;
  const char *tmp = str;
  switch (*str) {
    case '+':
      res = 1;
      if (*tmp--) {
        if (*tmp == '(') res = 16;
      }
      break;
    case '-':
      res = 2;
      if (*tmp--) {
        if (*tmp == '(') res = 17;
      }
      break;
    case '*':
      res = 3;
      break;
    case '/':
      res = 4;
      break;
    case '%':
      res = 5;
      break;
    case '^':
      res = 6;
      break;
    case 's':
      res = 11;
      break;
    case 'q':
      res = 7;
      break;
    case 'C':
      res = 8;
      break;
    case 'S':
      res = 9;
      break;
    case 'T':
      res = 10;
      break;
    case 'c':
      res = 12;
      break;
    case 't':
      res = 13;
      break;
    case 'l':
      res = 14;
      break;
    case 'L':
      res = 15;
      break;
    case '(':
      res = 18;
      break;
    case ')':
      res = 19;
      break;
  }
  return res;
}

int s21::Model::GetPriorities(int oper) {
  int res = 0;
  if (oper == 1 || oper == 2 || oper == 5) {
    res = 1;
  }
  if (oper > 2 && oper < 6) {
    res = 2;
  }
  if (oper == 6) {
    res = 3;
  }
  if (oper > 6 && oper < 16) {
    res = 4;
  }
  if (oper == 18 || oper == 19) {
    res = 5;
  }
  return res;
}

void s21::Model::Replace(std::string &str, const std::string a,
                         const std::string b) {
  size_t i = 0;
  while ((i = str.find(a)) != std::string::npos) {
    str.replace(i, a.size(), b);
  }
}

void s21::Model::ReplaceDot(std::string &str, const std::string a,
                            const std::string b, const std::string c) {
  int next_ = 0;
  int prev_ = 0;
  for (size_t i = 0; i < str.size(); i++) {
    if (str[i] == '.') {
      if (str[i + 1]) {
        next_ = isdigit(str[i + 1]);
        if (!next_) {
          str.replace(i, a.size(), c);
        }
        if (next_ && str[i - 1]) {
          prev_ = isdigit(str[i - 1]);
          if (!prev_) {
            str.replace(i, a.size(), b);
          }
        }
        if (i == 0) {
          str = "0" + str;
        }
      }
    }
  }
}

void s21::Model::ValueX(double value, std::list<s21::Model::Lexeme> &List) {
  std::list<Lexeme> tmp;
  for (auto i = List.begin(); i != List.end(); i++) {
    if ((*i).GetOperation() == 21) {
      tmp.push_back(Lexeme(value, 20, 0));
    } else {
      tmp.push_back(
          Lexeme((*i).GetValue(), (*i).GetOperation(), (*i).GetPriority()));
    }
  }
  List = tmp;
}

double s21::Model::Counter(std::list<s21::Model::Lexeme> &List) {
  std::list<double> result;
  for (auto &i_lexeme : List) {
    if (i_lexeme.GetOperation() == 20) {
      result.push_back(i_lexeme.GetValue());
    } else if (i_lexeme.GetOperation() == 16 || i_lexeme.GetOperation() == 17) {
      UnaryOperations(i_lexeme.GetOperation(), result);
    } else if (i_lexeme.GetOperation() < 7) {
      Operations(i_lexeme.GetOperation(), result);
    } else {
      Functions(i_lexeme.GetOperation(), result);
    }
  }
  return result.back();
}

void s21::Model::UnaryOperations(int oper, std::list<double> &res) {
  double result = 0;
  if (!res.empty()) {
    if (oper == 16) result = +res.back();
    if (oper == 17) result = -res.back();
    res.pop_back();
    res.push_back(result);
  }
}

void s21::Model::Operations(int oper, std::list<double> &res) {
  double result = 0.0;
  double second = res.back();
  res.pop_back();
  if (oper == 1) result = res.back() + second;
  if (oper == 2) result = res.back() - second;
  if (oper == 3) result = res.back() * second;
  if (oper == 4) {
    if (second == 0) {
      throw std::runtime_error("ERROR");
    }
    result = res.back() / second;
  }
  if (oper == 5) result = fmod(res.back(), second);
  if (oper == 6) result = pow(res.back(), second);
  res.pop_back();
  res.push_back(result);
}

void s21::Model::Functions(int oper, std::list<double> &res) {
  double result = 0;
  if (oper == 7) result = sqrt(res.back());
  if (oper == 8) result = acos(res.back());
  if (oper == 9) result = asin(res.back());
  if (oper == 10) result = atan(res.back());
  if (oper == 11) result = sin(res.back());
  if (oper == 12) result = cos(res.back());
  if (oper == 13) result = tan(res.back());
  if (oper == 14) result = log(res.back());
  if (oper == 15) result = log10(res.back());
  res.pop_back();
  res.push_back(result);
}

void s21::Model::Check(std::string &str) {
  if (str.find('e') == std::string::npos) {
    size_t end = str.find_last_not_of('0');
    if (end != std::string::npos && str[end] == '.') {
      end--;
    }
    str.erase(end + 1);
  }
}

void s21::Model::CheckE(std::string &str) {
  size_t pos = str.find('e');
  if (pos != std::string::npos) {
    std::string exponent = str.substr(pos);
    str = str.substr(0, pos);
    str = DoubleToString(std::stod(str));
    str += exponent;
  }
}

bool s21::Model::CheckPriority(std::list<s21::Model::Lexeme> &list_lexeme,
                               int i_priority) {
  bool flag = 0;
  if (!list_lexeme.empty()) {
    if (i_priority <= list_lexeme.back().GetPriority()) {
      flag = 1;
    }
    if (list_lexeme.back().GetPriority() == 3 && i_priority == 3) {
      flag = 0;
    }
  }
  return flag;
}

int s21::Model::CheckFunctions(const std::string &str) {
  bool res = 1;
  if (str.size() == 1) {
    if (str[0] == '.') return res = 0;
  }
  if (CheckBracket(&str[0]) && CheckSymbol(&str[str.length() - 1])) {
    for (size_t i = 0; i < (str.length() - 1); ++i) {
      if (isdigit(str[i])) {
        res = CheckDot(&str[i + 1]);
        if (res == 0) return res;
        while (isdigit(str[i]) || str[i] == '.') {
          i++;
        }
        if (i != str.length()) i--;
        if (i != str.length()) {
          if (str[i + 1] == '+' || str[i + 1] == '-' || str[i + 1] == '*' ||
              str[i + 1] == '/' || str[i + 1] == '%' || str[i + 1] == '^' ||
              str[i + 1] == ')')
            res = 1;
        }
      } else if (str[i] == '^' || str[i] == '%') {
        res = CheckMod(&str[i]);
      } else if (str[i] == '*' || str[i] == '/' || str[i] == '^' ||
                 str[i] == '(') {
        res = CheckSign(str[i + 1]);
      } else if (str[i] == ')') {
        if (str[i + 1] == '+' || str[i + 1] == '-' || str[i + 1] == '*' ||
            str[i + 1] == '/' || str[i + 1] == '%' || str[i + 1] == '^' ||
            str[i + 1] == ')')
          res = 1;
      } else if (str[i] == '-' || str[i] == '+') {
        res = CheckSign(str[i + 1]);
      } else if (str[i] == 'a') {
        if (!std::strncmp(&str[i], "acos(", 5)) res = 1;
        if (!std::strncmp(&str[i], "asin(", 5)) res = 1;
        if (!std::strncmp(&str[i], "atan(", 5)) res = 1;
        i += 3;
      } else if (str[i] == 'c') {
        if (!std::strncmp(&str[i], "cos(", 4)) res = 1;
        i += 2;
      } else if (str[i] == 't') {
        if (!std::strncmp(&str[i], "tan(", 4)) res = 1;
        i += 2;
      } else if (str[i] == 's' && str[i + 1] == 'i') {
        if (!std::strncmp(&str[i], "sin(", 4)) res = 1;
        i += 2;
      } else if (str[i] == 's' && str[i + 1] == 'q') {
        res = CheckSqrt(&str[i]);
        i += 3;
      } else if (str[i] == 'l' && str[i + 1] == 'o') {
        if (!std::strncmp(&str[i], "log(", 4)) res = 1;
        i += 2;
      } else if (str[i] == 'l' && str[i + 1] == 'n') {
        if (!std::strncmp(&str[i], "ln(", 3)) res = 1;
        i += 1;
      }
      if (res == 0) return 0;
    }
  } else {
    res = 0;
  }
  return res;
}

bool s21::Model::CheckDot(const char *str) {
  bool res = 0;
  int dot = 0;
  while (isdigit(*str) || *str == '.') {
    if (*str == '.') {
      dot++;
    }
    str += 1;
  }
  if (dot == 0 || dot == 1) {
    res = 1;
  }
  return res;
}

bool s21::Model::CheckBracket(const char *str) {
  bool flag = 1;
  int open = 0, cloze = 0;
  while (*str) {
    if (*str == '(') open++;
    if (*str == ')') cloze++;
    if (cloze > open) {
      flag = 0;
      break;
    }
    str += 1;
  }
  if (cloze != open) {
    flag = 0;
  }
  return flag;
}

bool s21::Model::CheckSymbol(const char *s) {
  bool flag = 0;
  if (*s == '.' || *s == ')' || *s == 'x' || isdigit(*s)) {
    flag = 1;
  }
  return flag;
}

bool s21::Model::CheckSign(const char s) {
  bool flag = 0;
  if (isdigit(s) || s == '(' || s == 'c' || s == 'C' || s == 'T' || s == 'S' ||
      s == 'L' || s == 'q' || s == 's' || s == 't' || s == 'a' || s == 'l' ||
      s == '+' || s == '-' || s == 'x')
    flag = 1;
  return flag;
}

bool s21::Model::CheckMod(const char *str) {
  bool flag = 0;
  const char *tmp = str;
  if (*tmp++) {
    if (*str--) {
      if (CheckSign(*tmp) && (isdigit(*str) || *str == ')' || *str == 'x'))
        flag = 1;
    }
  }
  return flag;
}

bool s21::Model::CheckSqrt(const char *str) {
  bool flag = 0;
  const char *tmp = str;
  tmp += 5;
  if ((!std::strncmp(str, "sqrt(", 5)) && (CheckSign(*tmp))) flag = 1;
  return flag;
}

}  // namespace s21
