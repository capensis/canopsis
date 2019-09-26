// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');
const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('statsDisplayModeModal');

const commands = {
  clickStatTitle() {
    return this.customClick('@statTitle');
  },

  selectDisplayModeType(index = 1) {
    return this.customClick('@statsDisplayModeType')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickDisplayModeParameterValue(level) {
    return this.customClick(this.el('@statsDisplayModeParameterValue', level));
  },

  clearDisplayModeParameterValue(level) {
    return this.customClearValue(this.el('@statsDisplayModeParameterValue', level));
  },

  setDisplayModeParameterValue(level, value) {
    return this.customSetValue(this.el('@statsDisplayModeParameterValue', level), value);
  },

  clickDisplayModeParameterColor(level) {
    return this.customClick(this.el('@statsDisplayModeParameterColor', level));
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('statsDisplayModeCancelButton'),
      submitButton: sel('statsDisplayModeSubmitButton'),
    }),

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s) .v-list__tile__content',

    statsDisplayModeType: `${sel('statsDisplayModeParameters')} .v-input__slot`,
    statsDisplayModeParameterValue: `${sel('statsDisplayMode-%s')} ${sel('displayModeValue')}`,
    statsDisplayModeParameterColor: `${sel('statsDisplayMode-%s')} ${sel('displayModeColorPicker')}`,
  },
  commands: [commands],
});
