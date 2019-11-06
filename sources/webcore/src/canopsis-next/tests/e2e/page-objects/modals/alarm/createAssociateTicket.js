// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createAssociateTicketModal');

const commands = {
  clickTicketNumber() {
    return this.customClick('@numberOfTicket');
  },

  clearTicketNumber() {
    return this.customClearValue('@numberOfTicket');
  },

  setTicketNumber(value) {
    return this.customSetValue('@numberOfTicket', value);
  },
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('createAssociateTicketCancelButton'),
      submitButton: sel('createAssociateTicketSubmitButton'),
    }),

    numberOfTicket: sel('createAssociateTicketNumberOfTicket'),
  },
  commands: [commands],
});
