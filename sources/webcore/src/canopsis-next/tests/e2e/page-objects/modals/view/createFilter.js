// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createFilterModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('createFilterSubmitButton'),
      cancelButton: sel('createFilterCancelButton'),
    }),
  },
  commands: [],
});
