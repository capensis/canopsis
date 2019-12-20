// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const modalSelector = sel('createChangeStateEventModal');

const commands = {
  clickNote() {
    return this.customClick('@createChangeStateEventNote');
  },

  clearNote() {
    return this.customClearValue('@createChangeStateEventNote');
  },

  setNote(value) {
    return this.customSetValue('@createChangeStateEventNote', value);
  },

  clickCriticity(index) {
    return this.customClick(this.el('@createChangeStateCriticity', index));
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('createChangeStateEventCancelButton'),
      submitButton: sel('createChangeStateEventSubmitButton'),
    }),

    createChangeStateEventNote: sel('createChangeStateEventNote'),
    createChangeStateCriticity: sel('stateCriticity-%s'),
  },
  commands: [commands],
});
