// https://nightwatchjs.org/guide/#working-with-page-objects

const { modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('variablesHelpModal');

const commands = {
  clickCopyFile(path) {
    return this.customClick(this.el('@variablesHelpCopyFile', path));
  },

  clickOutsideModal() {
    return this.moveTo('@modalSelector', -20, -20)
      .mouseButtonDown(0)
      .mouseButtonUp(0);
  },
};

module.exports = modalCreator(modalSelector, {
  elements: {
    modalSelector,
    variablesHelpCopyFile: sel('variablesHelpCopyFile-%s'),
  },
  commands: [commands],
});
