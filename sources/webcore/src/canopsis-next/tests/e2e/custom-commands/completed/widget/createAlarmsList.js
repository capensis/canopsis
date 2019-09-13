// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createAlarmsList({
  parameters: {
    ack,
    filters,
    enableHtml = false,
    liveReporting,
    ...parameters
  } = {},
  ...fields
}, callback = () => {}) {
  const alarms = this.page.widget.alarms();
  const createFilter = this.page.modals.common.createFilterModal();
  const liveReportingModal = this.page.modals.common.liveReportingModal();

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

  if (filters) {
    alarms.clickFilters()
      .setMixFilters(filters.isMix);

    if (filters.isMix) {
      alarms.setFiltersType(filters.type);
    }

    if (filters.groups) {
      alarms.clickAddFilter();
      createFilter.verifyModalOpened()
        .clearFilterTitle()
        .setFilterTitle(filters.title)
        .fillFilterGroups(filters.groups)
        .clickSubmitButton()
        .verifyModalClosed();
    }

    if (filters.selected) {
      filters.selected.forEach((element) => {
        alarms.selectFilter(element);
      });
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
