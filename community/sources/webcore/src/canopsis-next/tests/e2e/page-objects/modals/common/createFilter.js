// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../../helpers/el');
const { FILTERS_TYPE, VALUE_TYPES } = require('../../../constants');

const { elementsWrapperCreator, modalCreator } = require('../../../helpers/page-object-creators');

const modalSelector = sel('createFilterModal');

const commands = {
  selectGroup([firstIndex, ...path]) {
    const childrens = path.map(item => this.el('@filterGroupLayout', item));
    const parent = this.el('@filterGroup', firstIndex);

    return [parent, ...childrens].join(' ');
  },

  clearFilterTitle() {
    return this.customClearValue('@filterTitle');
  },

  setFilterTitle(value) {
    return this.customSetValue('@filterTitle', value);
  },

  setFilterType(groupSelector, type) {
    return this.getAttribute(
      this.el('@radioAndInput', groupSelector),
      'aria-checked',
      ({ value }) => {
        if (value === 'true' && type === FILTERS_TYPE.OR) {
          this.customClick(this.el('@radioOr', groupSelector));
        } else if (value === 'false' && type === FILTERS_TYPE.AND) {
          this.customClick(this.el('@radioAnd', groupSelector));
        }
      },
    );
  },

  clickAddRule(groupSelector) {
    return this.customClick(this.el('@addRule', groupSelector));
  },

  clickDeleteRule(groupSelector, index) {
    return this.customClick(this.el('@deleteRule', groupSelector, index));
  },

  clickAddGroup(groupSelector) {
    return this.customClick(this.el('@addGroup', groupSelector));
  },

  clickDeleteGroup(groupSelector) {
    return this.customClick(this.el('@deleteGroup', groupSelector));
  },

  selectFieldRule(groupSelector, ruleIndex, index = 1) {
    return this.customClick(this.el('@fieldRule', groupSelector, ruleIndex))
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  selectOperatorRule(groupSelector, ruleIndex, index = 1) {
    return this.customClick(this.el('@operatorRule', groupSelector, ruleIndex))
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clearInputRule(groupSelector, ruleIndex) {
    return this.customClearValue(this.el('@inputRule', groupSelector, ruleIndex));
  },

  setInputRule(groupSelector, ruleIndex, value) {
    return this.customSetValue(this.el('@inputRule', groupSelector, ruleIndex), value);
  },

  setInputRuleSwitch(groupSelector, ruleIndex, checked) {
    return this.getAttribute(
      this.el('@mixedInputSwitchField', groupSelector, ruleIndex),
      'aria-checked',
      ({ value }) => {
        if (value !== String(checked)) {
          this.customClick(this.el('@mixedInputSwitch', groupSelector, ruleIndex), checked);
        }
      },
    );
  },

  selectValueType(groupSelector, ruleIndex, index = 1) {
    return this.customClick(this.el('@inputType', groupSelector, ruleIndex))
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  fillRuleFields(path, rulesIndex, {
    rule, value, valueType, operator, groups,
  }) {
    const groupSelector = this.selectGroup(path);

    this.selectFieldRule(groupSelector, rulesIndex, rule)
      .selectOperatorRule(groupSelector, rulesIndex, operator)
      .selectValueType(groupSelector, rulesIndex, valueType);

    switch (valueType) {
      case VALUE_TYPES.NUMBER:
      case VALUE_TYPES.STRING:
        this.clearInputRule(groupSelector, rulesIndex)
          .setInputRule(groupSelector, rulesIndex, value);
        break;
      case VALUE_TYPES.BOOLEAN:
        this.setInputRuleSwitch(groupSelector, rulesIndex, value);
        break;
      default:
    }

    if (groups) {
      this.clickAddGroup(groupSelector)
        .fillFilterGroups(groups, path);
    }

    return this;
  },

  fillFilterGroup(path, group) {
    const selector = this.selectGroup(path);

    this.setFilterType(selector, group.type);

    group.items.forEach((item, index) => {
      this.clickAddRule(selector)
        .fillRuleFields(path, index + 1, item);
    });

    return this;
  },

  fillFilterGroups(groups = [], rootPath = []) {
    groups.forEach((group, index) => {
      const path = [...rootPath, index + 1];

      if (index !== 0) {
        this.clickAddGroup(this.selectGroup([...rootPath, index]));
      }

      this.fillFilterGroup(path, group);
    });

    return this;
  },

  el,
};

module.exports = modalCreator(modalSelector, {
  elements: {
    ...elementsWrapperCreator(modalSelector, {
      submitButton: sel('createFilterSubmitButton'),
      cancelButton: sel('createFilterCancelButton'),
    }),

    filterGroup: `${sel('filterGroup')}:nth-of-type(%s)`,
    filterGroupLayout: `${sel('filterGroupLayout')} ${sel('filterGroup')}:nth-of-type(%s)`,

    filterTitle: sel('filterTitle'),

    radioAnd: `%s input${sel('radioAnd')} + .v-input--selection-controls__ripple`,
    radioAndInput: `%s input${sel('radioAnd')}`,
    radioOr: `%s input${sel('radioOr')} + .v-input--selection-controls__ripple`,

    addRule: `%s ${sel('addRule')}`,
    deleteRule: `%s ${sel('filterRule')}:nth-of-type(%s) ${sel('deleteRule')}`,

    addGroup: `%s ${sel('addGroup')}`,
    deleteGroup: `%s ${sel('deleteGroup')}`,

    fieldRule: `%s ${sel('filterRule')}:nth-of-type(%s) ${sel('fieldRule')} .v-input__slot`,
    operatorRule: `%s ${sel('filterRule')}:nth-of-type(%s) ${sel('operatorRule')} .v-input__slot`,
    inputRule: `%s ${sel('filterRule')}:nth-of-type(%s) input${sel('mixedInput')}`,
    inputType: `%s ${sel('filterRule')}:nth-of-type(%s) .mixed-field__type-selector`,

    mixedInputSwitchField: `%s ${sel('filterRule')}:nth-of-type(%s) input${sel('mixedInputSwitch')}`,
    mixedInputSwitch: `%s ${sel('filterRule')}:nth-of-type(%s) div${sel('mixedInputSwitch')} .v-input--selection-controls__ripple`,

    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',
  },
  commands: [commands],
});
