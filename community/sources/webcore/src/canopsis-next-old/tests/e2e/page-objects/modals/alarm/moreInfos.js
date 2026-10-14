// https://nightwatchjs.org/guide/#working-with-page-objects

const { modalCreator } = require('../../../helpers/page-object-creators');
const el = require('../../../helpers/el');

const modalSelector = sel('moreInfosModal');

const commands = {
  clickOutside() {
    return this.customClickOutside('@modalSelector');
  },

  getContentText(callback) {
    return this.getText('@moreInfosContent', ({ value }) => callback(value));
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    modalSelector,
    moreInfosContent: sel('moreInfosContent'),
  },
  commands: [commands],
});
