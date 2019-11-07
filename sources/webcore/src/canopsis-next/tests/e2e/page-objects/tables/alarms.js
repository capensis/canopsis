// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  openLiveReporting() {
    return this.customClick('@liveReportingButton');
  },

  clickResetLiveReporting() {
    return this.customClick('@resetAlarmsDateIntervalButton');
  },

  el,
};

module.exports = {
  elements: {
    liveReportingButton: sel('alarmsDateInterval'),
    resetAlarmsDateIntervalButton: `${sel('resetAlarmsDateInterval')} .v-chip__close`,
  },
  commands: [commands],
};
