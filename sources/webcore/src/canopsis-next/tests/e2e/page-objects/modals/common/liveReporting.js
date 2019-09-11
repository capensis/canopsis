// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');
const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('liveReportingModal');

const commands = {
  clickStartDateButton() {
    return this.customClick('@startDateButton');
  },

  clickEndDateButton() {
    return this.customClick('@endDateButton');
  },

  clearStartDate() {
    return this.customClearValue('@startDateField');
  },

  clickStartDate() {
    return this.customClick('@startDateField');
  },

  clickDatePickerDateTab() {
    return this.customClick('@datePickerDateTab');
  },

  selectCalendarDay(value) {
    this.api.useXpath()
      .click(this.el('@dateDayField', value));
    return this;
  },

  clickDatePickerHoursTab() {
    return this.customClick('@datePickerHoursTab');
  },

  selectCalendarHour(index) {
    this.api
      .moveToElement(this.el('@datePickerClock', index + 1), 5, 5)
      .mouseButtonDown(0)
      .mouseButtonUp(0)
      .pause(500);

    return this;
  },

  clickDatePickerMinutesTab() {
    return this.customClick('@datePickerMinutesTab');
  },

  selectCalendarMinute(index) {
    this.api
      .moveToElement(this.el('@datePickerClock', index + 1), 5, 5)
      .mouseButtonDown(0)
      .mouseButtonUp(0)
      .pause(500);

    return this;
  },

  setStartDate(value) {
    return this.customSetValue('@startDateField', value);
  },

  clearEndDate() {
    return this.customClearValue('@endDateField');
  },

  clickEndDate() {
    return this.customClick('@endDateField');
  },

  setEndDate(value) {
    return this.customSetValue('@endDateField', value);
  },

  selectRange(index = 1) {
    return this.customClick('@intervalRange')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('liveReportingApplyButton'),
      cancelButton: sel('liveReportingCancelButton'),
    }),
    dateDayField: './/*[@data-test=\'dateTimePickerCalendar\']//div[@class=\'date-time-picker__body\']//div[@class=\'v-btn__content\' and contains(text(), \'%s\')]',

    datePickerDateTab: sel('datePickerDateTab'),
    datePickerHoursTab: sel('datePickerHoursTab'),
    datePickerMinutesTab: sel('datePickerMinutesTab'),
    datePickerSecondsTab: sel('datePickerSecondsTab'),
    datePickerClock: '.v-time-picker-clock__inner .v-time-picker-clock__item:nth-of-type(%s)',

    startDateField: `${sel('intervalStart')} ${sel('dateTimePickerTextField')}`,
    startDateButton: `${sel('intervalStart')} ${sel('dateTimePickerButton')}`,
    startDateCalendar: `${sel('intervalStart')} ${sel('dateTimePickerCalendar')}`,

    endDateField: `${sel('intervalStop')} ${sel('dateTimePickerTextField')}`,
    endDateButton: `${sel('intervalStop')} ${sel('dateTimePickerButton')}`,
    endDateCalendar: `${sel('intervalStop')} ${sel('dateTimePickerCalendar')}`,
    endDateDayField: './/div[@class=\'v-date-picker-table\']/[@class=\'v-btn__content\' and contains(text(), "%s")]',

    intervalRange: `${sel('intervalRange')} .v-input__control`,

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
