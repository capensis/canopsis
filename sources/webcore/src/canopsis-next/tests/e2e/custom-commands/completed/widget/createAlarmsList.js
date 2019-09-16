// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createAlarmsList({
  parameters: {
    ack,
    enableHtml = false,
    liveReporting,
    ...parameters
  } = {},
  ...fields
}, callback = () => {}) {
  const alarms = this.page.widget.alarms();
  const liveReportingModal = this.page.modals.common.liveReporting();

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

  if (liveReporting) {
    alarms.clickCreateLiveReporting();
    liveReportingModal.verifyModalOpened()
      .selectRange(liveReporting.range);

    if (liveReporting.calendarStartDate) {
      liveReportingModal.clickStartDateButton()
        .clickDatePickerDayTab()
        .selectCalendarDay(liveReporting.calendarStartDate.day)
        .clickDatePickerHoursTab()
        .selectCalendarHour(liveReporting.calendarStartDate.hour)
        .clickDatePickerMinutesTab()
        .selectCalendarMinute(liveReporting.calendarStartDate.minute);
    }

    if (liveReporting.calendarEndDate) {
      liveReportingModal.clickEndDateButton()
        .selectCalendarDay(liveReporting.calendarEndDate.day)
        .clickDatePickerHoursTab()
        .selectCalendarHour(liveReporting.calendarEndDate.hour)
        .clickDatePickerMinutesTab()
        .selectCalendarMinute(liveReporting.calendarEndDate.minute);
    }

    if (liveReporting.endDate) {
      liveReportingModal.clearEndDate()
        .clickEndDate()
        .setEndDate(liveReporting.endDate);
    }

    if (liveReporting.startDate) {
      liveReportingModal.clearStartDate()
        .clickStartDate()
        .setStartDate(liveReporting.startDate);
    }

    liveReportingModal.clickSubmitButton()
      .verifyModalClosed();
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
