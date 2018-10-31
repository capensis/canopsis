<template lang="pug">
  v-card.my-1.pa-0
    v-layout(justify-end)
      v-btn(
      @click="$emit('deleteRule')",
      color="red",
      small,
      flat,
      dark,
      fab
      )
        v-icon close
    v-layout(row, wrap, justify-space-around)
      v-flex(xs10, md3)
        v-combobox.my-2(
        :items="possibleFields",
        :value="rule.field",
        @input="updateField('field', $event)",
        solo-inverted,
        hide-details,
        dense,
        flat
        )
      v-flex(xs10, md3)
        v-select.my-2(
        :value="rule.operator",
        :items="operators",
        @change="updateField('operator', $event)",
        solo-inverted,
        hide-details,
        dense,
        flat
        )
      v-flex(xs10, md3)
        v-text-field.my-2(
        v-show="isShownTextField",
        :value="rule.input",
        @input="updateField('input', $event)",
        solo-inverted,
        hide-details,
        single-line,
        flat
        )
</template>

<script>
import formMixin from '@/mixins/form';

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
        return Object.values(this.$constants.FILTER_OPERATORS);
      },
    },
  },
  computed: {
    isShownTextField() {
      return ![
        this.$constants.FILTER_OPERATORS.isEmpty,
        this.$constants.FILTER_OPERATORS.isNotEmpty,
        this.$constants.FILTER_OPERATORS.isNull,
        this.$constants.FILTER_OPERATORS.isNotNull,
      ].includes(this.rule.operator);
    },
  },
};
</script>
