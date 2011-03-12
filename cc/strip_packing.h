#ifndef __STRIP_PACKING_H__
#define __STRIP_PACKING_H__

#include <algorithm>
#include <cmath>
#include <cstdlib>
#include <cstring>
#include <ctime>
#include <iostream>
#include <map>
#include <set>
#include <string>
#include <vector>

#include <cairo.h>
#include <gtk/gtk.h>

// Main classes.

struct Rect {
  double x, y, w, h;
  Rect() {}
  Rect(double x, double y, double w, double h) : x(x), y(y), w(w), h(h) {}
  inline void PrintInfo() {
    std::cout << "x,y,w,h = " << x << " " << y << " " << w << " " << h << "\n";
  }
};

struct Bin : public Rect {
  double top;
  Bin() {}
  Bin(double w, double h, double top, double x = 0, double y = 0) : 
      Rect(x, y, w, h), top(top) {}
  Bin(Rect r, double top) : Rect(r), top(top) {}
  double VacantSpace() const;
};

typedef std::vector<Rect> RectVec;
typedef RectVec::iterator RectIter;

struct Context;
double rand_double();

struct SavedRect {
  Rect r;
  std::string color;
  bool solid;
  SavedRect(Rect r, std::string c, bool s) : r(r), color(c), solid(s) {}
};

class Algorithm {
 public:
  std::vector<SavedRect> saved_rects;
  double total_area;
  double solution_height;
  Algorithm() : total_area(0), solution_height(0) {}
  virtual void Pack(int n, double xbe, double ybe, Context* context) = 0;
  inline void NextRect(Rect* r) {
    r->h = rand_double();
    r->w = rand_double();
    total_area += r->w * r->h;
  }
  inline void SaveRect(Rect* r) {
    saved_rects.push_back(SavedRect(*r, "black", true));
  }
  inline void RecalcSolutionHeightSingle(Rect* r) {
    solution_height = std::max<double>(solution_height, r->y + r->h);
  }
  inline void RecalcSolutionHeight() {
    for (std::vector<SavedRect>::iterator i = saved_rects.begin(); 
         i != saved_rects.end(); ++i) {
      RecalcSolutionHeightSingle(&i->r);
    }
  }
};

class Renderer {
 public:
  GtkWidget* window_;
  cairo_t* cairo_;
  
  void ShowAll(Context* context);
};

struct Options {
  int n;
  int m;
  int t;
  bool render;
  bool render_bins;
  bool validate;
  bool save_rects;
  bool sqrt_divisor;
  std::string algo;
  Options() : n(100), m(1), render(false), render_bins(false), algo("kp1"), 
    validate(false), save_rects(false), t(1), sqrt_divisor(false) {}
  void Parse(int argc, char** argv);
};

struct Context {
  Algorithm* algo;
  Renderer* ren;
  Options opt;
  Context(int argc, char** argv);
  int seed;

  void InitAlgo();
  void DestroyAlgo();
  void Render();
  bool Validate();
  void ThreadInit();
};

// Algorithms definitions.

typedef std::map<int, std::set<Bin> > MapOfSets;
typedef std::set<Bin> SetOfBins;

class Kp1Algo : public Algorithm {
 public:
  virtual void Pack(int n, double xbe, double ybe, Context* context);
  void InitParams(int n);
  int RectType(Rect* r);
  int ComplType(int t);
  double WidthType(int t);
  void PackToNewShelfInFrame(Rect* r, Bin* f, int j);
  bool PackToTopBin(Rect* r, int j);
  void SaveBins(Context* context);

 protected:
  MapOfSets bins_;
  double delta_, u_;
  int d_;
  Bin frame_;
};

class Kp2MspBalanced : public Kp1Algo {
 public:
  void Pack(int n, double xbe, double ybe, Context* context);
  void PackWithSize(int n, double xbe, double ybe, Context* context);
  void InitFrames(double xbe, double ybe, int m);

 private:
  std::set<Bin> frames_;
};

struct PyramidPos {
  std::set<Rect>::iterator i;
  double pos;
  PyramidPos() {}
  PyramidPos(std::set<Rect>::iterator i, double pos) : i(i), pos(pos) {}
};

class PyramidAlgo : public Algorithm {
 public:
  void Pack(int n, double xbe, double ybe, Context* context);
  void ConvertCooToComplPyramid(Rect* r, int ind);
  void PackOnTop(Rect* r);
  PyramidPos PackToPyramid(std::set<Rect>* s, Rect* r);
  void PerformPacking(std::set<Rect>* s, PyramidPos* p, Rect* r);
  void InitParams(int n);

 protected:
  double shift_;
  double h_;
  double top_h_;
};

class SimplePyramidAlgo : public Algorithm {
 public:
  void InitParams(int n);
  void Pack(int n, double xbe, double ybe, Context* context);
  void ConvertCooToComplPyramid(Rect* r, int ind);
  double YOfType(int t);
  double WOfType(int t);
  int TypeOfW(double w);
  void PackToPyramid(Bin* bins, Rect* r);
  void PackOnTop(Rect* r);

 private:
  double step_;
  int nsteps_;
  double h_;
  double top_h_;
  double w_step_;
};

// Utility & inline functions.
extern __thread random_data __random_data;

inline double rand_double() {
  int32_t rand;
  random_r(&__random_data, &rand);
  return rand / (float(RAND_MAX) + 1);
}

inline void PackToBin(Bin *bin, Rect* r) {
	r->x = bin->x;
	r->y = bin->y + bin->top;
	bin->top += r->h;
}

inline bool double_less(double a, double b) {
  return (a + 1e-5) <= b;
}

inline bool double_eq(double a, double b) {
  return !double_less(a, b) && !double_less(b, a);
}

#endif  // __STRIP_PACKING_H__