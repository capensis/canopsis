// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../helpers/el');

const commands = {
  verifyPageElementsBeforeById(id) {
    return this.waitForElementVisible(this.el('@viewPageById', id));
  },

  moveTab(id) {
    return this.dragAndDrop(
      this.el('@tab', id),
      this.el('@draggableWrap'),
    );
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

  clickAddWidgetButton() {
    return this.waitForElementVisible('@addWidgetButton')
      .customClick('@addWidgetButton');
  },

  clickEditViewButton() {
    return this.waitForElementVisible('@editViewButton')
      .customClick('@editViewButton');
  },

  clickTab(id) {
    return this.waitForElementVisible(this.el('@tab', id))
      .customClick(this.el('@tab', id));
  },

  clickEditTab(id) {
    return this.waitForElementVisible(this.el('@editTab', id))
      .customClick(this.el('@editTab', id));
  },

  clickCopyTab(id) {
    return this.waitForElementVisible(this.el('@copyTab', id))
      .customClick(this.el('@copyTab', id));
  },

  clickDeleteTab(id) {
    return this.waitForElementVisible(this.el('@deleteTab', id))
      .customClick(this.el('@deleteTab', id));
  },

  clickSubmitMoveTab() {
    return this.waitForElementVisible('@submitMoveTab')
      .customClick('@submitMoveTab');
  },

  verifySettingsWrapperBefore() {
    return this.waitForElementVisible('@settingsWrapper')
      .assert.visible('@settingsViewButton');
  },

  clickDeleteRowButton() {
    return this.customClick('@deleteRowButton');
  },

  clickDeleteWidgetButton() {
    return this.customClick('@deleteWidgetButton');
  },

  clickCopyWidgetButton() {
    return this.customClick('@copyWidgetButton');
  },

  clickEditWidgetButton() {
    return this.customClick('@editWidgetButton');
  },
  el,
};

module.exports = {
  elements: {
    deleteRowButton: sel('deleteRowButton'),
    deleteWidgetButton: sel('deleteWidgetButton'),
    copyWidgetButton: sel('copyWidgetButton'),
    editWidgetButton: sel('editWidgetButton'),
    viewPageById: sel('view-page-%s'),
    controlViewLayout: `${sel('controlViewLayout')} .v-speed-dial`,
    menuViewButton: `${sel('controlViewLayout')} .v-speed-dial ${sel('menuViewButton')}`,
    addWidgetButton: `${sel('controlViewLayout')} .v-speed-dial__list ${sel('addWidgetButton')}`,
    addViewButton: `${sel('controlViewLayout')} .v-speed-dial__list ${sel('addViewButton')}`,
    editViewButton: `${sel('controlViewLayout')} .v-speed-dial__list ${sel('editViewButton')}`,
    submitMoveTab: sel('submitMoveTab'),
    tab: sel('tab-%s'),
    editTab: `${sel('tab-%s')} ${sel('editTab')}`,
    copyTab: `${sel('tab-%s')} ${sel('copyTab')}`,
    deleteTab: `${sel('tab-%s')} ${sel('deleteTab')}`,
    draggableWrap: sel('draggable-wrap'),
  },
  commands: [commands],
};
