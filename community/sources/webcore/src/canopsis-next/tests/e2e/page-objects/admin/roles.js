// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@dataTable')
      .assert.visible('@addButton');
  },

  verifyPageRoleBefore(id) {
    return this.waitForElementVisible('@dataTable')
      .assert.visible(sel(`role-${id}`));
  },

  verifyMassDeleteButton() {
    return this.waitForElementVisible('@massDeleteButton');
  },

  clickAddButton() {
    return this.customClick('@addButton');
  },

  clickOptionCheckbox(id) {
    return this.customClick(this.el('@optionCheckbox', sel(`role-${id}`)));
  },

  clickEditButton(id) {
    return this.customClick(this.el('@editButton', sel(`role-${id}`)));
  },

  clickDeleteButton(id) {
    return this.customClick(this.el('@deleteButton', sel(`role-${id}`)));
  },

  clickMassDeleteButton() {
    return this.customClick('@massDeleteButton');
  },

  clickPrevButton() {
    return this.customClick('@prevButton');
  },

  clickNextButton() {
    return this.customClick('@nextButton');
  },

  clickRefreshButton() {
    return this.customClick('@refreshButton');
  },

  setSearchingText(value) {
    return this.customSetValue('@searchingTextField', value);
  },

  clickSubmitSearchButton() {
    return this.customClick('@submitSearchButton');
  },

  clickClearSearchButton() {
    return this.customClick('@clearSearchButton');
  },

  selectRange(index = 5) {
    return this.customClick('@selectRangeField')
      .waitForElementVisible(this.el('@selectRangeItemOption', index))
      .customClick(this.el('@selectRangeItemOption', index));
  },

  el,
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/roles`;
  },
  elements: {
    dataTable: '.v-datatable',
    dataTableUserItem: '.v-datatable tbody tr',
    searchingTextField: sel('searchingTextField'),
    submitSearchButton: sel('submitSearchButton'),
    clearSearchButton: sel('clearSearchButton'),
    addButton: sel('addButton'),
    editButton: `%s ${sel('editButton')}`,
    deleteButton: `%s ${sel('deleteButton')}`,
    massDeleteButton: sel('massDeleteButton'),
    refreshButton: sel('refreshButton'),
    selectRangeField: '.v-datatable .v-datatable__actions__select .v-input__control',
    selectRangeItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    prevButton: '.v-datatable .v-datatable__actions__range-controls .v-btn[aria-label="Previous page"]',
    nextButton: '.v-datatable .v-datatable__actions__range-controls .v-btn[aria-label="Next page"]',
    optionCheckbox: `%s .v-input${sel('optionCheckbox')} .v-input--selection-controls__ripple`,
  },
  commands: [commands],
};
