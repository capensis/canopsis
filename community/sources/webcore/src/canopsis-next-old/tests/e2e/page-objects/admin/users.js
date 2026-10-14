// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');
const sel = require('../../helpers/sel');

function getUserSelector(id) {
  return sel(`user-${id}`);
}

const commands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@dataTable')
      .assert.visible('@addButton');
  },

  verifyPageUserBefore(id) {
    return this.waitForElementVisible('@dataTable')
      .assert.visible(getUserSelector(id));
  },

  verifyMassDeleteButton() {
    return this.waitForElementVisible('@massDeleteButton');
  },

  clickAddButton() {
    return this.customClick('@addButton');
  },

  clickOptionCheckbox(id) {
    return this.customClick(this.el('@optionCheckbox', getUserSelector(id)));
  },

  clickEditButton(id) {
    return this.customClick(this.el('@editButton', getUserSelector(id)));
  },

  clickDeleteButton(id) {
    return this.customClick(this.el('@deleteButton', getUserSelector(id)));
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

  setSearchingText(value) {
    return this.customSetValue('@searchingTextField', value);
  },

  clickSubmitSearchButton() {
    return this.customClick('@submitSearchButton');
  },

  clickClearSearchButton() {
    return this.customClick('@clearSearchButton');
  },

  selectRange(idx = 5) {
    return this.customClick('@selectRangeField')
      .waitForElementVisible(this.el('@selectRangeItemOption', idx))
      .customClick(this.el('@selectRangeItemOption', idx));
  },

  el,
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/users`;
  },
  elements: {
    dataTable: '.v-datatable',
    dataTableUserItem: '.v-datatable tbody tr',
    addButton: sel('addButton'),
    searchingTextField: sel('searchingTextField'),
    submitSearchButton: sel('submitSearchButton'),
    clearSearchButton: sel('clearSearchButton'),
    editButton: `%s ${sel('editButton')}`,
    deleteButton: `%s ${sel('deleteButton')}`,
    massDeleteButton: sel('massDeleteButton'),
    selectRangeField: '.v-datatable .v-datatable__actions__select .v-input__control',
    selectRangeItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    prevButton: '.v-datatable .v-datatable__actions__range-controls .v-btn[aria-label="Previous page"]',
    nextButton: '.v-datatable .v-datatable__actions__range-controls .v-btn[aria-label="Next page"]',
    optionCheckbox: `%s .v-input${sel('optionCheckbox')} .v-input--selection-controls__ripple`,
  },
  commands: [commands],
};
