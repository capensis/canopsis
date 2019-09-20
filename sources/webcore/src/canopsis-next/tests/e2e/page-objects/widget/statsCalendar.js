// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitStatsCalendar() {
    return this.customClick('@submitStatsCalendar');
  },

  el,
};

module.exports = {
  elements: {
    submitStatsCalendar: sel('submitStatsCalendarButton'),
  },
  commands: [commands],
};
