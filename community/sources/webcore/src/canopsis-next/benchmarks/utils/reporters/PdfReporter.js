const fs = require('fs');
const path = require('path');

const { jsPDF: JsPDF } = require('jspdf');

const { getMetricsGroupedData } = require('../metrics-files');
const { logInfo } = require('../logger');

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
              const [maxValue, minValue] = [firstValue, secondValue].sort((a, b) => b - a);

              const diffNumber = ((maxValue / minValue) * 100) - 100;

              diff = `${firstValue > secondValue ? '+' : '-'}${diffNumber.toFixed()}%`;
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

        const diffWidth = 70;
        const nameWidth = 340;
        const versionWidth = nameWidth;

        const width = diffWidth + nameWidth + (versionWidth * 2);
        const height = 1000;
        const doc = new JsPDF('p', 'px', [width, height]);

        const name = `Compare ${firstVersionName} and ${secondVersionName}(${property})`;

        doc.text(name, 10, 20, { align: 'left' });
        doc.table(
          10,
          30,
          rows,
          [
            { name: 'name', width: nameWidth },
            { name: firstVersionName, width: versionWidth },
            { name: secondVersionName, width: versionWidth },
            { name: 'diff', width: diffWidth },
          ],
          {
            fontSize: 14,
          },
        );

        if (!fs.existsSync(this.outputPath)) {
          fs.mkdirSync(this.outputPath);
          logInfo(`Create folder: ${this.outputPath}`);
        }

        const filePath = path.resolve(this.outputPath, `${name}.pdf`);

        doc.save(filePath);
        logInfo(`Save file: ${filePath}`);
      });
    });
  }
}

module.exports = PdfReporter;
