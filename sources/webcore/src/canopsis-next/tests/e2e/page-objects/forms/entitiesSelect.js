// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickRemoveEntity(id) {
    return this.customClick(this.el('@entitiesSelect', this.el('@removeEntity', id)));
  },

  clickClearEntities() {
    return this.customClick(this.el('@entitiesSelect', this.el('@clearEntities')));
  },

  clickSearch() {
    return this.customClick(this.el('@entitiesSelect', this.el('@searchEntity')));
  },

  clearSearch() {
    return this.customClearValue(this.el('@entitiesSelect', this.el('@searchEntity')));
  },

  setSearch(value) {
    return this.customSetValue(
      this.el('@entitiesSelect', this.el('@searchEntity')),
      value,
    );
  },

  clickSubmitSearch() {
    return this.customClick(this.el('@entitiesSelect', this.el('@submitSearchButton')));
  },

  clickAddCollection() {
    return this.customClick(this.el('@entitiesSelect', this.el('@addCollectionEntities')));
  },

  setAllCheckbox(checked) {
    return this.getAttribute(
      this.el('@entitiesSelect', this.el('@allCheckboxInput')),
      'aria-checked',
      ({ value }) => {
        if (value !== String(checked)) {
          this.customClick(this.el('@entitiesSelect', this.el('@allCheckbox')));
        }
      },
    );
  },

  setRowCheckbox(id, checked) {
    return this.getAttribute(
      this.el('@entitiesSelect', this.el('@rowCheckboxInput', id)),
      'aria-checked',
      ({ value }) => {
        if (value !== String(checked)) {
          this.customClick(this.el('@entitiesSelect', this.el('@rowCheckbox', id)));
        }
      },
    );
  },

  clickAddEntity(id) {
    return this.customClick(this.el('@entitiesSelect', this.el('@addEntityButton', id)));
  },

  setItemPerPage(index) {
    return this.customClick(this.el('@entitiesSelect', this.el('@itemsPerPage')))
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickPreviousPage() {
    return this.customClick(this.el('@entitiesSelect', this.el('@previousPage')));
  },

  clickNextPage() {
    return this.customClick(this.el('@entitiesSelect', this.el('@nextPage')));
  },

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    entitiesSelect: `${sel('entitiesSelect')} [aria-expanded="true"] %s`,

    removeEntity: `${sel('removeEntity-%s')} .v-chip__close`,
    clearEntities: sel('clearEntities'),
    searchEntity: sel('searchEntity'),
    searchEntityButton: sel('submitSearchButton'),
    addCollectionEntities: sel('addCollectionEntities'),
    allCheckboxInput: `${sel('contextEntitiesTable')} thead tr th:first-of-type .v-input--selection-controls__input input`,
    allCheckbox: `${sel('contextEntitiesTable')} thead tr th:first-of-type .v-input--selection-controls__input`,
    rowCheckbox: `${sel('contextEntitiesTable')} ${sel('contextRowCheckbox-%s')}`,
    rowCheckboxInput: `${sel('contextEntitiesTable')} ${sel('contextRowCheckbox-%s')} input`,
    addEntityButton: `${sel('contextEntitiesTable')} ${sel('contextRowAdd-%s')}`,
    itemsPerPage: `${sel('contextEntitiesTable')} .v-datatable__actions .v-select__slot`,
    previousPage: `${sel('contextEntitiesTable')} .v-datatable__actions__range-controls [aria-label="Previous page"]`,
    nextPage: `${sel('contextEntitiesTable')} .v-datatable__actions__range-controls [aria-label="Next page"]`,
  },
  commands: [commands],
};
