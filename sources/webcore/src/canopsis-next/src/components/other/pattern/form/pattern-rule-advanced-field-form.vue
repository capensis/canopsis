<template lang="pug">
  v-layout(align-center)
    v-flex(xs3)
      v-select(
        v-field="form.operator",
        :items="availableOperators"
      )
    v-flex.pl-1(xs9)
      mixed-field(v-field="form.value")
    v-flex
      action-btn(
        type="delete",
        @click="$emit('delete')"
      )
</template>

<script>
import { uniq } from 'lodash';

import MixedField from '@/components/forms/fields/mixed-field.vue';
import ActionBtn from '@/components/common/buttons/action-btn.vue';

export default {
  components: { MixedField, ActionBtn },
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
