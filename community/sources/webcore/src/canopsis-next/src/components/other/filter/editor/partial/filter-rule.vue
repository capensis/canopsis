<template lang="pug">
  v-card.my-2.pa-0(data-test="filterRule")
    v-layout(justify-end)
      v-btn(
        data-test="deleteRule",
        @click="$emit('delete-rule')",
        color="red",
        small,
        flat,
        dark,
        fab
      )
        v-icon close
    v-layout.px-2(row, wrap, justify-space-around)
      v-flex.pa-1(data-test="fieldRule", xs12, md4)
        filter-rule-field.my-2(
          v-field="rule.field",
          :items="possibleFields",
          :selected-field="selectedField"
        )
      v-flex.pa-1(data-test="operatorRule", xs12, md3)
        v-combobox.my-2(v-field="rule.operator", v-bind="operatorProps")
      v-flex.pa-1(data-test="inputRule", xs12, md5)
        template(v-if="isOperatorForArray")
          v-layout(v-for="(input, index) in rule.input", :key="input.key", row, align-top)
            c-mixed-field.mt-2(
              v-show="isShownInputField",
              v-bind="getValueProps(input.value, input.key)",
              v-validate="valueValidationRules",
              @input="updateField(`input[${index}].value`, $event)"
            )
            v-btn.mt-3(icon, small, @click="removeInput(index)")
              v-icon(color="error", small) close
          v-layout(row, justify-center)
            v-btn(icon, @click="addInput")
              v-icon(color="primary") add
        template(v-else)
          c-mixed-field.my-2(
            v-show="isShownInputField",
            v-bind="getValueProps(rule.input)",
            v-validate="valueValidationRules",
            @input="updateField('input', $event)"
          )
</template>

<script>
import { get, omit } from 'lodash';

import { PATTERN_INPUT_TYPES, FILTER_OPERATORS, FILTER_OPERATORS_FOR_ARRAY } from '@/constants';

import uid from '@/helpers/uid';

import { formMixin } from '@/mixins/form';

import FilterRuleField from './filter-rule-field.vue';

/**
 * Component representing a rule in MongoDB filter
 */
export default {
  inject: ['$validator'],
  components: { FilterRuleField },
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
    possibleFields: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      default: 'rule',
    },
    operators: {
      type: Array,
      default() {
        return Object.values(FILTER_OPERATORS);
      },
    },
  },
  computed: {
    ruleValueName() {
      return `${this.name}.value`;
    },

    selectedField() {
      if (!this.rule.field) {
        return {};
      }

      return this.possibleFields.find(({ value, additionalFieldProps }) => {
        if (additionalFieldProps) {
          return this.rule.field.startsWith(`${value}.`) || value === this.rule.field;
        }

        return value === this.rule.field;
      }) || {};
    },

    operatorProps() {
      const { operatorProps = {} } = this.selectedField;

      return {
        items: this.operators,
        soloInverted: true,
        hideDetails: true,
        dense: true,
        flat: true,

        ...operatorProps,
      };
    },

    valueProps() {
      const { valueProps = {} } = this.selectedField;

      return {
        soloInverted: true,
        dense: true,
        flat: true,
        value: this.rule.input,
        name: this.ruleValueName,
        errorMessages: this.errors.collect(this.ruleValueName),

        ...valueProps,
      };
    },

    valueValidationRules() {
      const { valueValidationRules = {} } = this.selectedField;

      return valueValidationRules;
    },

    isOperatorForArray() {
      return [FILTER_OPERATORS.hasOneOf, FILTER_OPERATORS.hasNot].includes(this.rule.operator);
    },

    isShownInputField() {
      return ![
        FILTER_OPERATORS.isEmpty,
        FILTER_OPERATORS.isNotEmpty,
        FILTER_OPERATORS.isNull,
        FILTER_OPERATORS.isNotNull,
      ].includes(this.rule.operator);
    },
  },
  watch: {
    'rule.operator': {
      handler(value, oldValue) {
        const valueForArray = FILTER_OPERATORS_FOR_ARRAY.includes(value);
        const oldValueForArray = FILTER_OPERATORS_FOR_ARRAY.includes(oldValue);

        if (valueForArray && !oldValueForArray) {
          this.updateField('input', [this.getKeyedInput(this.rule.input)]);
        } else if (!valueForArray && oldValueForArray) {
          this.updateField('input', get(this.rule.input, '0.value', ''));
        }
      },
    },
  },
  methods: {
    getValueProps(value, nameSuffix) {
      const { valueProps = {} } = this.selectedField;
      const name = nameSuffix ? `${this.ruleValueName}.${nameSuffix}` : this.ruleValueName;

      return {
        value,
        name,

        soloInverted: true,
        dense: true,
        flat: true,
        errorMessages: this.errors.collect(name),
        types: valueProps.types
          ? valueProps.types.filter(type => type.value !== PATTERN_INPUT_TYPES.array)
          : [
            { value: PATTERN_INPUT_TYPES.string },
            { value: PATTERN_INPUT_TYPES.number },
            { value: PATTERN_INPUT_TYPES.boolean },
            { value: PATTERN_INPUT_TYPES.null },
          ],

        ...omit(valueProps, ['types']),
      };
    },

    getKeyedInput(value = '') {
      return { value, key: uid() };
    },

    addInput() {
      this.updateField('input', [...this.rule.input, this.getKeyedInput()]);
    },

    removeInput(index) {
      this.updateField('input', this.rule.input.filter((item, i) => i !== index));
    },
  },
};
</script>

<style lang="scss" scoped>
  .input-field {
    border-left: solid 1px #cccccc;
  }

  .type-icon {
    color: inherit;
    opacity: .6;
  }

  .small-avatar {
    min-width: 30px;

    & /deep/ .v-avatar {
      width: 20px !important;
      height: 20px !important;
    }
  }

  .switch-field /deep/ .v-label {
    text-transform: capitalize;
  }
</style>
