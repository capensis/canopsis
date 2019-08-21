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

  clickColumnAdd() {
    return this.customClick('@columnAdd');
  },

  clickDefaultNumberOfElementsPerPage() {
    return this.customClick('@defaultNumberOfElementsPerPage')
      .defaultPause();
  },

  selectElementsPerPage(index = 1) {
    return this.customClick('@elementsPerPage')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickFilterOnOpenResolved() {
    return this.customClick('@filterOnOpenResolved')
      .defaultPause();
  },

  clickOpenFilter() {
    return this.customClick('@openFilter');
  },

  clickResolvedFilter() {
    return this.customClick('@resolvedFilter');
  },

  clickFilters() {
    return this.customClick('@filters')
      .defaultPause();
  },

  clickMixFilters() {
    return this.customClick('@mixFilters');
  },

  clickAndFilters() {
    return this.customClick('@andFilters');
  },

  clickOrFilters() {
    return this.customClick('@orFilters');
  },

  selectFilters(index = 1) {
    return this.customClick('@selectFilters')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickInfoPopupButton() {
    return this.customClick('@infoPopupButton');
  },

  clickCreateEditMore() {
    return this.customClick('@createEditMore');
  },

  clickDeleteMore() {
    return this.customClick('@deleteMore');
  },

  clickEnableHtml() {
    return this.customClick('@enableHtml');
  },

  clickAckGroup() {
    return this.customClick('@ackGroup')
      .defaultPause();
  },

  clickIsAckNoteRequired() {
    return this.customClick('@isAckNoteRequired');
  },

  clickIsMultiAckEnabled() {
    return this.customClick('@isMultiAckEnabled');
  },

  clickFastAckOutput() {
    return this.customClick('@fastAckOutput')
      .defaultPause();
  },

  clickFastAckOutputSwitch() {
    return this.customClick('@fastAckOutputSwitch');
  },

  setFastAckOutputText(value) {
    return this.customSetValue('@fastAckOutputText', value);
  },

  el,
};

module.exports = {
  elements: {
    advancedSettings: sel('advancedSettings'),
    submitAlarms: sel('submitAlarms'),
    defaultSortColumn: sel('defaultSortColumn'),
    sortColumn: `${sel('sortContainer')} div.v-input:nth-child(1) .v-input__slot`,
    sortOrder: `${sel('sortContainer')} div.v-input:nth-child(2) .v-input__slot`,
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    columnNames: sel('columnNames'),
    columnUp: `${sel('settings-column-%s')} ${sel('upButton')}`,
    columnDown: `${sel('settings-column-%s')} ${sel('downButton')}`,
    columnClose: `${sel('settings-column-%s')} ${sel('closeButton')}`,
    columnLabel: `${sel('settings-column-%s')} ${sel('labelField')}`,
    columnValue: `${sel('settings-column-%s')} ${sel('valueField')}`,
    columnHtml: `${sel('settings-column-%s')} div${sel('htmlSwitch')} .v-input--selection-controls__ripple`,
    columnAdd: `${sel('columnNames')} ${sel('addButton')}`,
    defaultNumberOfElementsPerPage: sel('defaultNumberOfElementsPerPage'),
    elementsPerPage: `${sel('elementsPerPage')} .v-input__slot`,
    filterOnOpenResolved: sel('filterOnOpenResolved'),
    openFilter: `div${sel('openFilter')} .v-input--selection-controls__ripple`,
    resolvedFilter: `div${sel('resolvedFilter')} .v-input--selection-controls__ripple`,
    filters: sel('filters'),
    addFilter: sel('addFilter'),
    editFilter: sel('editFilter-%s'),
    deleteFilter: sel('deleteFilter-%s'),
    mixFilters: `div${sel('mixFilters')} .v-input--selection-controls__ripple`,
    andFilters: `div${sel('andFilters')} .v-input--selection-controls__ripple`,
    orFilters: `div${sel('orFilters')} .v-input--selection-controls__ripple`,
    selectFilters: `${sel('selectFilters')} .v-input__slot`,
    infoPopupButton: sel('infoPopupButton'),
    createEditMore: sel('createEditMore'),
    deleteMore: sel('deleteMore'),
    enableHtml: `${sel('isHtmlEnabledOnTimeLine')} ${sel('switcherLayout')} .v-input--selection-controls__ripple`,
    isAckNoteRequired: `${sel('isAckNoteRequired')} ${sel('switcherLayout')} .v-input--selection-controls__ripple`,
    isMultiAckEnabled: `${sel('isMultiAckEnabled')} ${sel('switcherLayout')} .v-input--selection-controls__ripple`,
    ackGroup: sel('ackGroup'),
    fastAckOutput: sel('fastAckOutput'),
    fastAckOutputSwitch: `${sel('fastAckOutput')} .v-input--switch .v-input--selection-controls__ripple`,
    fastAckOutputText: `${sel('fastAckOutput')} .v-text-field input`,
  },
  commands: [commands],
};
