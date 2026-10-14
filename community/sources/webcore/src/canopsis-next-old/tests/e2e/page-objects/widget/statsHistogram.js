// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitStatsHistogram() {
    return this.customClick('@submitStatsHistogramButton');
  },

  el,
};

module.exports = {
  elements: {
    submitStatsHistogramButton: sel('submitStatsHistogramButton'),
  },
  commands: [commands],
};
