<template lang="pug">
  v-container(fluid, class='filterRule pa-0')
    v-layout(justify-end)
      v-btn(
          @click="handleDeleteRuleClick",
          fab,
          small,
          flat,
          dark,
          color="red"
        )
          v-icon close
    v-layout(row, wrap, justify-space-around)
      v-flex(xs10, md3)
        v-select.my-2(
          solo-inverted,
          flat,
          hide-details,
          dense,
          @input="$emit('update:field', $event)",
          :items="possibleFields",
          :value="field",
          combobox
        )
      v-flex(xs10, md3)
        v-select.my-2(
          solo-inverted,
          hide-details,
          flat,
          dense,
          @change="$emit('update:operator', $event)",
          :value="operator",
          :items="operators",
          item-text="value"
        )
      v-flex(
        xs10,
        md3,
        v-if=`
          operator != 'is empty'
          && operator != 'is not empty'
          && operator != 'is null'
          && operator != 'is not null'`
      )
        v-text-field.my-2(
          solo-inverted,
          flat,
          hide-details,
          single-line,
          :value="input",
          @input="$emit('update:input', $event)"
        )
</template>

<script>
/**
 * Component representing a rule in MongoDB filter
 *
 * @prop {Number} [index] - Index of the group
 * @prop {Array} [operators] - List of all possible operators. Ex : 'equal', 'not equal', 'contains', ...
 * @prop {Array} [possibleFields] - List of all possible fields to filter on
 * @prop {String} [field] - Selected field
 * @prop {String} [operator] - Selected operator
 * @prop {String} [input] - Input value
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
    field: {
      type: String,
      required: true,
    },
    operator: {
      type: String,
      required: true,
    },
    input: {
      type: String,
      required: true,
    },
  },
  methods: {
    /**
     * Invoked on a click on 'Delete' button. Emit a 'deleteRuleClick' event to the parent,
     * that will actually delete the rule
     */
    handleDeleteRuleClick() {
      this.$emit('deleteRuleClick', this.index);
    },
  },
};
</script>

<style scoped>
  .filterRule {
    border: 1px solid lightgray;
    margin: 0.2em 0;
  }
</style>
