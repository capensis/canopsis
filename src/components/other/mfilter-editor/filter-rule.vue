<template lang='pug'>
  v-container
    v-form(ref="ruleForm")
      v-layout
        v-flex(xs12)
          v-select(
            @change="$emit('update:field', $event)"
            :items="possibleFields"
            :value="field"
          )

        v-flex(xs12)
          v-select(
            @change="$emit('update:operator', $event)"
            :value="operator"
            :items="operators"
            item-text="value"
          )

        v-flex(
          xs12
          v-if=`
            operator != 'is empty'
            && operator != 'is not empty'
            && operator != 'is null'
            && operator != 'is not null'`
        )
          v-text-field(
            :value="input"
            @input="$emit('update:input', $event)"
          )

        v-flex(xs2)
          v-btn(@click="handleDeleteRuleClick") Delete
</template>

<script>
export default {
  name: 'filter-rule',
  props: {
    index: Number,
    operators: Array,
    possibleFields: Array,
    field: String,
    operator: String,
    input: String,
    isValid: Boolean,
  },
  methods: {
    /**
     * @description Invoked on a click on 'Delete' button. Emit a 'deleteRuleClick' event to the parent,
     * that will actually delete the rule
     */
    handleDeleteRuleClick() {
      this.$emit('deleteRuleClick', this.index);
    },
  },
};
</script>
