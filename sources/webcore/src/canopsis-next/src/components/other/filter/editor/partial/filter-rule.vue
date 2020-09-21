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
          v-field="rule.field",
          :items="possibleFields",
          solo-inverted,
          hide-details,
          dense,
          flat,
          :return-object="false",
          item-text="name",
          item-value="value",
          :loading="filterHintsPending"
        )
          template(slot="item", slot-scope="props")
            v-list-tile-content {{ props.item.name }} ({{ props.item.value }})
      v-flex.pa-1(data-test="operatorRule", xs12, md3)
        v-combobox.my-2(
          v-field="rule.operator",
          :items="operators",
          solo-inverted,
          hide-details,
          dense,
          flat
        )
      v-flex.pa-1(data-test="inputRule", xs12, md5)
        template(v-if="isOperatorForArray")
          v-layout(v-for="(input, index) in rule.input", :key="input.key", row, align-center)
            mixed-field.my-2(
              v-field="rule.input[index].value",
              v-show="isShownInputField",
              solo-inverted,
              hide-details,
              flat
            )
            v-btn(icon, small, @click="removeInput(index)")
              v-icon(color="error", small) close
          v-layout.mt-2(row, justify-center)
            v-btn(icon, @click="addInput")
              v-icon(color="primary") add
        template(v-else)
          mixed-field.my-2(
            v-field="rule.input",
            v-show="isShownInputField",
            solo-inverted,
            hide-details,
            flat
          )
</template>

<script>
import { isBoolean, isNumber, get } from 'lodash';

import { FILTER_OPERATORS, FILTER_OPERATORS_FOR_ARRAY, FILTER_INPUT_TYPES } from '@/constants';

import uid from '@/helpers/uid';

import formMixin from '@/mixins/form';
import filterHintsMixin from '@/mixins/entities/filter-hint';

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
  mixins: [formMixin, filterHintsMixin],
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
    isOperatorForArray() {
      return [FILTER_OPERATORS.in, FILTER_OPERATORS.notIn].includes(this.rule.operator);
    },

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
