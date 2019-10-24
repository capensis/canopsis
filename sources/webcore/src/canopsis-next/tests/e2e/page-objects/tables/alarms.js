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

  openLiveReporting() {
    return this.customClick('@liveReportingButton');
  },

  clickAlarmListHeaderCell(header) {
    this.api
      .useXpath()
      .customClick(this.el('@alarmListHeaderCell', header));

    return this;
  },

  setAllCheckbox(checked) {
    return this.getAttribute('@selectAllCheckboxInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@selectAllCheckbox');
      }
    });
  },

  setRowCheckbox(checked) {
    // TODO add set for checkbox without input
    return this.getAttribute('@selectAllCheckboxInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@selectAllCheckbox');
      }
    });
  },

  clickOnAlarmRow(id) {
    return this.customClick(this.el('@alarmListRow', id));
  },

  clickOnAlarmRowCell(id, column) {
    return this.customClick(this.el('@alarmListRow', id, column));
  },

  setItemPerPage(index) {
    return this.customClick('@itemsPerPage')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    searchInput: `${sel('alarmListSearch')} ${sel('searchingTextField')}`,
    searchButton: `${sel('alarmListSearch')} ${sel('submitSearchButton')}`,
    resetSearchButton: `${sel('alarmListSearch')} ${sel('clearSearchButton')}`,
    helpInformationButton: `${sel('alarmListSearch')} ${sel('alarmListSearchHelp')}`,
    helpInformation: `${sel('alarmListSearch')} ${sel('alarmListSearchHelp')}`,

    topPaginationPrevious: `${sel('topPagination')} ${sel('paginationPreviewButton')}`,
    topPaginationNext: `${sel('topPagination')} ${sel('paginationNextButton')}`,

    bottomPaginationPrevious: `${sel('vPagination')} li:first-child`,
    bottomPaginationPage: './/*[@data-test=\'vPagination\']//button[@class=\'v-pagination__item\' and contains(text(), \'%s\')]',
    bottomPaginationNext: `${sel('vPagination')} li:last-child`,

    itemsPerPage: `${sel('alarmsWidget')} ${sel('itemsPerPage')} .v-select__slot`,

    mixFilters: `${sel('alarmsFilterSelector')} div${sel('mixFilters')} .v-input--selection-controls__ripple`,
    mixFiltersInput: `${sel('alarmsFilterSelector')} input${sel('mixFilters')}`,
    selectFilters: `${sel('alarmsFilterSelector')} ${sel('selectFilters')} .v-input__slot`,
    andFilters: `${sel('alarmsFilterSelector')} ${sel('andFilters')} + .v-input--selection-controls__ripple`,
    andFiltersInput: `${sel('alarmsFilterSelector')} input${sel('andFilters')}`,
    orFilters: `${sel('alarmsFilterSelector')} ${sel('orFilters')} + .v-input--selection-controls__ripple`,
    showFiltersListButton: sel('showFiltersListButton'),

    liveReportingButton: sel('alarmsDateInterval'),

    alarmsListTable: sel('alarmsListTable'),

    alarmListHeaderCell: './/*[@data-test=\'alarmsListTable\']//thead//tr//th[@role=\'columnheader\']//span[contains(text(), \'%s\')]',
    selectAllCheckboxInput: `${sel('alarmsListTable')} thead tr th:first-of-type .v-input--selection-controls__input input`,
    selectAllCheckbox: `${sel('alarmsListTable')} thead tr th:first-of-type .v-input--selection-controls__input`,

    alarmListRow: `${sel('alarmListRow-%s')}`,
    alarmListRowCheckbox: `${sel('alarmListRow-%s')} ${sel('rowCheckbox')} ${sel('vCheckboxFunctional')}`,
    alarmListRowColumn: `${sel('alarmListRow-%s')} ${sel('alarmValue-%s')}`,

    actionsPanelItem: `${sel('alarmsWidget')} ${sel('massActionsPanel')} ${sel('actionsPanelItem')}:nth-of-type(%s)`,
  },
  commands: [commands],
};
