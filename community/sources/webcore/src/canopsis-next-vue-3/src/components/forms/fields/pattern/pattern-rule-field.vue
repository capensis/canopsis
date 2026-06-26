<template>
  <c-layout class="gap-3" column>
    <c-layout class="gap-3">
      <div :style="firstColumnFlexStyle">
        <c-layout class="gap-3">
          <div :style="attributeColumnFlexStyle">
            <pattern-attribute-field
              :value="rule.attribute"
              :items="attributes"
              :name="name"
              :disabled="disabled || disabledField"
              return-object
              required
              @update:value="updateAttribute"
            />
          </div>
          <div
            v-if="isAnyInfosRule"
            :style="gridFlexStyle(8)"
          >
            <c-infos-attribute-field
              :value="rule"
              :items="infos"
              :name="name"
              :disabled="disabled"
              :combobox="isInfosRule"
              row
              @update:value="updateModel"
            />
          </div>
          <div
            v-else-if="isObjectRule"
            :style="gridFlexStyle(6)"
          >
            <c-text-field
              :model-value="rule.dictionary"
              :name="objectDictionaryName"
              :disabled="disabled"
              :label="$t('common.dictionary')"
              required
              @update:model-value="updateField('dictionary', $event)"
            />
          </div>
        </c-layout>
      </div>
      <div
        v-if="rule.attribute"
        :style="secondColumnFlexStyle"
      >
        <c-layout class="gap-3">
          <div
            v-if="isInfosValueField"
            :style="gridFlexStyle(1)"
          >
            <c-input-type-field
              :value="rule.fieldType"
              :label="$t('common.type')"
              :types="inputTypesWithDefinedType"
              :disabled="disabled"
              :name="name"
              @update:value="updateField('fieldType', $event)"
            />
          </div>
          <div
            v-if="shownOperatorField"
            :style="operatorFlexStyle"
          >
            <pattern-operator-field
              :value="rule.operator"
              :operators="operators"
              :disabled="disabled || disabledOperator"
              :name="operatorFieldName"
              required
              @update:value="updateField('operator', $event)"
            />
          </div>
          <div
            v-if="rule.operator && operatorHasValue"
            :style="valueFlexStyle"
          >
            <component
              :is="valueComponent.is"
              v-bind="valueComponent.props"
            />
          </div>
        </c-layout>
      </div>
    </c-layout>
    <c-layout>
      <c-alert
        v-if="notDefinedType && !disabled"
        type="warning"
      >
        {{ $t('pattern.notDefinedType') }}
      </c-alert>
    </c-layout>
  </c-layout>
</template>

<script>
import { isFunction } from 'lodash';

import {
  PATTERN_FIELD_TYPES,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  DEFAULT_PATTERN_FIELD_TYPES,
} from '@/constants';

import {
  getFieldType,
  isDateRuleType,
  isDurationRuleType,
  isExtraInfosRuleType,
  isInfosRuleType,
  isObjectRuleType,
  isOperatorHasValue,
  isNumberRuleType,
} from '@/helpers/entities/pattern/form';

import { gridFlexStyle } from '@/helpers/layout';

import { formMixin } from '@/mixins/form';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

import PatternAttributeField from './pattern-attribute-field.vue';
import PatternOperatorField from './pattern-operator-field.vue';
import PatternRuleFieldDateValue from './pattern-rule-field-date-value.vue';

export default {
  emits: ['update:rule'],
  components: {
    DateTimePickerTextField,
    PatternAttributeField,
    PatternOperatorField,
    PatternRuleFieldDateValue,
  },
  mixins: [formMixin],
  props: {
    rule: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    infos: {
      type: Array,
      default: () => [],
    },
    operators: {
      type: Array,
      default: () => [],
    },
    inputTypes: {
      type: Array,
      default: () => DEFAULT_PATTERN_FIELD_TYPES,
    },
    valueField: {
      type: Object,
      required: false,
    },
    type: {
      type: String,
      default: PATTERN_RULE_TYPES.string,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledField: {
      type: Boolean,
      default: false,
    },
    disabledOperator: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'rule',
    },
    itemDisabled: {
      type: Function,
      required: false,
    },
  },
  computed: {
    operatorFieldName() {
      return `${this.name}.operator`;
    },

    valueFieldName() {
      return `${this.name}.value`;
    },

    inputType() {
      return getFieldType(this.rule.value);
    },

    isInfosRule() {
      return isInfosRuleType(this.type);
    },

    isExtraInfosRule() {
      return isExtraInfosRuleType(this.type);
    },

    isObjectRule() {
      return isObjectRuleType(this.type);
    },

    isAnyInfosRule() {
      return this.isInfosRule || this.isExtraInfosRule;
    },

    isAnyInfosRuleOrAlias() {
      return this.isAnyInfosRule || this.rule.alias;
    },

    isInfosValueField() {
      return this.rule.field === PATTERN_RULE_INFOS_FIELDS.value || this.rule.alias;
    },

    isDateRule() {
      return isDateRuleType(this.type) || this.rule.fieldType === PATTERN_FIELD_TYPES.timestamp;
    },

    isDurationRule() {
      return isDurationRuleType(this.type);
    },

    isNumberRule() {
      return isNumberRuleType(this.type);
    },

    inputTypesWithDefinedType() {
      const item = this.isInfosRule
        ? this.infos.find(info => info.value === this.rule.dictionary)
        : this.attributes.find(attribute => attribute.value === this.rule.attribute);

      if (!item?.definedType) {
        return this.inputTypes;
      }

      return this.inputTypes.map(type => ({
        ...type,
        defined: type.value === item.definedType,
      }));
    },

    valueComponent() {
      const valueProps = {
        value: this.rule.value,
        modelValue: this.rule.value,
        required: true,
        disabled: this.disabled,
        name: this.valueFieldName,
        label: this.$t('common.value'),
        'onUpdate:value': this.updateValue,
        'onUpdate:modelValue': this.updateValue,
      };

      if (this.isDateRule) {
        return {
          is: 'pattern-rule-field-date-value',
          props: {
            value: this.rule.range,
            operator: this.rule.operator,
            disabled: this.disabled,
            name: this.valueFieldName,
            'onUpdate:value': this.updateRange,
          },
        };
      }

      if (this.valueField) {
        const valueFieldProps = isFunction(this.valueField.props)
          ? this.valueField.props.call(this, this.rule)
          : this.valueField.props;

        const valueFieldOn = isFunction(this.valueField.on)
          ? this.valueField.on.call(this, this.rule)
          : this.valueField.on;

        return {
          is: this.valueField.is,
          props: {
            disabled: this.disabled,
            ...valueProps,
            ...valueFieldProps,
            ...valueFieldOn,
          },
        };
      }

      if (this.isDurationRule) {
        return {
          is: 'c-duration-field',
          props: {
            duration: this.rule.duration,
            disabled: this.disabled,
            name: this.valueFieldName,
            medium: true,
            'onUpdate:duration': this.updateDuration,
          },
        };
      }

      if (this.isNumberRule) {
        return {
          is: 'c-number-field',
          props: valueProps,
        };
      }

      return {
        is: 'c-mixed-input-field',
        props: {
          inputType: this.inputType,
          types: this.inputTypes,
          ...valueProps,
        },
      };
    },

    operatorHasValue() {
      const hasValue = isOperatorHasValue(this.rule.operator);

      if (this.isAnyInfosRule) {
        return this.isInfosValueField && hasValue;
      }

      return hasValue;
    },

    shownOperatorField() {
      return this.operators.length > 0 || this.operators[0] !== this.rule.operator;
    },

    objectDictionaryName() {
      return `${this.name}.dictionary`;
    },

    firstColumnFlexStyle() {
      return gridFlexStyle(this.isAnyInfosRule ? 5 : 4);
    },

    secondColumnFlexStyle() {
      return gridFlexStyle(this.isAnyInfosRuleOrAlias ? 7 : 8);
    },

    attributeColumnFlexStyle() {
      if (!this.isObjectRule && this.isAnyInfosRule) {
        return gridFlexStyle(4);
      }

      if (this.isObjectRule) {
        return gridFlexStyle(6);
      }

      return {};
    },

    operatorFlexStyle() {
      if (this.isDateRule) {
        return gridFlexStyle(5);
      }

      if (this.isAnyInfosRule) {
        return gridFlexStyle(4);
      }

      return gridFlexStyle(6);
    },

    valueFlexStyle() {
      if (this.isAnyInfosRule || this.isDateRule) {
        return gridFlexStyle(7);
      }

      return gridFlexStyle(6);
    },

    definedType() {
      return this.inputTypesWithDefinedType.find(type => type.defined);
    },

    notDefinedType() {
      return this.rule.fieldType && this.definedType && this.definedType.value !== this.rule.fieldType;
    },
  },
  methods: {
    gridFlexStyle,

    preparedItemDisabled(item) {
      return this.itemDisabled?.(this.rule, item);
    },

    updateAttribute(attribute) {
      const newRule = {
        ...this.rule,

        attribute: attribute?.value,
        alias: attribute?.alias ?? false,
      };

      if (attribute?.definedType) {
        newRule.fieldType = attribute.definedType;
      }

      this.updateModel(newRule);
    },

    updateDuration(duration) {
      this.updateField('duration', duration);
    },

    updateRange(duration) {
      this.updateField('range', duration);
    },

    updateValue(value) {
      this.updateField('value', value);
    },
  },
};
</script>
