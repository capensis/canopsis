<template>
  <c-layout class="flowchart-color-field" justify="space-between" align="center">
    <v-checkbox
      v-if="!hideCheckbox"
      :model-value="isFilled"
      :label="label || $t('common.color')"
      class="mt-0"
      color="primary"
      hide-details
      @update:model-value="updateIsFilled"
    />
    <span v-else class="v-label">
      {{ label }}
    </span>
    <v-col v-show="isFilled || hideCheckbox" class="flowchart-color-field__picker" cols="3">
      <c-color-picker-field
        :model-value="value"
        hide-label
        @update:model-value="updateModel"
      />
    </v-col>
  </c-layout>
</template>

<script>
import { COLORS } from '@/config';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  emits: ['input', 'update:value'],
  props: {
    value: {
      type: String,
      default: 'transparent',
    },
    palette: {
      type: Array,
      default: () => COLORS.flowchart.shapes,
    },
    label: {
      type: String,
      required: false,
    },
    hideCheckbox: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isFilled() {
      return this.value !== 'transparent';
    },

  },
  methods: {
    updateIsFilled(checked) {
      this.updateModel(checked ? this.palette[0] : 'transparent');
    },
  },
};
</script>

<style lang="scss">
.flowchart-color-field {
  &__picker {
    flex: 0 0 auto;
  }
}
</style>
