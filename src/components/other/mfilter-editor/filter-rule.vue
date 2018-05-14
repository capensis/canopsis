<template lang="pug">
  v-container(fluid, class='filterRule pa-0')
    v-layout(justify-end)
      v-btn(
          @click="handleDeleteRuleClick",
          fab,
          small,
          dark,
          color="red"
        )
          v-icon close
    v-layout(row, wrap, justify-space-around)
      v-flex(xs10, md3)
        v-select(
          @input="$emit('update:field', $event)",
          :items="possibleFields",
          :value="field",
          combobox
        )

      v-flex(xs10, md3)
        v-select(
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
        v-text-field(
          :value="input",
          @input="$emit('update:input', $event)",
        )
</template>

<script>
export default {
  name: 'filter-rule',
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
    border: 1px solid black;
    margin: 0.2em 0;
  }
</style>
