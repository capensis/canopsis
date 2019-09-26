// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitStatsCurves() {
    return this.customClick('@statsCurvesSubmitButton');
  },

  el,
};

module.exports = {
  elements: {
    statsCurvesSubmitButton: sel('statsCurvesSubmitButton'),
  },
  commands: [commands],
};
