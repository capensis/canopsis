<template lang="pug">
  div
    v-layout(v-for="(value, index) in values", :key="index", align-center)
      v-flex
        v-text-field(
          v-field="values[index]",
          :disabled="disabled",
          :label="$t('common.value')"
        )
      c-action-btn(v-if="!disabled", type="delete", @click="removeItemFromArray(index)")
    v-messages(:value="errorMessages", color="error")
    v-btn.mx-0(:disabled="disabled", color="primary", outline, @click="addItem") {{ $t('common.add') }}
</template>

<script>
import { formArrayMixin } from '@/mixins/form/array';

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
    errorMessages: {
      type: Array,
      default: () => [],
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
