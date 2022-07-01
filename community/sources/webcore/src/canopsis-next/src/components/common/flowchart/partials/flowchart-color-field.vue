<template lang="pug">
  v-layout(row, justify-between, align-center)
    v-checkbox.mt-0(
      :input-value="isFilled",
      :label="label || $t('flowchart.color')",
      hide-details,
      @change="updateIsFilled"
    )
    v-btn.ma-0(
      v-show="isFilled",
      :style="style",
      small,
      @click="showColorPickerModal"
    )
</template>

<script>
import { MODALS } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: 'transparent',
    },
    label: {
      type: String,
      required: false,
    },
  },
  computed: {
    isFilled() {
      return this.value !== 'transparent';
    },

    style() {
      return {
        backgroundColor: this.value,
      };
    },
  },
  methods: {
    updateIsFilled(checked) {
      this.updateModel(checked ? 'white' : 'transparent');
    },

    showColorPickerModal() {
      this.$modals.show({
        name: MODALS.colorPicker,
        config: {
          color: this.value,
          action: color => this.updateModel(color),
        },
      });
    },
  },
};
</script>
