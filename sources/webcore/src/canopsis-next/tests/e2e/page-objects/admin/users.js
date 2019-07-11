// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');
const sel = require('../../helpers/sel');

const commands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@dataTable')
      .assert.visible('@addButton');
  },

  verifyPageUserBefore(userSelector) {
    return this.waitForElementVisible('@dataTable')
      .assert.visible(this.sel(userSelector));
  },

  verifyCreateConfirmModal() {
    return this.waitForElementVisible('@confirmationModal')
      .assert.visible('@confirmButton');
  },

  verifyMassDeleteButton() {
    return this.waitForElementVisible('@massDeleteButton');
  },

  clickAddButton() {
    return this.customClick('@addButton');
  },

  clickOptionCheckbox(userSelector) {
    return this.customClick(this.el('@optionCheckbox', this.sel(userSelector)));
  },

  clickEditButton(userSelector) {
    return this.customClick(this.el('@editButton', this.sel(userSelector)));
  },

  clickDeleteButton(userSelector) {
    return this.customClick(this.el('@deleteButton', this.sel(userSelector)));
  },

  clickConfirmButton() {
    return this.customClick('@confirmButton');
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

  selectRange(idx = 5) {
    return this.customClick('@selectRangeField')
      .waitForElementVisible(this.el('@selectRangeItemOption', idx))
      .customClick(this.el('@selectRangeItemOption', idx));
  },

  el,
  sel,
};

module.exports = {
  url() {
    return `${process.env.VUE_DEV_SERVER_URL}admin/users`;
  },
  elements: {
    dataTable: '.v-datatable',
    addButton: sel('addButton'),
    editButton: `%s ${sel('editButton')}`,
    deleteButton: `%s ${sel('deleteButton')}`,
    confirmButton: sel('confirmButton'),
    massDeleteButton: sel('massDeleteButton'),
    selectRangeField: '.v-datatable .v-datatable__actions__select .v-input__control',
    selectRangeItemOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    prevButton: '.v-datatable .v-datatable__actions__range-controls .v-btn[aria-label="Previous page"]',
    nextButton: '.v-datatable .v-datatable__actions__range-controls .v-btn[aria-label="Next page"]',
    optionCheckbox: `%s .v-input${sel('optionCheckbox')} .v-input--selection-controls__ripple`,
    confirmationModal: sel('confirmationModal'),
  },
  commands: [commands],
};
