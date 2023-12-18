const Benchmark = require('./Benchmark');
const { enhanceBenchmarkFunction } = require('./iterate');

class BenchmarkLauncher {
  constructor() {
    this.benchmarksQueue = [];
    this.reporters = [];

    enhanceBenchmarkFunction(this.benchmark);
  }

  getReportData() {
    return this.benchmarksQueue.reduce((acc, benchmark) => {
      acc[benchmark.name] = benchmark.getReportData();

      return acc;
    }, {});
  }

  addReporter(...reporters) {
    this.reporters.push(...reporters);
  }

  report() {
    const data = this.getReportData();

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
