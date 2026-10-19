// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createPbehaviorModal');

module.exports = modalCreator(modalSelector, {
  elements: elementsWrapperCreator(modalSelector, {
    cancelButton: sel('createPbehaviorCancelButton'),
    submitButton: sel('createPbehaviorSubmitButton'),
  }),
});
