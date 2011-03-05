#include "strip_packing.h"

double PyramidAlgo::Pack(int n, double xbe, double ybe, Context* context) {
  shift_ = pow(n, 0.5);
  h_ = n / 4;
  top_h_ = 0;

  Rect r;
  std::set<Rect> s[2];
  for (int y = 0; y < n; ++y) {
    NextRect(&r);
    if (!PackToPyramid(&s[0], &r)) {
      if (PackToPyramid(&s[1], &r)) {
        ConvertCooToComplPyramid(&r);
      } else {
        PackOnTop(&r);
      }
    }
    SaveRect(&r);
  }
  RecalcSolutionHeight();
  if (context->opt.render_bins) {
    saved_rects.push_back(
        SavedRect(Rect(0, 0, 1, h_ + 2 * shift_), "red", false));
  }
  return solution_height;
}

void AppendAtBottom(std::set<Rect>* s, std::set<Rect>::iterator* i, Rect* r) {
  Rect b = **i;
  s->erase(*i);
  b.h += r->h;
  b.y -= r->h;
  r->x = 0;
  r->y = b.y;
  s->insert(b);
}

void AppendAtTop(std::set<Rect>* s, std::set<Rect>::iterator* i, Rect* r) {
  Rect b = **i;
  s->erase(*i);
  r->x = 0;
  r->y = b.y + b.h;
  b.h += r->h;
  s->insert(b);
}

bool PyramidAlgo::PackToPyramid(std::set<Rect>* s, Rect* r) {
  std::set<Rect>::iterator i, j;
  double cur = (1 - r->w) * h_ + shift_;
  double be, en;
  
  for (i = s->lower_bound(Rect(0, cur, 0, 0)); ; ++i) {
    if (s->end() == i) {
      en = 0; // No bottom limit.
    } else {
      en = i->y + i->h;
    }
    j = i;
    if (s->begin() == j) {
      be = h_ + shift_; // No limit at top.
    } else {
      --j;
      be = j->y;
    }
    be = std::min<double>(be, cur);
    en = std::max<double>(en, cur - 2 * shift_);
    
    // [be, en] boundaries determined.
    if (double_less(r->h, be - en)) {
      if ((j != i) && double_eq(j->y, be)) {
        AppendAtBottom(s, &j, r);
        return true;
      }
      if ((i != s->end()) && double_eq(i->y + i->h, en)) {
        AppendAtTop(s, &i, r);
        return true;
      } 
      r->x = 0;
      r->y = be - r->h;
      s->insert(*r);
      return true;
    }
    if (s->end() == i) {
      break;
    }
  }
  return false;
}

bool operator< (const Rect& a, const Rect& b) {
  return (a.y + a.h) > (b.y + b.h);
}

void PyramidAlgo::ConvertCooToComplPyramid(Rect* r) {
  Rect a;
  r->x = 1 - r->w;
  r->y = h_ + 2 * shift_ - r->y - r->h;
}

void PyramidAlgo::PackOnTop(Rect* r) {
  r->x = 0;
  r->y = top_h_ + h_ + 2 * shift_;
  top_h_ += r->h;
}