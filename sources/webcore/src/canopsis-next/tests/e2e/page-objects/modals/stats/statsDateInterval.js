// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');
const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('statsDateIntervalModal');

const commands = {
  selectPeriodUnit(index = 1) {
    return this.customClick('@intervalPeriodUnit')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  selectRange(index = 1) {
    return this.customClick('@intervalRange')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clearPeriodValue() {
    return this.customClearValue('@intervalPeriodValue');
  },

  clickPeriodValue() {
    return this.customClick('@intervalPeriodValue');
  },

  setPeriodValue(value) {
    return this.customSetValue('@intervalPeriodValue', value);
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('statsDateIntervalCancelButton'),
      submitButton: sel('statsDateIntervalSubmitButton'),
    }),

    intervalPeriodValue: `input${sel('intervalPeriodValue')}`,

    intervalPeriodUnit: `${sel('intervalPeriodUnit')} .v-input__control`,

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
