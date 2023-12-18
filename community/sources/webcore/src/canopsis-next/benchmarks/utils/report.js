const BenchmarkLauncher = require('./BenchmarkLauncher');
const ConsoleReporter = require('./ConsoleReporter');
const FileReporter = require('./FileReporter');

const benchmarkLauncher = new BenchmarkLauncher();

const runBenchmarks = async (options) => {
  await benchmarkLauncher.run(options);

  benchmarkLauncher.addReporter(
    new ConsoleReporter(),
    new FileReporter(),
  );

  benchmarkLauncher.report();
};

module.exports = {
  benchmark: benchmarkLauncher.benchmark,
  runBenchmarks,
};
