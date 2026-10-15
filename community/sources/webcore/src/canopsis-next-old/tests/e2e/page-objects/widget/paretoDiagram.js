// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitParetoDiagram() {
    return this.customClick('@paretoDiagramSubmitButton');
  },

  el,
};

module.exports = {
  elements: {
    paretoDiagramSubmitButton: sel('paretoDiagramSubmitButton'),
  },
  commands: [commands],
};
