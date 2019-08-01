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

  clickEditViewButton() {
    return this.customClick('@editViewButton');
  },

  clickSettingsViewButton() {
    return this.customClick('@settingsViewButton');
  },
};


module.exports = {
  elements: {
    addViewButton: sel('addViewButton'),
    editViewButton: sel('editViewButton'),
    settingsViewButton: sel('settingsViewButton'),
    settingsWrapper: '.v-speed-dial',
    controlsWrapper: '.v-speed-dial .v-speed-dial__list',
  },
  commands: [commands],
};
