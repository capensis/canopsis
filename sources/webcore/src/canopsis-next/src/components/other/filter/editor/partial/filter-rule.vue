<template lang="pug">
  v-card.my-2.pa-0(data-test="filterRule")
    v-layout(justify-end)
      v-btn(
      data-test="deleteRule",
      @click="$emit('deleteRule')",
      color="red",
      small,
      flat,
      dark,
      fab
      )
        v-icon close
    v-layout.px-2(row, wrap, justify-space-around)
      v-flex.pa-1(data-test="fieldRule", xs12, md4)
        v-combobox.my-2(
        :items="possibleFields",
        :value="rule.field",
        @input="updateField('field', $event)",
        solo-inverted,
        hide-details,
        dense,
        flat
        )
      v-flex.pa-1(data-test="operatorRule", xs12, md4)
        v-combobox.my-2(
        :value="rule.operator",
        :items="operators",
        @input="updateField('operator', $event)",
        solo-inverted,
        hide-details,
        dense,
        flat
        )
      v-flex.pa-1(data-test="inputRule", xs12, md4)
        mixed-field.my-2(
        v-show="isShownInputField",
        :value="rule.input",
        solo-inverted,
        hide-details,
        flat,
        @input="updateField('input', $event)"
        )
</template>

<script>
import { isBoolean, isNumber } from 'lodash';

import { FILTER_OPERATORS, FILTER_INPUT_TYPES } from '@/constants';

import formMixin from '@/mixins/form';

import MixedField from '@/components/forms/fields/mixed-field.vue';

/**
 * Component representing a rule in MongoDB filter
 *
 * @prop {Object} rule - Object of the rule
 * @prop {Array} possibleFields - List of all possible fields to filter on
 * @prop {Array} [operators=Object.values(FILTER_OPERATORS)] - List of all possible operators. Ex : 'equal', ...
 *
 * @event field#update
 * @event operator#update
 * @event input#update
 * @event deleteRule#click
 */
export default {
  components: { MixedField },
  mixins: [formMixin],
  model: {
    prop: 'rule',
    event: 'update:rule',
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
    operators: {
      type: Array,
      default() {
        return Object.values(FILTER_OPERATORS);
      },
    },
  },
  data() {
    return {
      types: [
        { text: 'String', value: FILTER_INPUT_TYPES.string },
        { text: 'Number', value: FILTER_INPUT_TYPES.number },
        { text: 'Boolean', value: FILTER_INPUT_TYPES.boolean },
      ],
    };
  },
  computed: {
    switchLabel() {
      return String(this.rule.input);
    },

    inputType() {
      if (isBoolean(this.rule.input)) {
        return FILTER_INPUT_TYPES.boolean;
      } else if (isNumber(this.rule.input)) {
        return FILTER_INPUT_TYPES.number;
      }

      return FILTER_INPUT_TYPES.string;
    },

    isInputTypeText() {
      return [FILTER_INPUT_TYPES.number, FILTER_INPUT_TYPES.string].includes(this.inputType);
    },

    getInputTypeIcon() {
      const TYPES_ICONS_MAP = {
        [FILTER_INPUT_TYPES.string]: 'title',
        [FILTER_INPUT_TYPES.number]: 'exposure_plus_1',
        [FILTER_INPUT_TYPES.boolean]: 'toggle_on',
      };

      return type => TYPES_ICONS_MAP[type];
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
  methods: {
    updateInputField(value) {
      const isInputTypeNumber = this.inputType === FILTER_INPUT_TYPES.number;

      this.updateField('input', isInputTypeNumber ? Number(value) : value);
    },
    updateInputTypeField(value) {
      switch (value) {
        case FILTER_INPUT_TYPES.number:
          this.updateField('input', Number(this.rule.input));
          break;
        case FILTER_INPUT_TYPES.boolean:
          this.updateField('input', Boolean(this.rule.input));
          break;
        case FILTER_INPUT_TYPES.string:
          this.updateField('input', String(this.rule.input));
          break;
      }
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
      width: 20px!important;
      height: 20px!important;
    }
  }

  .switch-field /deep/ .v-label {
    text-transform: capitalize;
  }
</style>
