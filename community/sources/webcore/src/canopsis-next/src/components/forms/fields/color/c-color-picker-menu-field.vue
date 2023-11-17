<template>
  <v-menu
    :close-on-content-click="false"
    :disabled="disabled"
    content-class="c-color-picker-menu-field__dropdown"
    bottom
    left
    offset-x
  >
    <template #activator="{ on }">
      <v-btn
        class="c-color-picker-menu-field__button ma-0 pa-0"
        v-on="on"
        :style="style"
        :disabled="disabled"
        block
      />
    </template>
    <c-color-chrome-picker-field v-model="colorObject" />
    <c-color-compact-picker-field v-model="colorObject" />
  </v-menu>
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
    min-width: unset !important;
    max-width: 80px;
    width: 80px;
    flex-shrink: 0;
  }

  &__dropdown {
    max-width: 245px;

    .vc-chrome,
    .vc-compact {
      box-shadow: none;
      width: 100%;
    }
  }
}
</style>
