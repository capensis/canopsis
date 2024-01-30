const fs = require('fs');
const path = require('path');

const { jsPDF: JsPDF } = require('jspdf');

const { getMetricsGroupedData } = require('../metrics-files');

class PdfReporter {
  constructor({ outputPath } = {}) {
    this.outputPath = outputPath;
  }

  report(...metrics) {
    const { groupedMetrics, allFileNames, allProperties, allMeasures, allBenchmarks } = getMetricsGroupedData(
      metrics,
    );

    const comparingVersions = allFileNames.reduce((acc, fileName, index) => {
      const nextFileName = allFileNames[index + 1];

      if (nextFileName) {
        acc.push([fileName, nextFileName]);
      }

      return acc;
    }, []);

    comparingVersions.forEach(([firstVersionName, secondVersionName]) => {
      allProperties.forEach((property) => {
        const rows = allBenchmarks.reduce((acc, benchmarkName) => {
          allMeasures.forEach((measureName) => {
            const firstValue = groupedMetrics?.[firstVersionName]?.[benchmarkName]?.[measureName]?.[property];
            const secondValue = groupedMetrics?.[secondVersionName]?.[benchmarkName]?.[measureName]?.[property];

            let diff = '';

            if (firstValue && secondValue) {
              const diffNumber = (1 - (secondValue / firstValue)) * 100;

              diff = `${diffNumber > 0 ? '+' : ''}${diffNumber.toFixed()}%`;
            }

            const firstString = firstValue ? String(firstValue) : '-';
            const secondString = secondValue ? String(secondValue) : '-';

            acc.push({
              name: `${benchmarkName}: ${measureName}`,
              [firstVersionName]: firstString,
              [secondVersionName]: secondString,
              diff,
            });
          });

          return acc;
        }, []);

        const doc = new JsPDF();

        const name = `Compare ${firstVersionName} and ${secondVersionName}(${property})`;

        doc.text(name, 10, 10, { align: 'left' });
        doc.table(
          10,
          15,
          rows,
          ['name', firstVersionName, secondVersionName, 'diff'],
          {
            fontSize: 10,
          },
        );

        if (!fs.existsSync(this.outputPath)) {
          fs.mkdirSync(this.outputPath);
        }

        doc.save(path.resolve(this.outputPath, `${name}.pdf`));
      });
    });
  }
}

module.exports = PdfReporter;
