// https://nightwatchjs.org/guide/#working-with-page-objects

const { API_ROUTES } = require('@/config');

const el = require('../../helpers/el');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../constants');

const commands = {
  clickSubmitAlarms() {
    return this.customClick('@submitAlarms');
  },

  setEnableHtml(checked = false) {
    return this.getAttribute('@enableHtmlInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@enableHtml');
      }
    });
  },

  checkEnableHtmlValue(expected) {
    return this.getAttribute('@enableHtmlInput', 'aria-checked', ({ value }) => {
      this.assert.equal(value, String(expected));
    });
  },

  clickAckGroup() {
    return this.customClick('@ackGroup');
  },

  setIsAckNoteRequired(checked = false) {
    return this.getAttribute('@isAckNoteRequiredInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@isAckNoteRequired');
      }
    });
  },

  setIsMultiAckEnabled(checked = false) {
    return this.getAttribute('@isMultiAckEnabledInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@isMultiAckEnabled');
      }
    });
  },

  checkIsMultiAckEnabled(expected) {
    return this.getAttribute('@isMultiAckEnabledInput', 'aria-checked', ({ value }) => {
      this.assert.equal(value, String(expected));
    });
  },

  clickFastAckOutput() {
    return this.customClick('@fastAckOutput');
  },

  setFastAckOutputSwitch(checked = false) {
    return this.getAttribute('@fastAckOutputSwitchInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@fastAckOutputSwitch');
      }
    });
  },

  checkFastAckOutputSwitch(expected) {
    return this.getAttribute('@fastAckOutputSwitchInput', 'aria-checked', ({ value }) => {
      this.assert.equal(value, String(expected));
    });
  },

  clickFastAckOutputText() {
    return this.customClick('@fastAckOutputField');
  },

  clearFastAckOutputText() {
    return this.customClearValue('@fastAckOutputField');
  },

  setFastAckOutputText(value) {
    return this.customSetValue('@fastAckOutputField', value);
  },

  clickCreateLiveReporting() {
    return this.customClick('@liveReportingCreateButton');
  },

  clickEditLiveReporting() {
    return this.customClick('@liveReportingEditButton');
  },

  clickDeleteLiveReporting() {
    return this.customClick('@liveReportingDeleteButton');
  },

  waitFirstAlarmsListXHR(triggerFunc, callback) {
    return this.waitForFirstXHR(
      API_ROUTES.alarmList,
      WAIT_FOR_FIRST_XHR_TIME,
      triggerFunc,
      ({ responseData, requestData }) => callback({
        responseData: JSON.parse(responseData),
        requestData: JSON.parse(requestData),
      }),
    );
  },

  waitAllAlarmsListXHR(triggerFunc, callback, delay = WAIT_FOR_FIRST_XHR_TIME) {
    return this.waitForXHR(
      API_ROUTES.alarmList,
      delay,
      triggerFunc,
      callback,
    );
  },

  waitFirstEventXHR(triggerFunc, callback) {
    return this.waitForFirstXHR(
      API_ROUTES.event,
      WAIT_FOR_FIRST_XHR_TIME,
      triggerFunc,
      ({ responseData }) => callback({
        responseData: JSON.parse(responseData),
      }),
    );
  },

  waitFirstPbehaviorXHR(triggerFunc, callback) {
    return this.waitForFirstXHR(
      API_ROUTES.pbehavior.pbehavior,
      WAIT_FOR_FIRST_XHR_TIME,
      triggerFunc,
      ({ responseData }) => callback(responseData),
    );
  },

  el,
};

module.exports = {
  elements: {
    submitAlarms: sel('submitAlarms'),

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    enableHtml: `${sel('isHtmlEnabledOnTimeLine')} div${sel('switcherField')}`,
    enableHtmlInput: `${sel('isHtmlEnabledOnTimeLine')} input${sel('switcherField')}`,

    isAckNoteRequired: `${sel('isAckNoteRequired')} div${sel('switcherField')}`,
    isAckNoteRequiredInput: `${sel('isAckNoteRequired')} input${sel('switcherField')}`,
    isMultiAckEnabled: `${sel('isMultiAckEnabled')} div${sel('switcherField')}`,
    isMultiAckEnabledInput: `${sel('isMultiAckEnabled')} input${sel('switcherField')}`,
    ackGroup: sel('ackGroup'),
    fastAckOutput: sel('fastAckOutput'),
    fastAckOutputSwitch: `div${sel('fastAckOutputSwitch')} .v-input--selection-controls__ripple`,
    fastAckOutputSwitchInput: `input${sel('fastAckOutputSwitch')}`,
    fastAckOutputField: sel('fastAckOutputField'),

    liveReportingCreateButton: `${sel('liveReporting')} + div > ${sel('createButton')}`,
    liveReportingEditButton: `${sel('liveReporting')} + div > ${sel('editButton')}`,
    liveReportingDeleteButton: `${sel('liveReporting')} + div > ${sel('deleteButton')}`,
  },
  commands: [commands],
};
