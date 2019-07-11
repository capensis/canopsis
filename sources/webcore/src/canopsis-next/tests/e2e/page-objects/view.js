// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../helpers/el');

const commands = {
  verifyPageElementsBeforeById(id) {
    return this.waitForElementVisible(this.el('@viewPageById', id));
  },

  el,
};

module.exports = {
  elements: {
    viewPageById: sel('view-page-%s'),
  },
  commands: [commands],
};
