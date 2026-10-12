// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createAckEventModal');

const commands = {
  clickSubmitButtonWithTicket() {
    return this.customClick('@submitButtonWithTicket');
  },

  clickTicketNumber() {
    return this.customClick('@createAckEventTicket');
  },

  clearTicketNumber() {
    return this.customClearValue('@createAckEventTicket');
  },

  setTicketNumber(value) {
    return this.customSetValue('@createAckEventTicket', value);
  },

  clickTicketNote() {
    return this.customClick('@createAckEventNote');
  },

  clearTicketNote() {
    return this.customClearValue('@createAckEventNote');
  },

  setTicketNote(value) {
    return this.customSetValue('@createAckEventNote', value);
  },

  setAckTicketResources(checked = false) {
    return this.getAttribute('@createAckEventResourceInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@createAckEventResource');
      }
    });
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('createAckEventCancelButton'),
      submitButton: sel('createAckEventSubmitButton'),
    }),
    submitButtonWithTicket: sel('createAckEventSubmitWithTicketButton'),

    createAckEventTicket: sel('createAckEventTicket'),

    createAckEventNote: sel('createAckEventNote'),

    createAckEventResourceInput: `input${sel('createAckEventResource')}`,
    createAckEventResource: sel('createAckEventResource'),
  },
  commands: [commands],
});
