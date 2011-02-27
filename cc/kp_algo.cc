#include "strip_packing.h"

//void Algorithm::CountRect(Rect* rect) {
 // total_area += rect->w * rect->h;
//}

double Bin::VacantSpace() const {
  return h - top;
}

bool operator< (const Bin& a, const Bin& b) {
  return a.VacantSpace() > b.VacantSpace();
}

double Kp1Algo::Pack(int n, double xbe, double ybe, Context* context) {
  frame_.top = 0;
  frame_.y = ybe;
  frame_.x = xbe;
  frame_.w = 1;
  
  InitParams(n);
  Rect r;
  for (int i = 0; i < n; ++i) {
    NextRect(&r);
    if (r.w > (1 - delta_)) {
      PackToBin(&frame_, &r);
    } else {
      int j = RectType(&r);
      if (!PackToTopBin(&r, j)) {
        PackToNewShelfInFrame(&r, &frame_, j);
      }
    }
    if (context->opt.save_rects) {
      saved_rects.push_back(SavedRect(r, "black", true));
    }
  }
  SaveBins(context);
  solution_height = frame_.y + frame_.top;
  return solution_height;
}

void Kp1Algo::SaveBins(Context* context) {
  if (!context->opt.render_bins) {
    return;
  }
  SetOfBins::iterator i;
  for (int y = 0; y <= d_ * 2 + 1; ++y) {
    for (i = bins_[y].begin(); i != bins_[y].end(); ++i) {
      saved_rects.push_back(SavedRect(*i, "red", false));
    }
  }
}

void Kp1Algo::InitParams(int n) {
  saved_rects.clear();
  bins_.clear();
  delta_ = pow(double(n), -1.0 / 3.0);
  u_ = pow(double(n), 1.0 / 3.0);
  d_ = int(1.0 / (2 * delta_));
}

// Returns true/false whether rectangle was packed into top-of-the-heap bin.
// @j is the type of rectangle @r.
bool Kp1Algo::PackToTopBin(Rect* r, int j) {
  if (0 == bins_[j].size()) {
    return false;
  }
  MapOfSets::iterator mi = bins_.find(j);
  SetOfBins::iterator i = mi->second.begin();
  Bin b = *i;
  if (b.VacantSpace() >= r->h) {
    mi->second.erase(i);
    PackToBin(&b, r);
    mi->second.insert(b);
    return true;
  }
  return false;
}

int Kp1Algo::RectType(Rect* r) {
  for (int y = 1; y <= d_; y++) {
    if (r->w <= (delta_ * double(y))) {
      return y;
    }
  }
  for (int y = d_; y >= 1; --y) {
    if (r->w <= (1 - delta_ * double(y))) {
      return d_ * 2 - y + 1;
    }
  }
}

void Kp1Algo::PackToNewShelfInFrame(Rect* r, Bin* f, int j) {
  int cj = ComplType(j);
  Bin b1(WidthType(j), u_, j, 0);
  Bin b2(WidthType(cj), u_, cj, 0);
  PackToBin(f, &b1);
  PackToBin(&b1, r);
  b2.y = b1.y;
  b2.x = b1.x + b1.w;
  bins_[j].insert(b1);
  bins_[cj].insert(b2);
}

int Kp1Algo::ComplType(int t) {
	return 2 * d_ - t + 1;
}

double Kp1Algo::WidthType(int t) {
	if (t <= d_) {
		return delta_ * double(t);
	}
	return 1.0 - double(ComplType(t)) * delta_;
}