// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../helpers/el');

const commands = {
  verifyPageElementsBeforeById(id) {
    return this.waitForElementVisible(this.el('@viewPageById', id));
  },

  clickMenuViewButton() {
    return this.waitForElementVisible('@controlViewLayout')
      .assert.visible('@menuViewButton')
      .customClick('@menuViewButton');
  },

  clickAddViewButton() {
    return this.waitForElementVisible('@addViewButton')
      .customClick('@addViewButton');
  },

  verifySettingsWrapperBefore() {
    return this.waitForElementVisible('@settingsWrapper')
      .assert.visible('@settingsViewButton');
  },
  el,
};

module.exports = {
  elements: {
    viewPageById: sel('view-page-%s'),
    controlViewLayout: `${sel('controlViewLayout')} .v-speed-dial`,
    menuViewButton: `${sel('controlViewLayout')} .v-speed-dial ${sel('menuViewButton')}`,
    addViewButton: `${sel('controlViewLayout')} .v-speed-dial__list ${sel('addViewButton')}`,
  },
  commands: [commands],
};
