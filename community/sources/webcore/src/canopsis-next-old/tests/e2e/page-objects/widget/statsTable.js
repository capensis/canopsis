// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitStatsTable() {
    return this.customClick('@submitStatsTable');
  },

  el,
};

module.exports = {
  elements: {
    submitStatsTable: sel('submitStatsTable'),
  },
  commands: [commands],
};
