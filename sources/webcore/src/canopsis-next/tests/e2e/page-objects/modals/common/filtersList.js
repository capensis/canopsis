// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('filtersListModal');

const commands = {
  clickOutside() {
    return this.customClickOutside('@filtersList');
  },
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('submitButton'),
    }),

    filtersList: sel('filtersListModal'),
  },
  commands: [commands],
});
