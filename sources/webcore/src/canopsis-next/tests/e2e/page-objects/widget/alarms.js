// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickSubmitAlarms() {
    return this.customClick('@submitAlarms');
  },

  clickFilters() {
    return this.customClick('@filters');
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

  setEnableHtml(checked = false) {
    return this.getAttribute('@enableHtmlInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@enableHtml');
      }
    });
  },

  clickAckGroup() {
    return this.customClick('@ackGroup');
  },

  setIsAckNoteRequired(checked = false) {
    return this.getAttribute('@isAckNoteRequiredInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@isAckNoteRequired');
      }
    });
  },

  setIsMultiAckEnabled(checked = false) {
    return this.getAttribute('@isMultiAckEnabledInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@isMultiAckEnabled');
      }
    });
  },

  clickFastAckOutput() {
    return this.customClick('@fastAckOutput');
  },

  setFastAckOutputSwitch(checked = false) {
    return this.getAttribute('@fastAckOutputSwitchInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@fastAckOutputSwitch');
      }
    });
  },

  clickFastAckOutputText() {
    return this.customClick('@fastAckOutputField');
  },

  clearFastAckOutputText() {
    return this.customClearValue('@fastAckOutputField');
  },

  setFastAckOutputText(value) {
    return this.customSetValue('@fastAckOutputField', value);
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
    fastAckOutputSwitch: `div${sel('fastAckOutputSwitch')} .v-input--selection-controls__ripple`,
    fastAckOutputSwitchInput: `input${sel('fastAckOutputSwitch')}`,
    fastAckOutputField: sel('fastAckOutputField'),
  },
  commands: [commands],
};
