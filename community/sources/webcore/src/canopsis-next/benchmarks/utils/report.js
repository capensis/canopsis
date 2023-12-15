/* eslint-disable no-console */

const benchmarksQueue = [];

const benchmark = (benchmarkName, benchmarkFunc) => {
  const metricsQueue = [];

  const measure = (measureName, measureFunc) => {
    metricsQueue.push({ name: measureName, func: measureFunc });
  };

  benchmarksQueue.push({
    func: async (options) => {
      await benchmarkFunc(measure, options);
      const benchmarkData = {};

      await metricsQueue.reduce((acc, { name, func }) => {
        const runMetric = async () => {
          console.info('Start metric: ', name);
          const report = data => benchmarkData[name] = data;

          await func(report, options);

          console.info('Finish metric: ', name);
        };

        return acc.then(runMetric);
      }, Promise.resolve());

      return benchmarkData;
    },
    name: benchmarkName,
  });
};

const runBenchmarks = async (options) => {
  const benchmarkReportedData = {};

  await benchmarksQueue.reduce(
    (acc, { name, func }) => {
      const runBenchmark = async () => {
        console.info('Start benchmark: ', name);

        benchmarkReportedData[name] = await func(options);

        console.info('Finish benchmark: ', name, '\n');
      };

      return acc.then(runBenchmark);
    },
    Promise.resolve(),
  );

  Object.entries(benchmarkReportedData).forEach(([name, data]) => {
    /* eslint-disable no-console */
    console.info(name);
    console.table(data);
  });

  return benchmarkReportedData;
};

module.exports = {
  benchmark,
  runBenchmarks,
};
