// https://nightwatchjs.org/guide/#working-with-page-objects

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const commands = {
  setFilterTitle(value) {
    return this.customSetValue('@filterTitle', value);
  },
  clearFilterTitle() {
    return this.customClearValue('@filterTitle');
  },
  clickSubmitFilter() {
    return this.customClick('@submitFilter');
  },
  clickRadioAnd() {
    return this.customClick('@radioAnd');
  },
  clickRadioOr() {
    return this.customClick('@radioOr');
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
};

const modalSelector = sel('createFilterModal');

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitFilter: sel('submitFilter'),
      filterTitle: sel('filterTitle'),
      radioAnd: `${sel('filter-group')} input${sel('radioAnd')} + .v-input--selection-controls__ripple`,
      radioOr: `${sel('filter-group')} input${sel('radioOr')} + .v-input--selection-controls__ripple`,
      addRule: sel('addRule'),
      addGroup: sel('addGroup'),
      deleteGroup: sel('deleteGroup'),
      fieldRule: `${sel('fieldRule')} .v-input__slot`,
      operatorRule: `${sel('operatorRule')} .v-input__slot`,
      optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
      inputRule: `${sel('inputRule')} .mixed-field div.v-input:nth-of-type(2) input`,
    }),
  },
  commands: [commands],
});
