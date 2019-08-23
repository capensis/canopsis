// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  setTabTitleField(value) {
    return this.customSetValue('@tabTitleField', value);
  },
  clickSubmitButton() {
    return this.customClick('@tabSubmitButton');
  },
};

const modalSelector = sel('createTabModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      tabTitleField: sel('tabTitleField'),
      tabSubmitButton: sel('tabSubmitButton'),
    }),
  },
  commands: [commands],
});
