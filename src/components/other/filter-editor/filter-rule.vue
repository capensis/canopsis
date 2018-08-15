<template lang="pug">
  v-card.my-1.pa-0
    v-layout(justify-end)
      v-btn(
      @click="$emit('deleteRule', index)",
      color="red",
      small,
      flat,
      dark,
      fab
      )
        v-icon close
    v-layout(row, wrap, justify-space-around)
      v-flex(xs10, md3)
        v-select.my-2(
        :items="possibleFields",
        :value="field",
        @input="$emit('update:field', $event)",
        solo-inverted,
        hide-details,
        combobox,
        dense,
        flat
        )
      v-flex(xs10, md3)
        v-select.my-2(
        :value="operator",
        :items="operators",
        @change="$emit('update:operator', $event)",
        solo-inverted,
        hide-details,
        dense,
        flat
        )
      v-flex(xs10, md3)
        v-text-field.my-2(
        v-show="isShownTextField",
        :value="input",
        @input="$emit('update:input', $event)",
        solo-inverted,
        hide-details,
        single-line,
        flat
        )
</template>

<script>
import { OPERATORS } from '@/constants';

/**
 * Component representing a rule in MongoDB filter
 *
 * @prop {Number} index - Index of the group
 * @prop {Array} operators - List of all possible operators. Ex : 'equal', 'not equal', 'contains', ...
 * @prop {Array} possibleFields - List of all possible fields to filter on
 * @prop {string} operator - Selected operator
 * @prop {string} field - Selected field
 * @prop {string} input - Input value
 *
 * @event field#update
 * @event operator#update
 * @event input#update
 * @event deleteRule#click
 */
export default {
  props: {
    index: {
      type: Number,
      required: true,
    },
    operators: {
      type: Array,
      required: true,
    },
    possibleFields: {
      type: Array,
      required: true,
    },
    operator: {
      type: String,
      required: true,
    },
    field: {
      type: String,
      required: true,
    },
    input: {
      type: String,
      required: true,
    },
  },
  computed: {
    isShownTextField() {
      return ![
        OPERATORS.isEmpty,
        OPERATORS.isNotEmpty,
        OPERATORS.isNull,
        OPERATORS.isNotNull,
      ].includes(this.operator);
    },
  },
};
</script>
