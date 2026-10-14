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
    const { activeEditModeButton } = this.elements;

    this.api.element(activeEditModeButton.locateStrategy, activeEditModeButton.selector, ({ status }) => {
      if (status === -1) {
        this.customClick('@editModeButton');
      }
    });

    return this;
  },

  clickSettingsViewButton() {
    const { activeSettingsViewButton } = this.elements;

    this.api.element(activeSettingsViewButton.locateStrategy, activeSettingsViewButton.selector, ({ status }) => {
      if (status === -1) {
        this.customClick('@settingsViewButton');
      }
    });

    return this;
  },
};

module.exports = {
  elements: {
    addViewButton: sel('addViewButton'),
    editModeButton: sel('editModeButton'),
    activeEditModeButton: `.v-btn--active${sel('editModeButton')}`,
    settingsViewButton: sel('settingsViewButton'),
    activeSettingsViewButton: `.v-btn--active${sel('settingsViewButton')}`,
    settingsWrapper: sel('settingsWrapper'),
    controlsWrapper: `${sel('settingsWrapper')} .v-speed-dial__list`,
  },
  commands: [commands],
};
