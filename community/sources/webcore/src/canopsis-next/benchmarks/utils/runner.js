const BenchmarkLauncher = require('./BenchmarkLauncher');
const ChartsReporter = require('./reporters/ChartsReporter');
const ConsoleReporter = require('./reporters/ConsoleReporter');
const FileReporter = require('./reporters/FileReporter');

const benchmarkLauncher = new BenchmarkLauncher();

const runBenchmarks = async (options) => {
  const { jsonName } = options;

  await benchmarkLauncher.run(options);

  benchmarkLauncher.addReporter(
    new ConsoleReporter(),
    new FileReporter({ name: jsonName }),
  );

  benchmarkLauncher.report();
};

const saveReportsCharts = () => {
  const chartsReporter = new ChartsReporter({ width: 500, height: 500 });

  chartsReporter.report(...FileReporter.readMetricFiles());
};

module.exports = {
  benchmark: benchmarkLauncher.benchmark,
  runBenchmarks,
  saveReportsCharts,
};
