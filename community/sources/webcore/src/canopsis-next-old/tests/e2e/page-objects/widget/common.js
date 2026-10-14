// https://nightwatchjs.org/guide/#working-with-page-objects

const { API_ROUTES } = require('@/config');

const el = require('../../helpers/el');
const { elementsWrapperCreator } = require('../../helpers/page-object-creators');
const { FILTERS_TYPE, WAIT_FOR_FIRST_XHR_TIME } = require('../../constants');

const sideBarSelector = sel('sideBarWrapper');

const commands = {
  el,

  setSlider(element, value) {
    return this.dragAndDrop(
      this.el('@sliderThumb', element),
      this.el('@sliderTicks', element, value),
    );
  },

  clickPeriodicRefresh() {
    return this.customClick('@periodicRefresh');
  },

  setPeriodicRefreshSwitch(checked = false) {
    return this.getAttribute('@periodicRefreshSwitchInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.clickPeriodicRefreshSwitch('@periodicRefreshSwitch');
      }
    });
  },

  clickPeriodicRefreshSwitch() {
    return this.customClick('@periodicRefreshSwitch');
  },

  clickPeriodicRefreshField() {
    return this.customClick('@periodicRefreshField');
  },

  clearPeriodicRefreshField() {
    return this.customClearValue('@periodicRefreshField');
  },

  setPeriodicRefreshField(value) {
    return this.customSetValue('@periodicRefreshField', value);
  },

  clickWidgetTitle() {
    return this.customClick('@settingsWidgetTitle');
  },

  setWidgetTitle(value) {
    return this.customSetValue('@widgetTitleField', value);
  },

  clearWidgetTitle() {
    return this.customClearValue('@widgetTitleField');
  },

  getWidgetTitle(id, callback) {
    return this.getText(this.el('@widgetTitle', id), callback);
  },

  clickCloseWidget() {
    return this.customSetValue('@closeWidget');
  },

  clickRowGridSize() {
    return this.customClick('@rowGridSize');
  },

  clearRow() {
    return this.customClearValue('@rowGridSizeCombobox');
  },

  setRow(value) {
    return this.customSetValue('@rowGridSizeCombobox', value)
      .customKeyup('@rowGridSizeCombobox', this.api.Keys.ENTER);
  },

  setRowSize(slider, value) {
    return this.setSlider(
      this.el('@rowSize', slider),
      // The unit is added, because along with 0, the slider has 13 elements.
      value + 1,
    );
  },

  clickWidgetLimit() {
    return this.customClick('@widgetLimit');
  },

  clearWidgetLimitField() {
    return this.customClearValue('@widgetLimitField');
  },

  setWidgetLimitField(limit) {
    return this.customSetValue('@widgetLimitField', limit);
  },

  clickAdvancedSettings() {
    return this.customClick('@advancedSettings');
  },

  clickAlarmList() {
    return this.customClick('@alarmsList');
  },

  clickDefaultSortColumn() {
    return this.customClick('@defaultSortColumn');
  },

  selectSortOrderBy(index = 1) {
    return this.customClick('@defaultSortColumnOrderByField')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  selectSortOrders(index = 1) {
    return this.customClick('@defaultSortColumnOrdersField')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  setColumn(size, value) {
    return this.customClick(this.el('@columnHeader', size))
      .setSlider(this.el('@column', size), value + 1);
  },

  clickMarginBlock() {
    return this.customClick('@marginBlock');
  },

  setMargin(position, value) {
    return this.customClick(this.el('@marginHeader', position))
      .setSlider(this.el('@margin', position), value + 1);
  },

  clickHeightFactor() {
    return this.customClick('@widgetHeightFactoryHeader');
  },

  setHeightFactor(value) {
    return this.setSlider(this.el('@widgetHeightFactory'), value);
  },

  clickModalType() {
    return this.customClick('@modalType');
  },

  clickModalTypeField(value = 1) {
    return this.customClick(this.el('@modalTypeField', value));
  },

  clickCreateFilter() {
    return this.customClick('@openWidgetFilterCreateModal');
  },

  clickEditFilter(name) {
    return this.customClick(this.el('@editFilter', name));
  },

  verifyFilterVisible(name) {
    return this.assert.visible(this.el('@filterItem', name));
  },

  verifyFilterDeleted(name) {
    return this.waitForElementNotPresent(this.el('@filterItem', name));
  },

  clickDeleteFilter(name) {
    return this.customClick(this.el('@deleteFilter', name));
  },

  clickCreateMoreInfos() {
    return this.customClick('@moreInfoTemplateCreateButton');
  },

  clickEditMoreInfos() {
    return this.customClick('@moreInfoTemplateEditButton');
  },

  clickDeleteMoreInfos() {
    return this.customClick('@moreInfoTemplateDeleteButton');
  },

  clickElementsPerPage() {
    return this.customClick('@elementsPerPage');
  },

  selectElementsPerPage(index = 1) {
    return this.customClick('@elementsPerPageField')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickColumnNames() {
    return this.customClick('@columnNames');
  },

  clickAddColumnName() {
    return this.customClick(this.el('@columnNameAddButton'));
  },

  clickDeleteColumnName(index) {
    return this.customClick(this.el('@columnNameDeleteButton', index));
  },

  clickColumnNameUpWard(index = 1) {
    return this.customClick(this.el('@columnNameUpWardButton', index));
  },

  clickColumnNameDownWard(index = 1) {
    return this.customClick(this.el('@columnNameDownWardButton', index));
  },

  clickColumnNameLabel(index = 1) {
    return this.customClick(this.el('@columnNameLabelField', index));
  },

  clearColumnNameLabel(index = 1) {
    return this.customClearValue(this.el('@columnNameLabelField', index));
  },

  setColumnNameLabel(index = 1, value) {
    return this.customSetValue(this.el('@columnNameLabelField', index), value);
  },

  clickColumnNameValue(index = 1) {
    return this.customClick(this.el('@columnNameValueField', index));
  },

  clearColumnNameValue(index = 1) {
    return this.customClearValue(this.el('@columnNameValueField', index));
  },

  setColumnNameValue(index = 1, value) {
    return this.customSetValue(this.el('@columnNameValueField', index), value);
  },

  clickColumnNameSwitch(index = 1) {
    return this.customClick(this.el('@columnNameSwitchField', index));
  },

  setColumnNameIsHtml(index, checked = false) {
    return this.getAttribute(this.el('@columnNameSwitchFieldInput', index), 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.clickColumnNameSwitch(index);
      }
    });
  },

  editColumnName(index = 1, { label, value, isHtml }) {
    this.clickColumnNameLabel(index)
      .clearColumnNameLabel(index, label)
      .setColumnNameLabel(index, label)
      .clickColumnNameValue(index)
      .clearColumnNameValue(index)
      .setColumnNameValue(index, value);

    if (typeof isHtml === 'boolean') {
      this.setColumnNameIsHtml(index, isHtml);
    }

    return this;
  },

  clickFilterOnOpenResolved() {
    return this.customClick('@filterOnOpenResolved');
  },

  setOpenFilter(checked = true) {
    return this.getAttribute('@openFilterInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@openFilter');
      }
    });
  },

  setResolvedFilter(checked = false) {
    return this.getAttribute('@resolvedFilterInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@resolvedFilter');
      }
    });
  },

  clickCreateOrEditInfoPopup() {
    return this.customClick('@widgetInfoPopup');
  },

  clickFilters() {
    return this.customClick('@filters');
  },

  setMixFilters(checked = false) {
    return this.getAttribute('@mixFiltersInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@mixFilters');
      }
    });
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

  clickAddFilter() {
    return this.customClick('@addFilter');
  },

  selectFilter(index = 1) {
    return this.customClick('@selectFilters')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  selectFilterByName(name) {
    this.customClick('@selectFilters')
      .api.useXpath()
      .customClick(this.el('@optionSelectXPath', name))
      .useCss();

    return this;
  },

  clickStatsSelect() {
    return this.customClick('@statsSelector');
  },

  clickAddStat() {
    return this.customClick('@addStatButton');
  },

  clickEditDateInterval() {
    return this.customClick('@editDateInterval');
  },

  clickStatItem(index) {
    return this.customClick(this.el('@statItem', index));
  },

  removeStatItem(index) {
    return this.customClick(this.el('@statItemDeleteButton', index));
  },

  editStatItem(index) {
    return this.customClick(this.el('@statItemEditButton', index));
  },

  clickStatsColor() {
    return this.customClick('@statsColor');
  },

  clickStatsColorItem(title) {
    return this.customClick(this.el('@statsColorPickButton', title));
  },

  clickAnnotationLine() {
    return this.customClick('@widgetStatsAnnotationLine');
  },

  setAnnotationLineEnabled(checked = false) {
    return this.getAttribute('@annotationEnabledInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@annotationEnabled');
      }
    });
  },

  clickAnnotationValue() {
    return this.customClick('@annotationValue');
  },

  clearAnnotationValue() {
    return this.customClearValue('@annotationValue');
  },

  setAnnotationValue(value) {
    return this.customSetValue('@annotationValue', value);
  },

  clickAnnotationLabel() {
    return this.customClick('@annotationLabel');
  },

  clearAnnotationLabel() {
    return this.customClearValue('@annotationLabel');
  },

  setAnnotationLabel(value) {
    return this.customSetValue('@annotationLabel', value);
  },

  clickAnnotationLineColor() {
    return this.customClick('@annotationLineColorButton');
  },

  clickAnnotationLabelColor() {
    return this.customClick('@annotationLabelColorButton');
  },

  clickStatsPointsStyles() {
    return this.customClick('@statsPointsStyles');
  },

  selectStatsPointsStylesType(title, index = 1) {
    return this.customClick(this.el('@statsPointsStyleSelect', title))
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickStatSelectButton() {
    return this.customClick('@statSelectButton');
  },

  assertWidgetRowClasses(id, classes) {
    return this.assert.cssClassPresent(this.el('@widgetRow', id), classes);
  },

  waitFirstUserPreferencesXHR(triggerFunc, callback) {
    return this.waitForFirstXHR(
      API_ROUTES.userPreferences,
      WAIT_FOR_FIRST_XHR_TIME,
      triggerFunc,
      ({ responseData, requestData }) => callback({
        responseData: JSON.parse(responseData),
        requestData: JSON.parse(requestData),
      }),
    );
  },
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    optionSelectXPath: './/*[contains(@class, "menuable__content__active")]//*[contains(@class, "v-select-list")]//span[contains(text(), "%s")]',

    periodicRefresh: sel('periodicRefresh'),
    periodicRefreshSwitchInput: `input${sel('periodicRefreshSwitch')}`,
    periodicRefreshSwitch: `.v-input${sel('periodicRefreshSwitch')} .v-input--selection-controls__ripple`,
    periodicRefreshField: sel('periodicRefreshField'),

    widgetTitle: sel('widgetTitle-%s'),
    settingsWidgetTitle: sel('settingsWidgetTitle'),
    widgetTitleField: sel('widgetTitleField'),
    closeWidget: sel('closeWidget'),

    widgetRow: sel('widgetRow-%s'),
    rowGridSize: sel('rowGridSize'),
    rowGridSizeCombobox: sel('rowGridSizeCombobox'),

    rowSize: `div${sel('slider-%s')}`,

    sliderThumb: '%s .v-slider__thumb',
    sliderTicks: '%s .v-slider__ticks:nth-child(%s)',

    widgetLimit: sel('widgetLimit'),
    widgetLimitField: `${sel('widgetLimit')} .v-text-field__slot input`,

    advancedSettings: sel('advancedSettings'),
    alarmsList: sel('widgetAlarmsList'),

    defaultSortColumn: sel('defaultSortColumn'),
    defaultSortColumnOrderByField: `${sel('defaultSortColumnOrderByLayout')} .v-input__slot`,
    defaultSortColumnOrdersField: `${sel('defaultSortColumnOrdersLayout')} .v-input__slot`,

    columnHeader: sel('column%s'),
    column: sel('column%s'),

    marginBlock: sel('widgetMargin'),
    marginHeader: sel('widget-margin-%s'),
    margin: sel('widget-margin-%s'),

    widgetHeightFactoryHeader: sel('widgetHeightFactory'),
    widgetHeightFactory: sel('widgetHeightFactory'),

    modalType: sel('modalType'),
    modalTypeField: `${sel('modalTypeGroup')} .v-radio:nth-of-type(%s) .v-label`,

    openWidgetFilterCreateModal: `${sel('widgetFilterEditor')} ${sel('createButton')}`,
    openWidgetFilterDeleteModal: `${sel('widgetFilterEditor')} ${sel('deleteButton')}`,
    openWidgetFilterEditModal: `${sel('widgetMoreInfoTemplate')} ${sel('editButton')}`,

    elementsPerPage: sel('elementsPerPage'),
    elementsPerPageField: `${sel('elementsPerPageFieldContainer')} .v-input__slot`,

    moreInfoTemplateCreateButton: `${sel('widgetMoreInfoTemplate')} ${sel('createButton')}`,
    moreInfoTemplateEditButton: `${sel('widgetMoreInfoTemplate')} ${sel('editButton')}`,
    moreInfoTemplateDeleteButton: `${sel('widgetMoreInfoTemplate')} ${sel('deleteButton')}`,

    columnNames: sel('columnNames'),
    columnNameAddButton: sel('columnNameAddButton'),

    columnNameUpWardButton: `${sel('columnName')}:nth-child(%s) ${sel('columnNameUpWard')}`,
    columnNameDownWardButton: `${sel('columnName')}:nth-child(%s) ${sel('columnNameDownWard')}`,
    columnNameLabelField: `${sel('columnName')}:nth-child(%s) ${sel('columnNameLabel')}`,
    columnNameValueField: `${sel('columnName')}:nth-child(%s) ${sel('columnNameValue')}`,
    columnNameSwitchFieldInput: `${sel('columnName')}:nth-child(%s) input${sel('columnNameSwitch')}`,
    columnNameSwitchField: `${sel('columnName')}:nth-child(%s) .v-input${sel('columnNameSwitch')} .v-input--selection-controls__ripple`,
    columnNameDeleteButton: `${sel('columnName')}:nth-child(%s) ${sel('columnNameDeleteButton')}`,

    filterOnOpenResolved: sel('filterOnOpenResolved'),
    openFilter: `div${sel('openFilter')} .v-input--selection-controls__ripple`,
    openFilterInput: `input${sel('openFilter')}`,
    resolvedFilter: `div${sel('resolvedFilter')} .v-input--selection-controls__ripple`,
    resolvedFilterInput: `input${sel('resolvedFilter')}`,

    widgetInfoPopup: sel('infoPopupButton'),

    filters: sel('filters'),

    ...elementsWrapperCreator(sideBarSelector, {
      mixFilters: `div${sel('mixFilters')} .v-input--selection-controls__ripple`,
      mixFiltersInput: `input${sel('mixFilters')}`,
      addFilter: sel('addFilter'),
      andFilters: `${sel('andFilters')} + .v-input--selection-controls__ripple`,
      andFiltersInput: `input${sel('andFilters')}`,
      orFilters: `${sel('orFilters')} + .v-input--selection-controls__ripple`,
      editFilter: sel('editFilter-%s'),
      deleteFilter: sel('deleteFilter-%s'),
      filterItem: sel('filterItem-%s'),

      selectFilters: `${sel('selectFilters')} .v-select__slot`,
    }),

    statsSelector: sel('statsSelector'),
    addStatButton: sel('addStatButton'),

    statSelectButton: `${sel('statSelector')} ${sel('selectButton')}`,

    editDateInterval: `${sel('dateInterval')} ${sel('editButton')}`,

    statItem: `${sel('statItem')}:nth-child(%s)`,
    statItemEditButton: `${sel('statItem')}:nth-child(%s) ${sel('statItemEditButton')}`,
    statItemDeleteButton: `${sel('statItem')}:nth-child(%s)  ${sel('statItemDeleteButton')}`,

    statsColor: sel('statsColor'),
    statsColorPickButton: sel('statsColorPickButton-%s'),

    widgetStatsAnnotationLine: sel('statsAnnotationLine'),
    annotationEnabledInput: `input${sel('annotationEnabled')}`,
    annotationEnabled: `div${sel('annotationEnabled')} .v-input__slot`,
    annotationValue: sel('annotationValue'),
    annotationLabel: sel('annotationLabel'),
    annotationLineColorButton: sel('annotationLineColorButton'),
    annotationLabelColorButton: sel('annotationLabelColorButton'),

    statsPointsStyles: sel('statsPointsStyle'),
    statsPointsStyleSelect: `${sel('statsPointsStyleSelect-%s')} .v-input__slot`,
  },
  commands: [commands],
};
