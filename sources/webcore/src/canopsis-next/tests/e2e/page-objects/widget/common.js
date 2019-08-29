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
    return this.customClick('@periodicRefresh')
      .defaultPause();
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
    return this.customClick('@widgetTitle')
      .defaultPause();
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
    return this.customClick('@rowGridSize')
      .defaultPause();
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
    return this.customClick('@widgetLimit')
      .defaultPause();
  },

  clearWidgetLimitField() {
    return this.customClearValue('@widgetLimitField');
  },

  setWidgetLimitField(limit) {
    return this.customSetValue('@widgetLimitField', limit);
  },

  clickAdvancedSettings() {
    return this.customClick('@advancedSettings')
      .defaultPause();
  },

  clickAlarmList() {
    return this.customClick('@alarmsList')
      .defaultPause();
  },

  clickDefaultSortColumn() {
    return this.customClick('@defaultSortColumn')
      .defaultPause();
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
    return this.customClick('@marginBlock')
      .defaultPause();
  },

  setMargin(position, value) {
    return this.customClick(this.el('@marginHeader', position))
      .setSlider(this.el('@margin', position), value + 1);
  },

  clickHeightFactor() {
    return this.customClick('@widgetHeightFactoryHeader')
      .defaultPause();
  },

  setHeightFactor(value) {
    return this.setSlider(this.el('@widgetHeightFactory'), value);
  },

  clickModalType() {
    return this.customClick('@modalType')
      .defaultPause();
  },

  clickModalTypeField(value = 1) {
    return this.customClick(this.el('@modalTypeField', value))
      .defaultPause();
  },

  clickEditFilter() {
    return this.customClick('@openWidgetFilterEditModal')
      .defaultPause();
  },

  clickDeleteFilter() {
    return this.customClick('@openWidgetFilterDeleteModal')
      .defaultPause();
  },

  clickCreateMoreInfos() {
    return this.customClick('@moreInfoTemplateCreateButton')
      .defaultPause();
  },

  clickEditMoreInfos() {
    return this.customClick('@moreInfoTemplateEditButton')
      .defaultPause();
  },

  clickElementsPerPage() {
    return this.customClick('@elementsPerPage')
      .defaultPause();
  },

  selectElementsPerPage(index = 1) {
    return this.customClick('@elementsPerPageField')
      .waitForElementVisible(this.el('@selectOption', index))
      .customClick(this.el('@selectOption', index));
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

    openWidgetFilterEditModal: sel('openWidgetFilterEditModal'),
    openWidgetFilterDeleteModal: sel('openWidgetFilterDeleteModal'),

    elementsPerPage: sel('elementsPerPage'),
    elementsPerPageField: `${sel('elementsPerPageFieldContainer')} .v-input__slot`,

    moreInfoTemplateCreateButton: `${sel('widgetMoreInfoTemplate')} ${sel('createButton')}`,
    moreInfoTemplateEditButton: `${sel('widgetMoreInfoTemplate')} ${sel('editButton')}`,
    moreInfoTemplateDeleteButton: `${sel('widgetMoreInfoTemplate')} ${sel('deleteButton')}`,
  },
  commands: [commands],
};
