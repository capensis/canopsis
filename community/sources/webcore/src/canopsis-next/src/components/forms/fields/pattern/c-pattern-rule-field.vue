<template lang="pug">
  v-layout(row)
    v-flex(:xs5="isInfosRule || isExtraInfosRule", xs4)
      v-layout(row)
        v-flex(:xs4="isInfosRule")
          c-pattern-attribute-field(
            v-field="rule.attribute",
            :items="attributes",
            :name="name",
            :disabled="disabled",
            required
          )
        v-flex.pl-3(v-if="isInfosRule", xs8)
          c-pattern-infos-attribute-field(
            v-field="rule",
            :items="infos",
            :name="name",
            :disabled="disabled"
          )
        v-flex.pl-3(v-else-if="isExtraInfosRule", xs8)
          c-pattern-extra-infos-attribute-field(
            v-field="rule",
            :items="infos",
            :name="name",
            :disabled="disabled"
          )

    v-flex(v-if="rule.attribute", :xs8="!isInfosRule && !isExtraInfosRule", xs7)
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
          v-flex.pl-3(v-if="isInfosValueField || isExtraInfosRule", xs1)
            c-input-type-field(
              :value="inputType",
              :label="$t('common.type')",
              :types="inputTypes",
              :disabled="disabled",
              :name="name",
              @input="updateType"
            )
          v-flex.pl-3(:xs6="!isInfosRule && !isExtraInfosRule", xs4)
            c-pattern-operator-field(
              v-field="rule.operator",
              :operators="operators",
              :disabled="disabled",
              :name="operatorFieldName",
              required
            )

          v-flex.pl-3(v-if="rule.operator && operatorHasValue", :xs7="isInfosRule || isExtraInfosRule", xs6)
            component(
              v-bind="valueComponent.props",
              v-on="valueComponent.on",
              :is="valueComponent.is"
            )
</template>

<script>
import {
  PATTERN_INPUT_TYPES,
  PATTERN_QUICK_RANGES,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  QUICK_RANGES,
} from '@/constants';

import { convertValueByType, getValueType, isOperatorHasValue } from '@/helpers/pattern';

import { formMixin } from '@/mixins/form';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
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
      required: false,
    },
    inputTypes: {
      type: Array,
      default: () => [
        { value: PATTERN_INPUT_TYPES.string },
        { value: PATTERN_INPUT_TYPES.number },
        { value: PATTERN_INPUT_TYPES.boolean },
        { value: PATTERN_INPUT_TYPES.array },
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
      return getValueType(this.rule.value);
    },

    isCustomRange() {
      return this.rule.range.type === QUICK_RANGES.custom.value;
    },

    isInfosRule() {
      return this.type === PATTERN_RULE_TYPES.infos;
    },

    isExtraInfosRule() {
      return this.type === PATTERN_RULE_TYPES.extraInfos;
    },

    isInfosValueField() {
      return this.rule.field === PATTERN_RULE_INFOS_FIELDS.value;
    },

    isDateRule() {
      return this.type === PATTERN_RULE_TYPES.date;
    },

    isDurationRule() {
      return this.type === PATTERN_RULE_TYPES.duration;
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
        return {
          is: this.valueField.is,
          props: {
            ...valueProps,
            ...this.valueField.props,
          },
          on: {
            ...valueHandlers,
            ...this.valueField.on,
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

      if (this.isInfosRule) {
        return this.isInfosValueField && hasValue;
      }

      return hasValue;
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
      this.updateValue(convertValueByType(this.rule.value, type));
    },
  },
};
</script>
