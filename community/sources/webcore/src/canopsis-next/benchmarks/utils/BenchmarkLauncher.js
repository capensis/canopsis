const Benchmark = require('./Benchmark');
const { enhanceBenchmarkFunction } = require('./iterate');

class BenchmarkLauncher {
  constructor() {
    this.benchmarksQueue = [];
    this.reporters = [];

    enhanceBenchmarkFunction(this.benchmark);
  }

  toJSON() {
    return this.benchmarksQueue.map(benchmark => benchmark.toJSON());
  }

  addReporter(...reporters) {
    this.reporters.push(...reporters);
  }

  report() {
    const data = this.toJSON();

    this.reporters.forEach(reporter => reporter.report(data));
  }

  benchmark = (benchmarkName, benchmarkFunc) => {
    this.benchmarksQueue.push(new Benchmark(benchmarkName, benchmarkFunc));
  };

  async run(options) {
    await this.benchmarksQueue.reduce(
      (acc, benchmark) => acc.then(() => benchmark.run(options)),
      Promise.resolve(),
    );
  }
}

module.exports = BenchmarkLauncher;
