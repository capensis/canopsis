// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');
const { FILTERS_TYPE } = require('../../constants');

const commands = {
  clickSubmitAlarms() {
    return this.customClick('@submitAlarms');
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

  clickCreateLiveReporting() {
    return this.customClick('@liveReportingCreateButton');
  },

  selectFilter(index = 1) {
    return this.customClick('@selectFilters')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickEditLiveReporting() {
    return this.customClick('@liveReportingEditButton');
  },

  clickDeleteLiveReporting() {
    return this.customClick('@liveReportingDeleteButton');
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

  el,
};

module.exports = {
  elements: {
    submitAlarms: sel('submitAlarms'),

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

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

    liveReportingCreateButton: `${sel('liveReporting')} + div > ${sel('createButton')}`,
    liveReportingEditButton: `${sel('liveReporting')} + div > ${sel('editButton')}`,
    liveReportingDeleteButton: `${sel('liveReporting')} + div > ${sel('deleteButton')}`,

    filters: sel('filters'),
    mixFilters: `div${sel('mixFilters')} .v-input--selection-controls__ripple`,
    mixFiltersInput: `input${sel('mixFilters')}`,
    addFilter: sel('addFilter'),
    andFilters: `${sel('andFilters')} + .v-input--selection-controls__ripple`,
    andFiltersInput: `input${sel('andFilters')}`,
    orFilters: `${sel('orFilters')} + .v-input--selection-controls__ripple`,
    editFilter: sel('editFilter-%s'),
    deleteFilter: sel('deleteFilter-%s'),
    selectFilters: `${sel('selectFilters')} .v-input__slot`,
  },
  commands: [commands],
};
