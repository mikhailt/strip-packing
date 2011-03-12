#include "strip_packing.h"

void SimplePyramidAlgo::Pack(int n, double xbe, double ybe, Context* context) {
  InitParams(n);

  Bin bins[2][nsteps_];
  for (int i = 0; i < 2; ++i) {
    double h_left = h_;
    for (int y = 1; y < nsteps_; ++y) {
      Rect r(0, YOfType(y), WOfType(y), step_);
      ConvertCooToComplPyramid(&r, i);
      bins[i][y] = Bin(r, 0);
      if (context->opt.render_bins) {
        saved_rects.push_back(SavedRect(r, "red", false));
      }
    }
  }

  Rect r;
  for (int y = 0; y < n; ++y) {
    NextRect(&r);
    
    int32_t i;
    random_r(&__random_data, &i);
    i %= 2;
    PackToPyramid(bins[i], &r);
    
    if (context->opt.save_rects) {
      SaveRect(&r);
    }
    RecalcSolutionHeightSingle(&r);
  }
}

void SimplePyramidAlgo::InitParams(int n) {
  h_ = n / 4;
  nsteps_ = h_ / sqrt(n);
  step_ = h_ / nsteps_;
  top_h_ = 0;
  w_step_ = double(1) / nsteps_;
}

void SimplePyramidAlgo::PackToPyramid(Bin* bins, Rect* r) {
  int t = TypeOfW(r->w);
  if (0 == t) {
    PackOnTop(r);
    return;
  }
  for ( ; t; --t) {
    if (bins[t].VacantSpace() >= r->h) {
      PackToBin(&bins[t], r);
      return;
    }
  }
  PackOnTop(r);
}

int SimplePyramidAlgo::TypeOfW(double w) {
  return (1 - w) / w_step_;
}

double SimplePyramidAlgo::YOfType(int t) {
  return (t - 1) * step_;
}

double SimplePyramidAlgo::WOfType(int t) {
  return 1 - t * w_step_;
} 

void SimplePyramidAlgo::ConvertCooToComplPyramid(Rect* r, int ind) {
  if (0 == ind) {
    return;
  }
  r->x = 1 - r->w;
  r->y = h_ - step_ - r->y - r->h;
}

void SimplePyramidAlgo::PackOnTop(Rect* r) {
  r->x = 0;
  r->y = top_h_ + h_ - step_;
  top_h_ += r->h;
}