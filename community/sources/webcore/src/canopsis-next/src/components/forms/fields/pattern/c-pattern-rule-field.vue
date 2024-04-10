<template lang="pug">
  v-layout(row)
    v-flex(:xs5="isAnyInfosRule", xs4)
      v-layout(row)
        v-flex(:xs4="isAnyInfosRule", :xs6="isObjectRule")
          c-pattern-attribute-field(
            v-field="rule.attribute",
            :items="attributes",
            :name="name",
            :disabled="disabled",
            required
          )
        v-flex.pl-3(v-if="isAnyInfosRule", xs8)
          c-infos-attribute-field(
            v-field="rule",
            :items="infos",
            :name="name",
            :disabled="disabled",
            :combobox="isInfosRule",
            row
          )
        v-flex.pl-3(v-else-if="isObjectRule", xs6)
          v-text-field(
            v-field="rule.dictionary",
            v-validate="'required'",
            :name="objectDictionaryName",
            :disabled="disabled",
            :label="$t('common.dictionary')",
            :error-messages="errors.collect(objectDictionaryName)"
          )
    v-flex(v-if="rule.attribute", :xs8="!isAnyInfosRule", xs7)
      v-layout(row)
        template(v-if="isDateRule")
          v-flex.pl-3(xs5)
            c-quick-date-interval-type-field(
              v-field="rule.range.type",
              :name="name",
              :disabled="disabled",
              :ranges="intervalRanges"
            )
          v-flex.pl-3(v-if="isCustomRange", xs7)
            c-date-time-interval-field(
              v-field="rule.range",
              :name="name",
              :disabled="disabled"
            )

        template(v-else)
          v-flex.pl-3(v-if="isInfosValueField", xs1)
            c-input-type-field(
              :value="rule.fieldType",
              :label="$t('common.type')",
              :types="inputTypes",
              :disabled="disabled",
              :name="name",
              @input="updateType"
            )
          v-flex.pl-3(v-if="shownOperatorField", :xs6="!isAnyInfosRule", xs4)
            c-pattern-operator-field(
              v-field="rule.operator",
              :operators="operators",
              :disabled="disabled",
              :name="operatorFieldName",
              required
            )

          v-flex.pl-3(v-if="rule.operator && operatorHasValue", :xs7="isAnyInfosRule", xs6)
            component(
              v-bind="valueComponent.props",
              v-on="valueComponent.on",
              :is="valueComponent.is"
            )
</template>

<script>
import { isFunction } from 'lodash';

import {
  PATTERN_FIELD_TYPES,
  PATTERN_QUICK_RANGES,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  QUICK_RANGES,
} from '@/constants';

import {
  convertValueByType,
  getFieldType,
  isDateRuleType,
  isDurationRuleType,
  isExtraInfosRuleType,
  isInfosRuleType,
  isObjectRuleType,
  isOperatorHasValue,
} from '@/helpers/pattern';

import { formMixin } from '@/mixins/form';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerTextField },
  mixins: [formMixin],
  model: {
    prop: 'rule',
    event: 'input',
  },
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
      default: () => [
        { value: PATTERN_FIELD_TYPES.string },
        { value: PATTERN_FIELD_TYPES.number },
        { value: PATTERN_FIELD_TYPES.boolean },
        { value: PATTERN_FIELD_TYPES.stringArray },
      ],
    },
    intervalRanges: {
      type: Array,
      default: () => PATTERN_QUICK_RANGES,
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
    name: {
      type: String,
      default: 'rule',
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

    isCustomRange() {
      return this.rule.range.type === QUICK_RANGES.custom.value;
    },

    isInfosRule() {
      return isInfosRuleType(this.type);
    },

    isExtraInfosRule() {
      return isExtraInfosRuleType(this.type);
    },

    isObjectRule() {
      return !!isObjectRuleType(this.type);
    },

    isAnyInfosRule() {
      return this.isInfosRule || this.isExtraInfosRule;
    },

    isInfosValueField() {
      return this.rule.field === PATTERN_RULE_INFOS_FIELDS.value;
    },

    isDateRule() {
      return isDateRuleType(this.type);
    },

    isDurationRule() {
      return isDurationRuleType(this.type);
    },

    valueComponent() {
      const valueProps = {
        value: this.rule.value,
        required: true,
        disabled: this.disabled,
        name: this.valueFieldName,
        label: this.$t('common.value'),
      };

      const valueHandlers = {
        input: this.updateValue,
      };

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
            ...valueProps,
            ...valueFieldProps,
          },
          on: {
            ...valueHandlers,
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
          },
          on: {
            input: this.updateDuration,
          },
        };
      }

      return {
        is: 'c-mixed-input-field',
        props: {
          inputType: this.inputType,
          types: this.inputTypes,
          ...valueProps,
        },
        on: valueHandlers,
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
      return this.operators.length !== 1 || this.operators[0] !== this.rule.operator;
    },

    objectDictionaryName() {
      return `${this.name}.dictionary`;
    },
  },
  methods: {
    updateDuration(duration) {
      this.updateField('duration', duration);
    },

    updateValue(value) {
      this.updateField('value', value);
    },

    updateType(type) {
      this.updateModel({
        ...this.rule,

        fieldType: type,
        value: convertValueByType(this.rule.value, type),
      });
    },
  },
};
</script>
