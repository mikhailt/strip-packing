#include <cstdlib>

#include "strip_packing.h"

Context::Context(int argc, char** argv) {
  timespec ts;
  
  clock_gettime(CLOCK_REALTIME, &ts);
  srand(ts.tv_sec + ts.tv_nsec);
  opt.Parse(argc, argv);
}

void Context::InitAlgo(std::string name) {
  if ("kp1" == name) {
    algo = new Kp1Algo;
  }
}

void Context::Render() {
  if (!opt.render) {
    return;
  }
  ren = new Renderer;
  ren->Init(this);
}