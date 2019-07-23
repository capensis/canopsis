// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../helpers/el');

const commands = {
  verifyPageElementsBeforeById(id) {
    return this.waitForElementVisible(this.el('@viewPageById', id));
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
    controlViewLayout: sel('controlViewLayout'),
    menuViewButton: `${sel('controlViewLayout')} ${sel('menuViewButton')}`,
    controlViewSpeedDial: '.v-speed-dial',
    controlViewSpeedDialList: '.v-speed-dial .v-speed-dial__list',
  },
  commands: [commands],
};
