<template lang="pug">
  v-layout(row)
    v-flex(:xs5="isInfosRule", xs4)
      v-layout(row)
        v-flex(:xs4="isInfosRule")
          c-pattern-attribute-field(
            v-field="rule.attribute",
            :items="attributes",
            :name="name",
            :disabled="disabled"
          )
        v-flex.pl-3(v-if="isInfosRule", xs8)
          c-pattern-infos-attribute-field(
            v-field="rule",
            :items="infos",
            :name="name",
            :disabled="disabled"
          )

    template(v-if="isDateRule")
      v-flex.pl-3(xs3)
        c-quick-date-interval-type-field(
          v-field="rule.range.type",
          :name="name",
          :disabled="disabled"
        )
      v-flex.pl-3(xs5)
        c-date-time-interval-field(
          v-if="isCustomRange",
          v-field="rule.range",
          :name="name",
          :disabled="disabled"
        )

    template(v-else-if="isInfosRule")
      v-flex.pl-3(v-if="isInfosValueField", xs1)
        c-input-type-field(
          :value="inputType",
          :label="$t('common.type')",
          :types="inputTypes",
          :disabled="disabled",
          :name="name",
          @input="updateType"
        )
      v-flex.pl-3(xs2)
        c-pattern-operator-field(
          v-field="rule.operator",
          :operators="operators",
          :disabled="disabled",
          :name="operatorFieldName"
        )
      v-flex.pl-3(v-if="isInfosValueField", xs4)
        c-mixed-input-field(
          v-field="rule.value",
          :input-type="inputType",
          :types="inputTypes",
          :label="$t('common.value')",
          :name="valueFieldName"
        )

    template(v-else)
      v-flex.pl-3(xs4)
        c-pattern-operator-field(
          v-field="rule.operator",
          :operators="operators",
          :disabled="disabled",
          :name="operatorFieldName"
        )

      v-flex.pl-3(xs4)
        template(v-if="operatorHasValue")
          c-duration-field(
            v-if="isDurationRule",
            v-field="rule.duration",
            :name="valueFieldName"
          )
          component(
            v-else,
            v-field="rule.value",
            v-bind="valueComponent.props",
            v-on="valueComponent.on",
            :is="valueComponent.is",
            :name="valueFieldName"
          )
</template>

<script>
import { PATTERN_INPUT_TYPES, PATTERN_RULE_INFOS_FIELDS, PATTERN_RULE_TYPES, QUICK_RANGES } from '@/constants';

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

    isInfosValueField() {
      return this.rule.field === PATTERN_RULE_INFOS_FIELDS.value;
    },

    isDateRule() {
      return this.type === PATTERN_RULE_TYPES.date;
    },

    isNumberRule() {
      return this.type === PATTERN_RULE_TYPES.number;
    },

    isDurationRule() {
      return this.type === PATTERN_RULE_TYPES.duration;
    },

    valueComponent() {
      if (this.valueField) {
        return this.valueField;
      }

      if (this.isInfosRule && this.isInfosValueField) {
        return {
          is: 'c-mixed-field',
          props: {
            class: 'mt-1',
          },
        };
      }

      if (this.isNumberRule) {
        return {
          is: 'c-number-field',
        };
      }

      return {
        is: 'v-text-field',
        props: {
          disabled: this.disabled,
          label: this.$t('common.value'),
        },
      };
    },

    operatorHasValue() {
      return isOperatorHasValue(this.rule.operator);
    },
  },
  methods: {
    updateType(type) {
      this.updateField('value', convertValueByType(this.rule.value, type));
    },

    updateInterval(interval) {
      this.updateModel({
        ...this.rule,
        ...interval,
      });
    },
  },
};
</script>
