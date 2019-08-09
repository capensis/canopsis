// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  el,

  clickPeriodicRefresh() {
    return this.customClick('@periodicRefresh')
      .defaultPause();
  },

  getPeriodicRefreshSwitchStatus() {
    let status = false;
    this.getAttribute('@periodicRefreshSwitchInput', 'aria-checked', ({ value }) => {
      status = value === 'true';
    });
    return status;
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
      .customKeyup('@rowGridSizeCombobox', 'ENTER');
  },
  setSlider(slider, value) {
    return this.dragAndDrop(
      this.el('@sliderThumb', slider),
      this.el('@sliderTicks', slider, value),
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
  },
  commands: [commands],
};
