<template lang="pug">
  v-layout(align-center)
    v-flex(xs3)
      v-select(
        v-field="form.operator",
        :items="availableOperators"
      )
    v-flex.pl-1(xs9)
      c-mixed-field(v-field="form.value")
    v-flex
      c-action-btn(
        type="delete",
        @click="$emit('delete')"
      )
</template>

<script>
import { uniq } from 'lodash';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    operators: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    availableOperators() {
      return uniq([...this.operators, this.form.operator]);
    },
  },
};
</script>
