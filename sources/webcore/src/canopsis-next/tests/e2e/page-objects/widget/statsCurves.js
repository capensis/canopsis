// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitStatsCurves() {
    return this.customClick('@submitStatsCurvesSubmitButton');
  },

  el,
};

module.exports = {
  elements: {
    submitStatsCurvesSubmitButton: sel('submitStatsCurvesSubmitButton'),
  },
  commands: [commands],
};
