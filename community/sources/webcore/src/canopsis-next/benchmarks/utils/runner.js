const path = require('path');

const BenchmarkLauncher = require('./BenchmarkLauncher');
const ChartsReporter = require('./reporters/ChartsReporter');
const ConsoleReporter = require('./reporters/ConsoleReporter');
const FileReporter = require('./reporters/FileReporter');
const PdfReporter = require('./reporters/PdfReporter');

const benchmarkLauncher = new BenchmarkLauncher();

const metricsFolderPath = path.resolve(process.cwd(), 'benchmarks', '__metrics__');
const chartsFolderPath = path.resolve(process.cwd(), 'benchmarks', '__reports__');

const runBenchmarks = async (options) => {
  const { jsonName } = options;

  await benchmarkLauncher.run(options);

  benchmarkLauncher.addReporter(
    new ConsoleReporter(),
    new FileReporter({ name: jsonName, outputPath: metricsFolderPath }),
  );

  benchmarkLauncher.report();
};

const compareMetric = ({ target, source }) => {
  const chartsReporter = new ChartsReporter({ width: 1000, height: 1000, outputPath: chartsFolderPath });
  const pdfReporter = new PdfReporter({ outputPath: chartsFolderPath });

  if (target && source) {
    const targetMetrics = FileReporter.readMetricFile(metricsFolderPath, target);
    const sourceMetrics = FileReporter.readMetricFile(metricsFolderPath, source);

    chartsReporter.report(targetMetrics, sourceMetrics);
    pdfReporter.report(targetMetrics, sourceMetrics);
  } else {
    const metrics = FileReporter.readMetricFiles(metricsFolderPath);

    chartsReporter.report(...metrics);
    pdfReporter.report(...metrics);
  }
};

module.exports = {
  benchmark: benchmarkLauncher.benchmark,
  runBenchmarks,
  compareMetric,
};
