// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('confirmAckModal');

const commands = {
  clickSubmitButtonWithTicket() {
    return this.customClick('@submitButtonWithTicket');
  },
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('confirmAckCancelButton'),
      submitButton: sel('confirmAckContinueButton'),
    }),
    submitButtonWithTicket: sel('confirmAckContinueWithTicketButton'),
  },
  commands: [commands],
});
