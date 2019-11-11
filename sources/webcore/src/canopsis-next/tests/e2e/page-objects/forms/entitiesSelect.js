// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickRemoveEntity(id) {
    return this.customClick(this.el('@removeEntity', id));
  },

  clickClearEntities() {
    return this.customClick('@clearEntities');
  },

  clickSearch() {
    return this.customClick('@searchEntity');
  },

  clearSearch() {
    return this.customClearValue('@searchEntity');
  },

  setSearch(value) {
    return this.customSetValue('@searchEntity', value);
  },

  clickSubmitSearch() {
    return this.customClick('@searchEntityButton');
  },

  clickAddCollection() {
    return this.customClick('@addCollectionEntities');
  },

  setAllCheckbox(checked) {
    return this.getAttribute('@allCheckboxInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@allCheckbox');
      }
    });
  },

  setRowCheckbox(id, checked) {
    return this.getAttribute(this.el('@rowCheckboxInput', id), 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick(this.el('@rowCheckbox', id));
      }
    });
  },

  clickAddEntity(id) {
    return this.customClick(this.el('@addEntityButton', id));
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

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    removeEntity: `${sel('removeEntity-%s')} .v-chip__close`,
    clearEntities: sel('clearEntities'),
    searchEntity: sel('searchEntity'),
    searchEntityButton: sel('searchEntityButton'),
    addCollectionEntities: sel('addCollectionEntities'),
    allCheckboxInput: `${sel('contextEntitiesTable')} thead tr th:first-of-type .v-input--selection-controls__input input`,
    allCheckbox: `${sel('contextEntitiesTable')} thead tr th:first-of-type .v-input--selection-controls__input`,
    rowCheckbox: `${sel('contextEntitiesTable')} ${sel('contextRowCheckbox-%s')}`,
    addEntityButton: `${sel('contextEntitiesTable')} ${sel('contextRowAdd-%s')}`,
    itemsPerPage: `${sel('contextEntitiesTable')} .v-datatable__actions .v-select__slot`,
    previousPage: `${sel('contextEntitiesTable')} .v-datatable__actions__range-controls [aria-label="Previous page"]`,
    nextPage: `${sel('contextEntitiesTable')} .v-datatable__actions__range-controls [aria-label="Next page"]`,
  },
  commands: [commands],
};
