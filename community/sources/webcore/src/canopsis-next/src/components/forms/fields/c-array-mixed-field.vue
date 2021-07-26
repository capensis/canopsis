<template lang="pug">
  div
    v-layout(v-for="(value, index) in values", :key="index")
      v-flex
        c-mixed-field(v-field="values[index]", :disabled="disabled", :types="types")
      c-action-btn(v-if="!disabled", type="delete", @click="removeItemFromArray(index)")
    v-btn.mx-0(color="primary", :disabled="disabled", outline, @click="addItem") {{ $t('common.add') }}
</template>

<script>
import formArrayMixin from '@/mixins/form/array';

export default {
  inject: ['$validator'],
  mixins: [
    formArrayMixin,
  ],
  model: {
    prop: 'values',
    event: 'change',
  },
  props: {
    values: {
      type: Array,
      default: () => [],
    },
    types: {
      type: Array,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    addItem() {
      this.addItemIntoArray('');
    },
  },
};
</script>
