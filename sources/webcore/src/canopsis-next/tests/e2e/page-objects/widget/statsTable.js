// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitStatsTable() {
    return this.customClick('@submitStatsTable');
  },

  getCellName(id, callback) {
    return this.getText(this.el('@statsTableNameCell', id), callback);
  },

  getCellValue(id, callback) {
    return this.getText(this.el('@statsTableValue', id), callback);
  },

  getCellSubValue(id, callback) {
    return this.getText(this.el('@statsTableSubValue', id), callback);
  },

  getAlarmsChipsValue(id, callback) {
    return this.getText(this.el('@alarmsChipsValue', id), callback);
  },

  el,
};

module.exports = {
  elements: {
    submitStatsTable: sel('submitStatsTable'),
    statsTableNameCell: sel('statsTableNameCell-%s'),
    statsTableValue: `${sel('statsTableCell-%s')} ${sel('statsTableValue')}`,
    statsTableSubValue: `${sel('statsTableCell-%s')} ${sel('statsTableSubValue')}`,
    alarmsChipsValue: `${sel('statsTableCell-%s')} ${sel('alarmsChipsValue')}`,
  },
  commands: [commands],
};
