// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');
const sel = require('../../helpers/sel');
const { FILTERS_TYPE } = require('../../constants');

const commands = {
  clickSearchInput() {
    return this.customClick('@searchInput');
  },

  clearSearchInput() {
    return this.customClearValue('@searchInput');
  },

  setSearchInput(value) {
    return this.customSetValue('@searchInput', value);
  },

  clickSearchButton() {
    return this.customClick('@searchButton');
  },

  clickSearchResetButton() {
    return this.customClick('@resetSearchButton');
  },

  moveToSearchInformation() {
    return this.moveToElement('@helpInformationButton', 5, 5);
  },

  clickNextPageTopPagination() {
    return this.customClick('@topPaginationPrevious');
  },

  clickPreviousPageTopPagination() {
    return this.customClick('@topPaginationNext');
  },

  clickNextPageBottomPagination() {
    return this.customClick('@bottomPaginationPrevious');
  },

  clickPreviousPageBottomPagination() {
    return this.customClick('@bottomPaginationPrevious');
  },

  clickOnPageBottomPagination(page) {
    this.api
      .useXpath()
      .customClick(this.el('@bottomPaginationPage', page));

    return this;
  },

  setMixFilters(checked = false) {
    return this.getAttribute('@mixFiltersInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@mixFilters');
      }
    });
  },

  selectFilter(index = 1) {
    return this.customClick('@selectFilters')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  setFiltersType(type) {
    return this.getAttribute('@andFiltersInput', 'aria-checked', ({ value }) => {
      if (value === 'true' && type === FILTERS_TYPE.OR) {
        this.customClick('@orFilters');
      } else if (value === 'false' && type === FILTERS_TYPE.AND) {
        this.customClick('@andFilters');
      }
    });
  },

  showFiltersList() {
    return this.customClick('@showFiltersListButton');
  },

  setItemPerPage(index) {
    return this.customClick('@itemsPerPage')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickTableHeaderCell(header) {
    this.api
      .useXpath()
      .customClick(this.el('@tableHeaderCell', header));

    return this;
  },

  setAllCheckbox(checked) {
    return this.getAttribute('@selectAllCheckboxInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@selectAllCheckbox');
      }
    });
  },

  setRowCheckbox(id, checked) {
    return this.getAttribute(this.el('@tableRowCheckboxInput', id), 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick(this.el('@tableRowCheckbox', id));
      }
    });
  },

  clickOnRow(id) {
    return this.customClick(this.el('@tableRow', id));
  },

  clickOnRowCell(id, column) {
    return this.customClick(this.el('@tableRow', id, column));
  },

  clickOnMassAction(index) {
    return this.customClick(this.el('@massActionsPanelItem', index));
  },

  clickOnSharedAction(id, index) {
    return this.customClick(this.el('@rowActionsSharedPanelItem', id, index));
  },

  clickOnDropDownActions(id, index) {
    return this
      .customClick(this.el('@rowMoreActionsButton', id))
      .customClick(this.el('@rowDropDownActions', index));
  },

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    tableHeaderCell: './/*[@data-test=\'tableWidget\']//thead//tr//th[@role=\'columnheader\']//span[contains(text(), \'%s\')]',
    selectAllCheckboxInput: `${sel('tableWidget')} thead tr th:first-of-type .v-input--selection-controls__input input`,
    selectAllCheckbox: `${sel('tableWidget')} thead tr th:first-of-type .v-input--selection-controls__input`,

    tableRow: `${sel('tableRow-%s')}`,
    tableRowCheckbox: `${sel('tableRow-%s')} ${sel('rowCheckbox')} ${sel('vCheckboxFunctional')}`,
    tableRowCheckboxInput: `${sel('tableRow-%s')} ${sel('rowCheckbox')} ${sel('vCheckboxFunctional')} input`,
    tableRowColumn: `${sel('tableRow-%s')} ${sel('alarmValue-%s')}`,

    searchInput: `${sel('tableSearch')} ${sel('searchingTextField')}`,
    searchButton: `${sel('tableSearch')} ${sel('submitSearchButton')}`,
    resetSearchButton: `${sel('tableSearch')} ${sel('clearSearchButton')}`,
    helpInformationButton: `${sel('tableSearch')} ${sel('tableSearchHelp')}`,
    helpInformation: `${sel('tableSearch')} ${sel('tableSearchHelpInfo')}`,

    topPaginationPrevious: `${sel('topPagination')} ${sel('paginationPreviewButton')}`,
    topPaginationNext: `${sel('topPagination')} ${sel('paginationNextButton')}`,

    bottomPaginationPrevious: `${sel('vPagination')} li:first-child`,
    bottomPaginationPage: './/*[@data-test=\'vPagination\']//button[@class=\'v-pagination__item\' and contains(text(), \'%s\')]',
    bottomPaginationNext: `${sel('vPagination')} li:last-child`,

    itemsPerPage: `${sel('tableWidget')} ${sel('itemsPerPage')} .v-select__slot`,

    mixFilters: `${sel('tableFilterSelector')} div${sel('mixFilters')} .v-input--selection-controls__ripple`,
    mixFiltersInput: `${sel('tableFilterSelector')} input${sel('mixFilters')}`,
    selectFilters: `${sel('tableFilterSelector')} ${sel('selectFilters')} .v-input__slot`,
    andFilters: `${sel('tableFilterSelector')} ${sel('andFilters')} + .v-input--selection-controls__ripple`,
    andFiltersInput: `${sel('tableFilterSelector')} input${sel('andFilters')}`,
    orFilters: `${sel('tableFilterSelector')} ${sel('orFilters')} + .v-input--selection-controls__ripple`,
    showFiltersListButton: sel('showFiltersListButton'),

    massActionsPanelItem: `${sel('tableWidget')} ${sel('massActionsPanel')} ${sel('actionsPanelItem')}:nth-of-type(%s) button`,
    rowActionsSharedPanelItem: `${sel('tableRow-%s')} ${sel('sharedActionsPanel')} .layout ${sel('actionsPanelItem')}:nth-of-type(%s) button`,
    rowMoreActionsButton: `${sel('tableRow-%s')} ${sel('sharedActionsPanel')} .layout ${sel('dropDownActionsButton')}`,
    rowDropDownActions: `.menuable__content__active ${sel('dropDownActions')} ${sel('actionsPanelItem')}:nth-of-type(%s)`,

  },
  commands: [commands],
};
