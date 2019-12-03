// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  openLiveReporting() {
    return this.customClick('@liveReportingButton');
  },

  clickResetLiveReporting() {
    return this.customClick('@resetAlarmsDateIntervalButton');
  },

  clickOnRowInfoPopupOpen(id) {
    return this.customClick(this.el('@tableRowInfoPopupOpen', id));
  },

  clickOnRowInfoPopupClose(id) {
    return this.customClick(this.el('@tableRowInfoPopupClose', id));
  },

  verifyRowInfoPopupVisible(id) {
    return this.assert.visible(this.el('@alarmInfoPopup', id));
  },

  verifyRowInfoPopupDeleted(id) {
    return this.waitForElementNotPresent(this.el('@alarmInfoPopup', id));
  },

  verifyAlarmTimeLineVisible(id) {
    return this.assert.visible(this.el('@alarmTimeLine', id));
  },

  verifyAlarmTimeLineDeleted(id) {
    return this.waitForElementNotPresent(this.el('@alarmTimeLine', id));
  },

  el,
};

module.exports = {
  elements: {
    liveReportingButton: sel('alarmsDateInterval'),
    resetAlarmsDateIntervalButton: `${sel('resetAlarmsDateInterval')} .v-chip__close`,
    alarmTimeLine: sel('alarmTimeLine-%s'),

    alarmInfoPopup: sel('alarmInfoPopup-%s'),

    tableRowInfoPopupOpen: `${sel('alarmInfoPopup-%s')} ${sel('alarmInfoPopupOpenButton')}`,
    tableRowInfoPopupClose: `${sel('alarmInfoPopup-%s')} ${sel('alarmInfoPopupCloseButton')}`,
  },
  commands: [commands],
};
