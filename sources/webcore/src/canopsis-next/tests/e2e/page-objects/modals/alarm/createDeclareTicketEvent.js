// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createDeclareTicketEventModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('declareTicketEventCancelButton'),
      submitButton: sel('declareTicketEventSubmitButton'),
    }),
  },
});
