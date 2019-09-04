// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

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

  togglePeriodicRefreshSwitch(checked = false) {
    return this.getAttribute('@periodicRefreshSwitchInput', 'aria-checked', ({ value }) => {
      if (value === 'false' && checked) {
        this.clickPeriodicRefreshSwitch('@periodicRefreshSwitch');
      }
    });
  },

  clickPeriodicRefreshSwitch() {
    return this.customClick('@periodicRefreshSwitch');
  },

  clearPeriodicRefreshField() {
    return this.customClearValue('@periodicRefreshField');
  },

  setPeriodicRefreshField(value) {
    return this.customSetValue('@periodicRefreshField', value);
  },

  clickWidgetTitle() {
    return this.customClick('@widgetTitle');
  },

  setWidgetTitleField(value) {
    return this.customSetValue('@widgetTitleField', value);
  },

  clearWidgetTitleField() {
    return this.customClearValue('@widgetTitleField');
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
      .waitForElementVisible(this.el('@selectOption', index))
      .customClick(this.el('@selectOption', index));
  },

  selectSortOrders(index = 1) {
    return this.customClick('@defaultSortColumnOrdersField')
      .waitForElementVisible(this.el('@selectOption', index))
      .customClick(this.el('@selectOption', index));
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

  clickEditFilter() {
    return this.customClick('@openWidgetFilterEditModal');
  },

  clickDeleteFilter() {
    return this.customClick('@openWidgetFilterDeleteModal');
  },

  clickCreateMoreInfos() {
    return this.customClick('@moreInfoTemplateCreateButton');
  },

  clickEditMoreInfos() {
    return this.customClick('@moreInfoTemplateEditButton');
  },

  clickElementsPerPage() {
    return this.customClick('@elementsPerPage');
  },

  selectElementsPerPage(index = 1) {
    return this.customClick('@elementsPerPageField')
      .waitForElementVisible(this.el('@selectOption', index))
      .customClick(this.el('@selectOption', index));
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

  toggleColumnNameSwitch(index, checked = false) {
    return this.getAttribute(this.el('@columnNameSwitchFieldInput', index), 'aria-checked', ({ value }) => {
      if (value === 'false' && checked) {
        this.clickColumnNameSwitch(index);
      }
    });
  },

  editColumnName(index = 1, { label, value, isHtml = false }) {
    return this.clickColumnNameLabel(index)
      .clearColumnNameLabel(index, label)
      .setColumnNameLabel(index, label)
      .clickColumnNameValue(index)
      .clearColumnNameValue(index)
      .setColumnNameValue(index, value)
      .toggleColumnNameSwitch(index, isHtml);
  },

  clickFilterOnOpenResolved() {
    return this.customClick('@filterOnOpenResolved')
      .defaultPause();
  },

  toggleOpenFilter(checked = true) {
    return this.getAttribute('@openFilterInput', 'aria-checked', ({ value }) => {
      if (value === 'false' && checked) {
        this.customClick('@openFilter');
      }
    });
  },

  toggleResolvedFilter(checked = false) {
    return this.getAttribute('@resolvedFilterInput', 'aria-checked', ({ value }) => {
      if (value === 'false' && checked) {
        this.customClick('@resolvedFilter');
      }
    });
  },
};

module.exports = {
  elements: {
    selectOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    periodicRefresh: sel('periodicRefresh'),
    periodicRefreshSwitchInput: `input${sel('periodicRefreshSwitch')}`,
    periodicRefreshSwitch: `.v-input${sel('periodicRefreshSwitch')} .v-input--selection-controls__ripple`,
    periodicRefreshField: sel('periodicRefreshField'),

    widgetTitle: sel('widgetTitle'),
    widgetTitleField: sel('widgetTitleField'),
    closeWidget: sel('closeWidget'),

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
  },
  commands: [commands],
};
