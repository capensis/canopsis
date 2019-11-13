// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickAddInfo() {
    return this.customClick('@addManageInfos');
  },

  setItemPerPage(index) {
    return this.customClick('@itemsPerPage')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickPreviousPage() {
    return this.customClick('@previousPage');
  },

  clickNextPage() {
    return this.customClick('@nextPage');
  },

  clickRowEditInfo(value) {
    return this.customClick(this.el('@manageInfosRowEditInfo', value));
  },

  clickRowDeleteInfo(value) {
    return this.customClick(this.el('@manageInfosRowDeleteInfo', value));
  },

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    addManageInfos: sel('addManageInfos'),
    manageInfosRowEditInfo: `${sel('manageInfosTable')} ${sel('infoEditButton-%s')}`,
    manageInfosRowDeleteInfo: `${sel('manageInfosTable')} ${sel('infoDeleteButton-%s')}`,
    itemsPerPage: `${sel('manageInfosTable')} .v-datatable__actions .v-select__slot`,
    previousPage: `${sel('manageInfosTable')} .v-datatable__actions__range-controls [aria-label="Previous page"]`,
    nextPage: `${sel('manageInfosTable')} .v-datatable__actions__range-controls [aria-label="Next page"]`,
  },
  commands: [commands],
};
