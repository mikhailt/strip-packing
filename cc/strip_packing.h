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
  virtual double Pack(int n, double xbe, double ybe, Context* context) = 0;
  inline void NextRect(Rect* r) {
    r->h = rand_double();
    r->w = rand_double();
    total_area += r->w * r->h;
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
  std::string algo;
  Options() : n(100), m(1), render(false), render_bins(false), algo("kp1"), 
    validate(false), save_rects(false), t(1) {}
  void Parse(int argc, char** argv);
};

struct Context {
  Algorithm* algo;
  Renderer* ren;
  Options opt;
  Context(int argc, char** argv);
  void InitAlgo(std::string name);
  void Render();
  void Validate();
};

// Algorithms definitions.

typedef std::map<int, std::set<Bin> > MapOfSets;
typedef std::set<Bin> SetOfBins;

class Kp1Algo : public Algorithm {
 public:
  virtual double Pack(int n, double xbe, double ybe, Context* context);
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
  double Pack(int n, double xbe, double ybe, Context* context);
  void PackWithSize(int n, double xbe, double ybe, Context* context);
  void InitFrames(double xbe, double ybe, int m);

 private:
  std::set<Bin> frames_;
};

// Utility inline functions.

inline double rand_double() {
  return rand()/(float(RAND_MAX) + 1);
}

inline void PackToBin(Bin *bin, Rect* r) {
	r->x = bin->x;
	r->y = bin->y + bin->top;
	bin->top += r->h;
}

#endif  // __STRIP_PACKING_H__