<template lang="pug">
  v-menu(:close-on-content-click="false", :disabled="disabled", bottom, left, offset-x)
    template(#activator="{ on }")
      v-btn.c-color-picker-menu-field__button.ma-0.pa-0(
        v-on="on",
        :style="style",
        :disabled="disabled",
        block
      )
    c-color-chrome-picker-field(v-model="colorObject")
</template>

<script>
import { MODALS } from '@/constants';

import { colorToHex } from '@/helpers/color';

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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    colorObject: {
      get() {
        return colorToHex(this.value);
      },

      set(value) {
        this.updateModel(value.hex);
      },
    },

    style() {
      return {
        backgroundColor: this.value && `${this.value} !important`,
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
.c-color-picker-menu-field {
  &__button {
    min-width: unset;
    max-width: 80px;
    flex-shrink: 0;
  }
}
</style>
