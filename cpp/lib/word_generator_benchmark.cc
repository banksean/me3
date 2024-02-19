#include "benchmark/benchmark.h"
#include "cpp/lib/word_generator.h"

static const auto generator = sandbox::create_generator();

static void benchmark_helloworld_generator(benchmark::State& state) {
  for (auto _ : state) {
    benchmark::DoNotOptimize(generator->next(state.range(0)));
  }
}

BENCHMARK(benchmark_helloworld_generator)->RangeMultiplier(2)->Range(1, 100);