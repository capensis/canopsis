const groupData = (arrayData, prepare) => arrayData.reduce((acc, { name, data }) => {
  acc[name] = prepare ? prepare(data) : data;

  return acc;
}, {});

const getMetricsGroupedData = (metrics) => {
  const groupedMetrics = groupData(
    metrics,
    benchmarkData => groupData(
      benchmarkData,
      groupData,
    ),
  );

  const {
    allProperties,
    allMeasures,
    allBenchmarks,
    allFileNames,
  } = Object.entries(groupedMetrics).reduce((acc, [fileName, benchmarksByName]) => {
    acc.allFileNames.push(fileName);

    Object.entries(benchmarksByName).forEach(([benchmarkName, measuresByName]) => {
      if (!acc.allBenchmarks.includes(benchmarkName)) {
        acc.allBenchmarks.push(benchmarkName);
      }

      Object.entries(measuresByName).forEach(([measureName, properties]) => {
        if (!acc.allMeasures.includes(measureName)) {
          acc.allMeasures.push(measureName);
        }

        if (!properties) {
          return;
        }

        Object.keys(properties).forEach((property) => {
          if (!acc.allProperties.includes(property)) {
            acc.allProperties.push(property);
          }
        });
      });
    });

    return acc;
  }, {
    allProperties: [],
    allMeasures: [],
    allBenchmarks: [],
    allFileNames: [],
  });

  return {
    groupedMetrics,

    allProperties,
    allMeasures,
    allBenchmarks,
    allFileNames,
  };
};

module.exports = {
  getMetricsGroupedData,
};
