// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');
const { FILTERS_TYPE } = require('../../../constants');

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createFilterModal');

const commands = {
  setFilterTitle(value) {
    return this.customSetValue('@filterTitle', value);
  },

  clearFilterTitle() {
    return this.customClearValue('@filterTitle');
  },

  setFilterType(type) {
    return this.getAttribute('@radioAndInput', 'aria-checked', ({ value }) => {
      if (value === 'true' && type === FILTERS_TYPE.OR) {
        this.customClick('@radioOr');
      } else if (value === 'false' && type === FILTERS_TYPE.AND) {
        this.customClick('@radioAnd');
      }
    });
  },

  clickAddRule() {
    return this.customClick('@addRule');
  },

  clickAddGroup() {
    return this.customClick('@addGroup');
  },

  clickDeleteGroup() {
    return this.customClick('@deleteGroup');
  },

  selectFieldRule(index = 1) {
    return this.customClick('@fieldRule')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  selectOperatorRule(index = 1) {
    return this.customClick('@operatorRule')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  setInputRule(value) {
    return this.customSetValue('@inputRule', value);
  },

  clearInputRule() {
    return this.customClearValue('@inputRule');
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('createFilterSubmitButton'),
      cancelButton: sel('createFilterCancelButton'),
    }),

    filterTitle: sel('filterTitle'),

    radioAnd: `${sel('filter-group')} input${sel('radioAnd')} + .v-input--selection-controls__ripple`,
    radioAndInput: `${sel('filter-group')} input${sel('radioAnd')}`,
    radioOr: `${sel('filter-group')} input${sel('radioOr')} + .v-input--selection-controls__ripple`,

    addRule: sel('addRule'),

    addGroup: sel('addGroup'),
    deleteGroup: sel('deleteGroup'),

    fieldRule: `${sel('fieldRule')} .v-input__slot`,
    operatorRule: `${sel('operatorRule')} .v-input__slot`,
    inputRule: `${sel('inputRule')} .mixed-field div.v-input:nth-of-type(2) input`,

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
