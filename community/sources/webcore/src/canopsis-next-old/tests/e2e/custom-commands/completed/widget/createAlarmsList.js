// http://nightwatchjs.org/guide#usage

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
  const commonWidget = this.page.widget.common();
  const liveReportingModal = this.page.modals.common.liveReporting();
  const dateIntervalField = this.page.fields.dateInterval();

  this.completed.widget.setCommonFields({
    ...fields,
    parameters: {
      advanced: true,
      ...parameters,
    },
  });

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
    liveReportingModal.verifyModalOpened();

    if (liveReporting.range) {
      dateIntervalField.selectRange(liveReporting.range);
    }

    if (liveReporting.calendarStartDate) {
      dateIntervalField.clickStartDateButton()
        .clickDatePickerDayTab()
        .selectCalendarDay(liveReporting.calendarStartDate.day)
        .clickDatePickerHoursTab()
        .selectCalendarHour(liveReporting.calendarStartDate.hour)
        .clickDatePickerMinutesTab()
        .selectCalendarMinute(liveReporting.calendarStartDate.minute);
    }

    if (liveReporting.calendarEndDate) {
      dateIntervalField.clickEndDateButton()
        .selectCalendarDay(liveReporting.calendarEndDate.day)
        .clickDatePickerHoursTab()
        .selectCalendarHour(liveReporting.calendarEndDate.hour)
        .clickDatePickerMinutesTab()
        .selectCalendarMinute(liveReporting.calendarEndDate.minute);
    }

    if (liveReporting.endDate) {
      dateIntervalField.clearEndDate()
        .clickEndDate()
        .setEndDate(liveReporting.endDate);
    }

    if (liveReporting.startDate) {
      dateIntervalField.clearStartDate()
        .clickStartDate()
        .setStartDate(liveReporting.startDate);
    }

    liveReportingModal.clickSubmitButton()
      .verifyModalClosed();
  }

  commonWidget.waitFirstUserPreferencesXHR(
    () => alarms.clickSubmitAlarms(),
    ({ responseData: response, requestData: request }) => callback({
      response,
      request,
    }),
  );
};
