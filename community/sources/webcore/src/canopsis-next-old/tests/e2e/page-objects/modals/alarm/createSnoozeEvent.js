// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createSnoozeEventModal');

const commands = {
  clickDurationValue() {
    return this.customClick('@createSnoozeEventDurationValue');
  },

  clearDurationValue() {
    return this.customClearValue('@createSnoozeEventDurationValue');
  },

  setDurationValue(value) {
    return this.customSetValue('@createSnoozeEventDurationValue', value);
  },

  setDurationType(index) {
    return this.customClick('@createSnoozeEventDurationType')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('createSnoozeEventCancelButton'),
      submitButton: sel('createSnoozeEventSubmitButton'),
    }),
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    createSnoozeEventDurationValue: `${sel('durationField')} ${sel('durationValue')}`,
    createSnoozeEventDurationType: `${sel('durationField')} ${sel('durationType')} .v-select__slot`,
  },
  commands: [commands],
});
