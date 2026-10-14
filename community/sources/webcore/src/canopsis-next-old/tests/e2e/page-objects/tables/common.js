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

  keyupSearchEnter() {
    return this.customKeyup('@searchInput', this.api.Keys.ENTER);
  },

  clickSearchResetButton() {
    return this.customClick('@resetSearchButton');
  },

  moveToSearchInformation() {
    return this.moveToElement('@helpInformationButton', 5, 5);
  },

  moveOutsideSearchInformation() {
    this
      .moveToElement('@helpInformationButton', 0, 0)
      .api.moveTo(null, -50, -50)
      .pause(500);

    return this;
  },

  verifySearchInformationVisible() {
    return this.assert.visible('@helpInformation');
  },

  verifySearchInformationHidden() {
    return this.assert.hidden('@helpInformation');
  },

  getTopPaginationPage(callback) {
    return this.getText('@topPaginationPage', ({ value }) => callback(value));
  },

  clickNextPageTopPagination() {
    return this.customClick('@topPaginationNext');
  },

  clickPreviousPageTopPagination() {
    return this.customClick('@topPaginationPrevious');
  },

  clickNextPageBottomPagination() {
    return this.customClick('@bottomPaginationNext');
  },

  clickPreviousPageBottomPagination() {
    return this.customClick('@bottomPaginationPrevious');
  },

  clickOnPageBottomPagination(page) {
    this.api
      .useXpath()
      .customClick(this.el('@bottomPaginationPage', page))
      .useCss();

    return this;
  },

  setMixFilters(checked = false) {
    return this.getAttribute('@mixFiltersInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@mixFilters');
      }
    });
  },

  clickFilterSelect() {
    return this.customClick('@selectFilters');
  },

  clickFilter(name) {
    this
      .clickFilterSelect()
      .api.useXpath()
      .customClick(this.el('@filterOptionXPath', name))
      .useCss();

    return this;
  },

  verifyFilterVisible(name) {
    this
      .clickFilterSelect()
      .api.useXpath()
      .assert.visible(this.el('@filterOptionXPath', name))
      .useCss();

    return this.clickOutsideFiltersOptions();
  },

  verifyFilterDeleted(name) {
    this
      .clickFilterSelect()
      .api.useXpath()
      .waitForElementNotPresent(this.el('@filterOptionXPath', name))
      .useCss();

    return this.clickOutsideFiltersOptions();
  },

  clearFilters() {
    return this.customClick('@clearFilters');
  },

  assertActiveFilters(count) {
    return this.assert.elementsCount('@activeOptionSelect', count);
  },

  clickOutsideFiltersOptions() {
    return this.customClickOutside('@activeMenu');
  },

  selectFilter(name, checked = false) {
    return this
      .customClick('@selectFilters')
      .getAttribute(
        this.el('@filterOptionInput', name),
        'aria-checked',
        ({ value }) => {
          if (value !== String(checked)) {
            this.customClick(this.el('@filterOption', name));
          }
        },
      );
  },

  clickOutsideFilter() {
    return this.customClickOutside('@selectFilters');
  },

  checkSelectedFilter(name, checked) {
    return this.getAttribute(
      this.el('@filterOptionInput', name),
      'aria-checked',
      ({ value }) => this.assert.equal(value, String(checked)),
    );
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
      .customClick(this.el('@tableHeaderCellContent', header))
      .useCss();

    return this;
  },

  getTableHeaderTextByIndex(index, callback) {
    // Increment because we have a checkbox in first column
    return this.getText(this.el('@tableHeaderCellContentByIndex', index + 1), callback);
  },

  checkTableHeaderSort(header, sortDirection) {
    this.api
      .useXpath()
      .getAttribute(this.el('@tableHeaderCell', header), 'aria-sort', ({ value }) => {
        this.assert.equal(value, sortDirection);
      })
      .useCss();

    return this;
  },

  verifyTableColumnVisible(header) {
    this.api
      .useXpath()
      .assert.visible(this.el('@tableHeaderCell', header))
      .useCss();

    return this;
  },

  verifyTableColumnDeleted(header) {
    this.api
      .useXpath()
      .waitForElementNotPresent(this.el('@tableHeaderCell', header))
      .useCss();

    return this;
  },

  setAllCheckbox(checked) {
    return this.getAttribute('@selectAllCheckboxInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@selectAllCheckbox');
      }
    });
  },

  moveOutsideMassActionsPanel() {
    this
      .moveToElement('@massActionsPanel', 0, 0)
      .api.moveTo(null, -5, -5)
      .pause(500);

    return this;
  },

  checkRowCheckboxValue(id, callback) {
    return this.getAttribute(this.el('@tableRowCheckboxInput', id), 'aria-checked', ({ value }) => {
      callback(value);
    });
  },

  setRowCheckbox(id, checked) {
    return this.checkRowCheckboxValue(id, (value) => {
      if (value !== String(checked)) {
        this.customClick(this.el('@tableRowCheckbox', id));
      }
    });
  },

  assertActiveCheckboxCount(count) {
    return this.assert.elementsCount('@tableCheckboxActiveInput', count);
  },

  clickOnRow(id) {
    return this.customClick(this.el('@tableRow', id));
  },

  clickOnRowCell(id, column) {
    return this.customClick(this.el('@tableRowColumn', id, column));
  },

  getCellTextByColumnName(id, column, callback) {
    return this.getText(this.el('@tableRowColumn', id, column), callback);
  },

  clickOnMassAction(index) {
    return this.customClick(this.el('@massActionsPanelItem', index));
  },

  clickOnSharedAction(id, index) {
    return this.customClick(this.el('@rowActionsSharedPanelItem', id, index));
  },

  verifySharedActionVisible(id, index) {
    return this.assert.visible(this.el('@rowActionsSharedPanelItem', id, index));
  },

  verifySharedActionDeleted(id, index) {
    return this.waitForElementNotPresent(this.el('@rowActionsSharedPanelItem', id, index));
  },

  clickOnDropDownDots(id) {
    return this.customClick(this.el('@rowMoreActionsButton', id));
  },

  clickOnDropDownAction(index) {
    return this.customClick(this.el('@rowDropDownAction', index));
  },

  assertDropDownAction(count) {
    return this.assert.elementsCount('@rowDropDownActions', count);
  },

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    activeOptionSelect: '.menuable__content__active .v-select-list [role="listitem"] .v-list__tile--active',

    activeMenu: '.menuable__content__active',

    filterOption: sel('filterOption-%s'),
    filterOptionXPath: './/*[contains(@class, "menuable__content__active")]//*[contains(@class, "v-select-list")]//span[contains(text(), "%s")]',
    filterOptionInput: `${sel('filterOption-%s')} input`,

    optionSelectInput: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s) input',

    tableHeaderCell: './/*[@data-test=\'tableWidget\']//thead//tr//th[@role=\'columnheader\'][span[contains(text(), \'%s\')]]',
    tableHeaderCellContent: './/*[@data-test=\'tableWidget\']//thead//tr//th[@role=\'columnheader\']//span[contains(text(), \'%s\')]',
    tableHeaderCellContentByIndex: `${sel('tableWidget')} thead tr th[role="columnheader"]:nth-of-type(%s) span`,
    selectAllCheckboxInput: `${sel('tableWidget')} thead tr th:first-of-type .v-input--selection-controls__input input`,
    selectAllCheckbox: `${sel('tableWidget')} thead tr th:first-of-type .v-input__slot`,

    tableRow: `${sel('tableRow-%s')}`,
    tableCheckboxActiveInput: `${sel('rowCheckbox')} ${sel('vCheckboxFunctional')} input[aria-checked="true"]`,
    tableRowCheckbox: `${sel('tableRow-%s')} ${sel('rowCheckbox')} ${sel('vCheckboxFunctional')}`,
    tableRowCheckboxInput: `${sel('tableRow-%s')} ${sel('rowCheckbox')} ${sel('vCheckboxFunctional')} input`,
    tableRowColumn: `${sel('tableRow-%s')} ${sel('alarmValue-%s')}`,

    searchInput: `${sel('tableSearch')} ${sel('searchingTextField')}`,
    searchButton: `${sel('tableSearch')} ${sel('submitSearchButton')}`,
    resetSearchButton: `${sel('tableSearch')} ${sel('clearSearchButton')}`,
    helpInformationButton: `${sel('tableSearch')} ${sel('tableSearchHelp')}`,
    helpInformation: sel('tableSearchHelpInfo'),

    topPaginationPrevious: `${sel('topPagination')} ${sel('paginationPreviewButton')}`,
    topPaginationNext: `${sel('topPagination')} ${sel('paginationNextButton')}`,
    topPaginationPage: `${sel('topPagination')} .v-pagination span:nth-of-type(1)`,

    bottomPaginationPrevious: `${sel('vPagination')} li:first-child button`,
    bottomPaginationPage: './/*[@data-test=\'vPagination\']//button[@class=\'v-pagination__item\' and contains(text(), \'%s\')]',
    bottomPaginationNext: `${sel('vPagination')} li:last-child button`,

    itemsPerPage: `${sel('tableWidget')} ${sel('itemsPerPage')} .v-select__slot`,

    mixFilters: `${sel('tableFilterSelector')} div${sel('mixFilters')} .v-input--selection-controls__ripple`,
    mixFiltersInput: `${sel('tableFilterSelector')} input${sel('mixFilters')}`,
    selectFilters: `${sel('tableFilterSelector')} ${sel('selectFilters')} .v-input__slot`,
    clearFilters: `${sel('tableFilterSelector')} ${sel('selectFilters')} .v-input__icon--clear`,
    andFilters: `${sel('tableFilterSelector')} ${sel('andFilters')} + .v-input--selection-controls__ripple`,
    andFiltersInput: `${sel('tableFilterSelector')} input${sel('andFilters')}`,
    orFilters: `${sel('tableFilterSelector')} ${sel('orFilters')} + .v-input--selection-controls__ripple`,
    showFiltersListButton: sel('showFiltersListButton'),

    massActionsPanel: sel('massActionsPanel'),
    massActionsPanelItem: `${sel('tableWidget')} ${sel('massActionsPanel')} ${sel('actionsPanelItem')}:nth-of-type(%s) button`,
    rowActionsSharedPanelItem: `${sel('tableRow-%s')} ${sel('sharedActionsPanel')} .layout ${sel('actionsPanelItem')}:nth-of-type(%s) button`,
    rowMoreActionsButton: `${sel('tableRow-%s')} ${sel('sharedActionsPanel')} .layout ${sel('dropDownActionsButton')}`,
    rowDropDownAction: `.menuable__content__active ${sel('dropDownActions')} ${sel('actionsPanelItem')}:nth-of-type(%s)`,
    rowDropDownActions: `.menuable__content__active ${sel('dropDownActions')} ${sel('actionsPanelItem')}`,

  },
  commands: [commands],
};
