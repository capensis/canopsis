// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  el,

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

  setSlider(slider, value) {
    return this.dragAndDrop(
      this.el('@sliderThumb', slider),
      // The unit is added, because along with 0, the slider has 13 elements.
      this.el('@sliderTicks', slider, value + 1),
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

  clickDefaultSortColumn() {
    return this.customClick('@defaultSortColumn')
      .defaultPause();
  },

  selectSortOrderBy(index = 1) {
    return this.customClick('@defaultSortColumnOrderByField')
      .waitForElementVisible(this.el('@defaultSortColumnOrderByOption', index))
      .customClick(this.el('@defaultSortColumnOrderByOption', index));
  },

  selectSortOrders(index = 1) {
    return this.customClick('@defaultSortColumnOrdersField')
      .waitForElementVisible(this.el('@defaultSortColumnOrdersOption', index))
      .customClick(this.el('@defaultSortColumnOrdersOption', index));
  },

  setColumn(size, value) {
    return this.customClick(this.el('@column', size))
      .dragAndDrop(
        this.el('@columnThumb', size),
        this.el('@columnTicks', size, value + 1),
      );
  },

  clickMarginBlock() {
    return this.customClick('@marginBlock')
      .defaultPause();
  },

  setMargin(position, value) {
    return this.customClick(this.el('@margin', position))
      .dragAndDrop(
        this.el('@marginThumb', position),
        this.el('@marginTicks', position, value + 1),
      );
  },

  clickHeightFactor() {
    return this.customClick('@widgetHeightFactory')
      .defaultPause();
  },

  setHeightFactor(value) {
    return this.dragAndDrop(
      this.el('@heightFactoryThumb'),
      this.el('@heightFactoryTicks', value),
    );
  },
};

module.exports = {
  elements: {
    periodicRefresh: sel('periodicRefresh'),
    periodicRefreshSwitchInput: `input${sel('periodicRefreshSwitch')}`,
    periodicRefreshSwitch: `.v-input${sel('periodicRefreshSwitch')} .v-input--selection-controls__ripple`,
    periodicRefreshField: sel('periodicRefreshField'),

    widgetTitle: sel('widgetTitle'),
    widgetTitleField: sel('widgetTitleField'),
    closeWidget: sel('closeWidget'),

    rowGridSize: sel('rowGridSize'),
    rowGridSizeCombobox: sel('rowGridSizeCombobox'),

    sliderThumb: `div${sel('slider-%s')} .v-slider__thumb`,
    sliderTicks: `div${sel('slider-%s')} .v-slider__ticks:nth-child(%s)`,

    widgetLimit: sel('widgetLimit'),
    widgetLimitField: `${sel('widgetLimit')} .v-text-field__slot input`,

    advancedSettings: sel('advancedSettings'),

    defaultSortColumn: sel('defaultSortColumn'),
    defaultSortColumnOrderByField: `${sel('defaultSortColumnOrderByLayout')} .v-input__slot`,
    defaultSortColumnOrderByOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
    defaultSortColumnOrdersField: `${sel('defaultSortColumnOrdersLayout')} .v-input__slot`,
    defaultSortColumnOrdersOption: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    column: `${sel('column%s')} .v-list__group__header`,
    columnThumb: `${sel('column%s')} .v-slider__thumb`,
    columnTicks: `${sel('column%s')} .v-slider__ticks:nth-child(%s)`,

    marginBlock: `${sel('widgetMarginBlock')} .v-list__group__header`,
    margin: `${sel('widget-margin-%s')} .v-list__group__header`,
    marginThumb: `${sel('widget-margin-%s')} .v-slider__thumb`,
    marginTicks: `${sel('widget-margin-%s')} .v-slider__ticks:nth-child(%s)`,

    widgetHeightFactory: `${sel('widgetHeightFactory')} .v-list__group__header`,
    heightFactoryThumb: `${sel('widgetHeightFactory')} .v-slider__thumb`,
    heightFactoryTicks: `${sel('widgetHeightFactory')} .v-slider__ticks:nth-child(%s)`,
  },
  commands: [commands],
};
