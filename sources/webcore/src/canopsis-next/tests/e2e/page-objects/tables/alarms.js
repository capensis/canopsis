// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');
const { scopedPageObject } = require('../../helpers/page-object-creators');


const commands = {
  clickLiveReportingOpenButton() {
    return this.customClick('@liveReportingButton');
  },

  clickLiveReportingResetButton() {
    return this.customClick('@resetAlarmsDateIntervalButton');
  },

  clickOnRowInfoPopupOpenButton(id) {
    return this.customClick(this.el('@tableRowInfoPopupOpen', id));
  },

  clickOnRowInfoPopupCloseButton(id) {
    return this.customClick(this.el('@tableRowInfoPopupClose', id));
  },

  verifyRowInfoPopupVisible(id) {
    return this.assert.visible(this.el('@alarmInfoPopup', id));
  },

  verifyRowInfoPopupDeleted(id) {
    return this.assert.hidden(this.el('@alarmInfoPopup', id));
  },

  verifyAlarmTimeLineVisible(id) {
    return this.assert.visible(this.el('@alarmTimeLine', id));
  },

  verifyAlarmTimeLineDeleted(id) {
    return this.waitForElementNotPresent(this.el('@alarmTimeLine', id));
  },

  moveToExtraDetailsOpenButton(id, type) {
    this
      .moveToElement(this.el('@tableRowExtraDetailsOpen', id, type), 0, 0)
      .api.pause(1000);

    return this;
  },

  moveOutsideToExtraDetailsOpenButton(id, type) {
    this
      .moveToElement(this.el('@tableRowExtraDetailsOpen', id, type), 0, 0)
      .api.moveTo(null, -50, -50)
      .pause(500);

    return this;
  },

  verifyRowExtraDetailsVisible(id) {
    return this.assert.visible(this.el('@tableRowExtraDetailsContent', id));
  },

  verifyRowExtraDetailsDeleted(id) {
    return this.assert.hidden(this.el('@tableRowExtraDetailsContent', id));
  },

  el,
};

module.exports = scopedPageObject({
  elements: {
    liveReportingButton: sel('alarmsDateInterval'),
    resetAlarmsDateIntervalButton: `${sel('resetAlarmsDateInterval')} .v-chip__close`,
    alarmTimeLine: sel('alarmTimeLine-%s'),

    alarmInfoPopup: sel('alarmInfoPopup-%s'),

    tableRowInfoPopupOpen: `${sel('tableRow-%s')} ${sel('alarmInfoPopupOpenButton')}`,
    tableRowInfoPopupClose: `${sel('alarmInfoPopup-%s')} ${sel('alarmInfoPopupCloseButton')}`,
    tableRowInfoPopupContent: `${sel('alarmInfoPopup-%s')} ${sel('alarmInfoPopupContent')}`,

    extraDetailsPopup: sel('extraDetails-%s'),

    tableRowExtraDetailsOpen: `${sel('tableRow-%s')} ${sel('extraDetailsOpenButton-%s')}`,
    tableRowExtraDetailsContent: sel('extraDetailsContent-%s'),
  },
  commands: [commands],
});
