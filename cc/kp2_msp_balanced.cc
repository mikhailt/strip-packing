#include "strip_packing.h"

void Kp2MspBalanced::PackWithSize(int n, double xbe, double ybe, Context* context) {
  InitParams(n);
  Rect r;
  for (int i = 0; i < n; ++i) {
    std::set<Bin>::iterator j = frames_.begin();
    Bin best_frame = *j;
    frames_.erase(j);
    
    NextRect(&r);
    if (r.w > (1 - delta_)) {
      PackToBin(&best_frame, &r);
    } else {
      int j = RectType(&r);
      if (!PackToTopBin(&r, j)) {
        PackToNewShelfInFrame(&r, &best_frame, j);
      }
    }
    if (context->opt.save_rects) {
      saved_rects.push_back(SavedRect(r, "black", true));
    }
    
    solution_height = std::max<double>(solution_height, r.y + r.h);
    frames_.insert(best_frame);
  }
  SaveBins(context);
}

void Kp2MspBalanced::Pack(int n, double xbe, double ybe, Context* context) {
  solution_height = ybe;
  InitFrames(xbe, ybe, context->opt.m);
  int s = 0;
  for (int cnt = 2; s < context->opt.n; cnt *= 2) {
    int cur_n = std::min<int>(context->opt.n - s, cnt);
    s += cur_n;
    PackWithSize(cur_n, xbe, ybe, context);
  }
}

bool operator< (const Bin& a, const Bin& b);

void Kp2MspBalanced::InitFrames(double xbe, double ybe, int m) {
  static const double MAX_Y = 10000000000.0; // 10 billions.
  for (int y = 0; y < m; ++y) {
    frames_.insert(Bin(1, MAX_Y, 0, xbe + y, ybe));
  }
}