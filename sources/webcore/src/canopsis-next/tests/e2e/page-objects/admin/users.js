// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');
const sel = require('../../helpers/sel');

const commands = {
  verifyPageElementsBefore() {
    return this.waitForElementVisible('@dataTable')
      .assert.visible('@addButton');
  },

  verifyPageUserBefore(user) {
    return this.waitForElementVisible('@dataTable')
      .assert.visible(this.sel(user));
  },

  verifyCreateConfirmModal() {
    return this.waitForElementVisible('@createConfirmModal')
      .assert.visible('@confirmButton');
  },

  verifyMassDeleteButton() {
    return this.waitForElementVisible('@massDeleteButton');
  },

  clickAddButton() {
    return this.customClick('@addButton');
  },

  clickOptionCheckbox(user) {
    return this.customClick(this.el('@optionCheckbox', this.sel(user)));
  },

  clickEditButton(user) {
    return this.customClick(this.el('@editButton', this.sel(user)));
  },

  clickDeleteButton(user) {
    return this.customClick(this.el('@deleteButton', this.sel(user)));
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
    createConfirmModal: sel('createConfirmModal'),
    createUserModal: sel('createUserModal'),
  },
  commands: [commands],
};
