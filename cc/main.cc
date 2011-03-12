#include <iomanip>

#include <pthread.h>

#include "strip_packing.h"

__thread random_data __random_data;

struct SolutionInfo {
  double h;
  double area;
  bool valid;
  SolutionInfo(double h, double area, bool valid) : h(h), area(area), 
      valid(valid) {}
};



pthread_mutex_t sols_mu;
std::vector<SolutionInfo> sols;

static const int nworkers = 8;
Context* context[nworkers];

void* Worker(void* arg) {
  bool valid;
  int ind = reinterpret_cast<long long>(arg);
  
  context[ind]->ThreadInit();

  int cnt = context[ind]->opt.t / nworkers;
  if (ind < (context[ind]->opt.t % nworkers)) {
    ++cnt;
  }

  for (int y = 0; y < cnt; ++y) {
    context[ind]->InitAlgo();
    context[ind]->algo->Pack(context[ind]->opt.n, 0, 0, context[ind]);
    if ((0 == y) && (0 == ind)) {
      context[ind]->Render();
    }
    valid = context[ind]->Validate();
    
    pthread_mutex_lock(&sols_mu);
    sols.push_back(SolutionInfo(context[ind]->algo->solution_height,
                   context[ind]->algo->total_area, valid));
    pthread_mutex_unlock(&sols_mu);
    
    context[ind]->DestroyAlgo();
  }
  return NULL;
}


int main(int argc, char** argv) {
  std::cout << std::fixed;
  std::cout << std::setprecision(4);
  
  timespec ts;
  clock_gettime(CLOCK_REALTIME, &ts);
  srand(ts.tv_sec + ts.tv_nsec);
  
  pthread_mutex_init(&sols_mu, NULL);
  for (int y = 0; y < nworkers; ++y) {
    context[y] = new Context(argc, argv);
  }
  
  pthread_t t[nworkers];
  pthread_attr_t attr;
  pthread_attr_init(&attr);
  pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_JOINABLE);
  for (int y = 0; y < nworkers; ++y) {
    pthread_create(&t[y], &attr, Worker, reinterpret_cast<void*>(y));
  }
  for (int y = 0; y < nworkers; ++y) {
    pthread_join(t[y], NULL);
  }

  double coef_s = 0;
  double divisor;
  double sh;
  double vacant_area;
  double coef;
  if (context[0]->opt.sqrt_divisor) {
    divisor = powl(context[0]->opt.n, 1.0 / 2.0);
  } else {
    divisor = powl(context[0]->opt.n, 2.0 / 3.0);
  }
  std::cout << "n = " << context[0]->opt.n << "\n";
  if (context[0]->opt.sqrt_divisor)  {
    std::cout << "n^(1/2) = ";
  } else {
    std::cout << "n^(2/3) = ";
  }
  std::cout << divisor << "\n";
  int y = 0;
  for (std::vector<SolutionInfo>::iterator i = sols.begin(); 
       i != sols.end(); ++i, ++y) {
    std::cout << "=============\niteration " << y << "\n=============\n";
    sh = i->h;
    std::cout << "height of solution = " << sh << "\n";
    std::cout << "total area = " << i->area << "\n";
    vacant_area = double(context[0]->opt.m) * sh - i->area;
    std::cout << "uncovered area = " << vacant_area << "\n";
    coef = vacant_area / divisor;
    std::cout << "coefficient = " << coef << "\n";
    coef_s += coef;
    if (context[0]->opt.validate) {
      std::cout << "Validation: ";
      if (i->valid) {
        std::cout << "OK\n";
      } else {
        std::cout << "FAIL\n";
      }
    }
  }
  
  std::cout << "Average coef = " << coef_s / double(context[0]->opt.t) << "\n";
}