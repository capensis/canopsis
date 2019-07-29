// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');


const commands = {
  el,

  verifySettingsWrapperBefore() {
    return this.waitForElementVisible('@settingsWrapper')
      .assert.visible('@settingsViewButton');
  },

  verifyControlsWrapperBefore() {
    return this.waitForElementVisible('@controlsWrapper')
      .assert.visible('@addViewButton');
  },

  clickAddViewButton() {
    return this.customClick('@addViewButton');
  },

  clickEditModeButton() {
    return this.customClick('@editModeButton');
  },

  clickSettingsViewButton() {
    return this.customClick('@settingsViewButton');
  },
};


module.exports = {
  elements: {
    addViewButton: sel('addViewButton'),
    editModeButton: sel('editModeButton'),
    settingsViewButton: sel('settingsViewButton'),
    settingsWrapper: '.v-speed-dial',
    controlsWrapper: '.v-speed-dial .v-speed-dial__list',
  },
  commands: [commands],
};
