// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clearName() {
    return this.customClearValue('@pbehaviorName');
  },

  clickName() {
    return this.customClick('@pbehaviorName');
  },

  setName(value) {
    return this.customSetValue('@pbehaviorName', value);
  },

  clickStartDate() {
    return this.customClick('@pbehaviorStartDate');
  },

  clickEndDate() {
    return this.customClick('@pbehaviorEndDate');
  },

  clickFilter() {
    return this.customClick('@pbehaviorFilterButton');
  },

  selectType(index = 1) {
    return this.customClick('@pbehaviorType')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clearReason() {
    return this.customClearValue('@pbehaviorReason');
  },

  clickReason() {
    return this.customClick('@pbehaviorReason');
  },

  setReason(value) {
    return this.customSetValue('@pbehaviorReason', value);
  },

  selectReason(index = 1) {
    return this.customClick('@pbehaviorReason')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickSimple() {
    return this.customClick('@pbehaviorSimpleButton');
  },

  clickAdvanced() {
    return this.customClick('@pbehaviorAdvancedButton');
  },

  selectFrequency(index = 1) {
    return this.customClick('@pbehaviorFrequency')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  selectByWeekDay(index = 1, checked = false) {
    return this.customClick('@pbehaviorByWeekDay')
      .waitForElementVisible(this.el('@optionSelect', index))
      .getAttribute(this.el('@optionSelectInput', index), 'aria-checked', ({ value }) => {
        if (value !== String(checked)) {
          this.customClick(this.el('@optionSelect', index));
        }
      });
  },

  clearRepeat() {
    return this.customClearValue('@pbehaviorRepeat');
  },

  clickRepeat() {
    return this.customClick('@pbehaviorRepeat');
  },

  setRepeat(value) {
    return this.customSetValue('@pbehaviorRepeat', value);
  },

  clearInterval() {
    return this.customClearValue('@pbehaviorInterval');
  },

  clickInterval() {
    return this.customClick('@pbehaviorInterval');
  },

  setInterval(value) {
    return this.customSetValue('@pbehaviorInterval', value);
  },

  selectWeekStart(index = 1) {
    return this.customClick('@pbehaviorWeekStart')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  selectByMonth(index = 1, checked = false) {
    return this.customClick('@pbehaviorByMonth')
      .waitForElementVisible(this.el('@optionSelect', index))
      .getAttribute(this.el('@optionSelectInput', index), 'aria-checked', ({ value }) => {
        if (value !== String(checked)) {
          this.customClick(this.el('@optionSelect', index));
        }
      });
  },

  clearBySetPosition() {
    return this.customClearValue('@pbehaviorBySetPosition');
  },

  clickBySetPosition() {
    return this.customClick('@pbehaviorBySetPosition');
  },

  setBySetPosition(value) {
    return this.customSetValue('@pbehaviorBySetPosition', value);
  },

  clearByMonthDay() {
    return this.customClearValue('@pbehaviorByMonthDay');
  },

  clickByMonthDay() {
    return this.customClick('@pbehaviorByMonthDay');
  },

  setByMonthDay(value) {
    return this.customSetValue('@pbehaviorByMonthDay', value);
  },

  clearByYearDay() {
    return this.customClearValue('@pbehaviorByYearDay');
  },

  clickByYearDay() {
    return this.customClick('@pbehaviorByYearDay');
  },

  setByYearDay(value) {
    return this.customSetValue('@pbehaviorByYearDay', value);
  },

  clearByWeekNo() {
    return this.customClearValue('@pbehaviorByWeekNo');
  },

  clickByWeekNo() {
    return this.customClick('@pbehaviorByWeekNo');
  },

  setByWeekNo(value) {
    return this.customSetValue('@pbehaviorByWeekNo', value);
  },

  clearByHour() {
    return this.customClearValue('@pbehaviorByHour');
  },

  clickByHour() {
    return this.customClick('@pbehaviorByHour');
  },

  setByHour(value) {
    return this.customSetValue('@pbehaviorByHour', value);
  },

  clearByMinute() {
    return this.customClearValue('@pbehaviorByMinute');
  },

  clickByMinute() {
    return this.customClick('@pbehaviorByMinute');
  },

  setByMinute(value) {
    return this.customSetValue('@pbehaviorByMinute', value);
  },

  clearBySecond() {
    return this.customClearValue('@pbehaviorBySecond');
  },

  clickBySecond() {
    return this.customClick('@pbehaviorBySecond');
  },

  setBySecond(value) {
    return this.customSetValue('@pbehaviorBySecond', value);
  },

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    optionSelectInput: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s) input',

    pbehaviorName: sel('pbehaviorFormName'),

    pbehaviorStartDate: sel('startDateTimePicker'),
    pbehaviorEndDate: sel('stopDateTimePicker'),

    pbehaviorFilterButton: sel('pbehaviorFilterButton'),

    pbehaviorType: `${sel('pbehaviorType')} .v-input__slot`,

    pbehaviorReason: `${sel('pbehaviorReason')} .v-input__slot`,

    ruleSwitcherInput: `input${sel('pbehaviorRuleSwitcher')}`,
    ruleSwitcher: `div${sel('pbehaviorRuleSwitcher')}`,

    pbehaviorSimpleButton: sel('pbehaviorSimple'),
    pbehaviorAdvancedButton: sel('pbehaviorAdvanced'),

    pbehaviorFrequency: `${sel('pbehaviorFrequency')} .v-input__slot`,
    pbehaviorByWeekDay: `${sel('pbehaviorByWeekDay')} .v-input__slot`,
    pbehaviorRepeat: sel('pbehaviorRepeat'),
    pbehaviorInterval: sel('pbehaviorInterval'),

    pbehaviorWeekStart: `${sel('pbehaviorWeekStart')} .v-input__slot`,
    pbehaviorByMonth: `${sel('pbehaviorWeekStart')} .v-input__slot`,
    pbehaviorBySetPosition: sel('pbehaviorBySetPos'),
    pbehaviorByMonthDay: sel('pbehaviorByMonthDay'),
    pbehaviorByYearDay: sel('pbehaviorByYearDay'),
    pbehaviorByWeekNo: sel('pbehaviorByWeekNo'),
    pbehaviorByHour: sel('pbehaviorByHour'),
    pbehaviorByMinute: sel('pbehaviorByMinute'),
    pbehaviorBySecond: sel('pbehaviorBySecond'),

    addExdateButton: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    deleteExdateButton: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    exdateField: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    addComment: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    deleteCommentButton: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    commentField: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
};
