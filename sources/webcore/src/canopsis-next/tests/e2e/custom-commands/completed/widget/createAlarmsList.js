// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createAlarmsList({
  parameters: {
    ack,
    enableHtml = false,
    ...parameters
  } = {},
  ...fields
}, callback = () => {}) {
  const alarms = this.page.widget.alarms();

  this.completed.widget.setCommonFields({ ...fields, parameters });

  if (enableHtml) {
    alarms.setEnableHtml(enableHtml);
  }

  if (ack) {
    alarms.clickAckGroup()
      .setIsAckNoteRequired(ack.isAckNoteRequired)
      .setIsMultiAckEnabled(ack.isMultiAckEnabled);

    if (ack.fastAckOutput) {
      alarms.clickFastAckOutput()
        .setFastAckOutputSwitch(ack.fastAckOutput.enabled);
    }

    if (ack.fastAckOutput.enabled) {
      alarms.clickFastAckOutputText()
        .clearFastAckOutputText()
        .setFastAckOutputText(ack.fastAckOutput.output);
    }
  }

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => alarms.clickSubmitAlarms(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
