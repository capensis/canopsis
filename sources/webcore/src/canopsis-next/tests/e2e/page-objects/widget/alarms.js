// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitAlarms() {
    return this.customClick('@submitAlarms');
  },

  clickAdvancedSettings() {
    return this.customClick('@advancedSettings')
      .defaultPause();
  },

  clickDefaultSortColumn() {
    return this.customClick('@defaultSortColumn')
      .defaultPause();
  },

  selectSortColumn(index = 1) {
    return this.customClick('@sortColumn')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  selectSortOrder(index = 1) {
    return this.customClick('@sortOrder')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickColumnNames() {
    return this.customClick('@columnNames')
      .defaultPause();
  },

  clickColumnUp(index = 1) {
    return this.customClick(this.el('@columnUp', index));
  },

  clickColumnDown(index = 1) {
    return this.customClick(this.el('@columnDown', index));
  },

  clickColumnClose(index = 1) {
    return this.customClick(this.el('@columnClose', index));
  },

  clearColumnLabel(index = 1) {
    return this.customClearValue(this.el('@columnLabel', index));
  },

  setColumnLabel(index = 1, value) {
    return this.customSetValue(this.el('@columnLabel', index), value);
  },

  clearColumnValue(index = 1) {
    return this.customClearValue(this.el('@columnValue', index));
  },

  setColumnValue(index = 1, value) {
    return this.customSetValue(this.el('@columnValue', index), value);
  },

  clickColumnHtml(index = 1) {
    return this.customClick(this.el('@columnHtml', index));
  },

  clickColumnAdd(index = 1) {
    return this.customClick(this.el('@columnAdd', index));
  },

  el,
};

module.exports = {
  elements: {
    advancedSettings: sel('advancedSettings'),
    submitAlarms: sel('submitAlarms'),
    defaultSortColumn: sel('defaultSortColumn'),
    sortColumn: `${sel('sortContainer')} div.v-input:nth-child(1)`,
    sortOrder: `${sel('sortContainer')} div.v-input:nth-child(2)`,
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    columnNames: sel('columnNames'),
    columnUp: `${sel('settings-column-%s')} ${sel('upButton')}`,
    columnDown: `${sel('settings-column-%s')} ${sel('downButton')}`,
    columnClose: `${sel('settings-column-%s')} ${sel('closeButton')}`,
    columnLabel: `${sel('settings-column-%s')} ${sel('labelField')}`,
    columnValue: `${sel('settings-column-%s')} ${sel('valueField')}`,
    columnHtml: `${sel('settings-column-%s')} div${sel('htmlSwitch')} .v-input--selection-controls__ripple`,
    columnAdd: `${sel('columnNames')} ${sel('addButton')}`,
  },
  commands: [commands],
};
