// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createWatcherModal');

const commands = {
  clickName() {
    return this.customClick('@createWatcherDisplayName');
  },

  clearName() {
    return this.customClearValue('@createWatcherDisplayName');
  },

  setName(value) {
    return this.customSetValue('@createWatcherDisplayName', value);
  },

  clickTemplate() {
    return this.customClick('@createWatcherOutputTemplate');
  },

  clearTemplate() {
    return this.customClearValue('@createWatcherOutputTemplate');
  },

  setTemplate(value) {
    return this.customSetValue('@createWatcherOutputTemplate', value);
  },

  clickCreateWatcherTab(index) {
    return this.customClick(this.el('@createWatcherTab', index));
  },
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      cancelButton: sel('createWatcherCancelButton'),
      submitButton: sel('createWatcherSubmitButton'),
    }),

    createWatcherTab: `${sel('createWatcherTab')}:nth-of-type(%s)`,
    createWatcherDisplayName: sel('createWatcherDisplayName'),
    createWatcherOutputTemplate: sel('createWatcherOutputTemplate'),
  },
  commands: [commands],
});
