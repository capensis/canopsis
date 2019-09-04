// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitAlarms() {
    return this.customClick('@submitAlarms');
  },

  clickFilters() {
    return this.customClick('@filters')
      .defaultPause();
  },

  clickAddFilter() {
    return this.customClick('@addFilter');
  },

  clickMixFilters() {
    return this.customClick('@mixFilters');
  },

  clickAndFilters() {
    return this.customClick('@andFilters');
  },

  clickOrFilters() {
    return this.customClick('@orFilters');
  },

  selectFilters(index = 1) {
    return this.customClick('@selectFilters')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickInfoPopupButton() {
    return this.customClick('@infoPopupButton');
  },

  toggleEnableHtml(checked = false) {
    return this.getAttribute('@enableHtmlInput', 'aria-checked', ({ value }) => {
      if (value === 'false' && checked) {
        this.customClick('@enableHtml');
      }
    });
  },

  clickAckGroup() {
    return this.customClick('@ackGroup')
      .defaultPause();
  },

  clickIsAckNoteRequired() {
    return this.customClick('@isAckNoteRequired');
  },

  clickIsMultiAckEnabled() {
    return this.customClick('@isMultiAckEnabled');
  },

  clickFastAckOutput() {
    return this.customClick('@fastAckOutput')
      .defaultPause();
  },

  clickFastAckOutputSwitch() {
    return this.customClick('@fastAckOutputSwitch');
  },

  setFastAckOutputText(value) {
    return this.customSetValue('@fastAckOutputText', value);
  },

  el,
};

module.exports = {
  elements: {
    submitAlarms: sel('submitAlarms'),

    filters: sel('filters'),
    addFilter: sel('addFilter'),
    editFilter: sel('editFilter-%s'),
    deleteFilter: sel('deleteFilter-%s'),
    mixFilters: `div${sel('mixFilters')} .v-input--selection-controls__ripple`,
    andFilters: `${sel('andFilters')} + .v-input--selection-controls__ripple`,
    orFilters: `${sel('orFilters')} + .v-input--selection-controls__ripple`,
    selectFilters: `${sel('selectFilters')} .v-input__slot`,

    infoPopupButton: sel('infoPopupButton'),

    enableHtml: `${sel('isHtmlEnabledOnTimeLine')} div${sel('switcherField')}`,
    enableHtmlInput: `${sel('isHtmlEnabledOnTimeLine')} input${sel('switcherField')}`,

    isAckNoteRequired: `${sel('isAckNoteRequired')} div${sel('switcherField')}`,
    isAckNoteRequiredInput: `${sel('isAckNoteRequired')} input${sel('switcherField')}`,
    isMultiAckEnabled: `${sel('isMultiAckEnabled')} div${sel('switcherField')}`,
    isMultiAckEnabledInput: `${sel('isMultiAckEnabled')} input${sel('switcherField')}`,
    ackGroup: sel('ackGroup'),
    fastAckOutput: sel('fastAckOutput'),
    fastAckOutputSwitch: `${sel('fastAckOutput')} .v-input--switch .v-input--selection-controls__ripple`,
    fastAckOutputText: `${sel('fastAckOutput')} .v-text-field input`,
  },
  commands: [commands],
};
