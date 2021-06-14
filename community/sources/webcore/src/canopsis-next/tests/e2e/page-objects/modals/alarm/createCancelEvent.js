// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createCancelEventModal');

const commands = {
  clickTicketNote() {
    return this.customClick('@createCancelEventNote');
  },

  clearTicketNote() {
    return this.customClearValue('@createCancelEventNote');
  },

  setTicketNote(value) {
    return this.customSetValue('@createCancelEventNote', value);
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('createCancelEventCancelButton'),
      submitButton: sel('createCancelEventSubmitButton'),
    }),

    createCancelEventNote: sel('createCancelEventNote'),
  },
  commands: [commands],
});
