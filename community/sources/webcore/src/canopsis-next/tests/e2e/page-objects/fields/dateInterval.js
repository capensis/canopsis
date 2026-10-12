// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  selectRange(index = 1) {
    return this.customClick('@intervalRange')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickStartDateButton() {
    return this.customClick('@startDateButton');
  },

  clickDatePickerDayTab() {
    return this.customClick('@datePickerDateTab');
  },

  selectCalendarDay(value) {
    return this.customClickXpath(this.el('@dateDayField', value));
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

  clickEndDateButton() {
    return this.customClick('@endDateButton');
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

  clearStartDate() {
    return this.customClearValue('@startDateField');
  },

  clickStartDate() {
    return this.customClick('@startDateField');
  },

  setStartDate(value) {
    return this.customSetValue('@startDateField', value);
  },

  clickOutsideDateInterval() {
    return this.customClickOutside('@dateTimePickerCalendar');
  },

  el,
};

module.exports = {
  elements: {
    dateDayField: './/*[@data-test=\'dateTimePickerCalendar\']//div[@class=\'date-time-picker__body\']//div[@class=\'v-btn__content\' and contains(text(), \'%s\')]',

    dateTimePickerCalendar: sel('dateTimePickerCalendar'),

    datePickerDateTab: sel('datePickerDateTab'),
    datePickerHoursTab: sel('datePickerHoursTab'),
    datePickerMinutesTab: sel('datePickerMinutesTab'),
    datePickerClock: '.v-time-picker-clock__inner .v-time-picker-clock__item:nth-of-type(%s)',

    startDateField: `${sel('intervalStart')} ${sel('dateTimePickerTextField')}`,
    startDateButton: `${sel('intervalStart')} ${sel('dateTimePickerButton')}`,
    startDateCalendar: `${sel('intervalStart')} ${sel('dateTimePickerCalendar')}`,

    endDateField: `${sel('intervalStop')} ${sel('dateTimePickerTextField')}`,
    endDateButton: `${sel('intervalStop')} ${sel('dateTimePickerButton')}`,
    endDateCalendar: `${sel('intervalStop')} ${sel('dateTimePickerCalendar')}`,

    intervalRange: `${sel('intervalRange')} .v-input__control`,

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
};
