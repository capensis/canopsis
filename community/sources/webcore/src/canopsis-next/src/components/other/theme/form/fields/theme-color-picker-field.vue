<template lang="pug">
  v-layout.theme-color-picker-field(row, justify-space-between, align-center)
    v-layout(align-center)
      span.v-label.mr-2 {{ label }}
      c-help-icon(v-if="helpText", :text="helpText", top)
    v-btn.theme-color-picker-field__button.ma-0.pa-0(
      :style="style",
      block,
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
    helpText: {
      type: String,
      required: false,
    },
  },
  computed: {
    style() {
      return {
        backgroundColor: this.value,
      };
    },
  },
  methods: {
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

<style lang="scss">
.theme-color-picker-field {
  &__button {
    min-width: unset;
    max-width: 80px;
    flex-shrink: 0;
  }
}
</style>
