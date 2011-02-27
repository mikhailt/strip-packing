#include <iomanip>

#include "strip_packing.h"

int main(int argc, char** argv) {
  std::cout << std::fixed;
  std::cout << std::setprecision(4);
  Context context(argc, argv);
  double coef_s = 0;
  double divisor = pow(context.opt.n, 2.0 / 3.0);
  std::cout << "n = " << context.opt.n << "\n";
  std::cout << "n^(2/3) = " << divisor << "\n";
  for (int y = 0; y < context.opt.t; ++y) {
    context.InitAlgo(context.opt.algo);
    std::cout << "=============\niteration " << y << "\n=============\n";
    double h = context.algo->Pack(context.opt.n, 0, 0, &context);
    std::cout << "height of solution = " << h << "\n";
    std::cout << "total area = " << context.algo->total_area << "\n";
    double vacant_area = double(context.opt.m) * h - context.algo->total_area;
    std::cout << "uncovered area = " << vacant_area << "\n";
    double coef = vacant_area / divisor;
    std::cout << "coefficient = " << coef << "\n";
    coef_s += coef;
  }
  context.Render();
  std::cout << "Average coef = " << coef_s / double(context.opt.t) << "\n";
}