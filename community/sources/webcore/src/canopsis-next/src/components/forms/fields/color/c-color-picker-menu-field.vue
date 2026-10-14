<template>
  <v-menu
    v-model="shown"
    :close-on-content-click="false"
    :disabled="disabled"
    content-class="c-color-picker-menu-field__dropdown"
    location="bottom end"
  >
    <template #activator="{ props: activatorProps }">
      <v-btn
        v-bind="activatorProps"
        :style="style"
        :disabled="disabled"
        class="c-color-picker-menu-field__button ma-0 pa-0"
        block
      />
    </template>

    <v-card>
      <v-color-picker v-model="colorValue" mode="hex" hide-inputs />
    </v-card>
  </v-menu>
</template>

<script>
import { computed, ref } from 'vue';

import { colorToHex } from '@/helpers/color';

export default {
  props: {
    modelValue: {
      type: String,
      default: 'transparent',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const shown = ref(false);

    const colorValue = computed({
      get: () => colorToHex(props.modelValue),
      set: value => emit('update:modelValue', value),
    });

    const style = computed(() => ({
      backgroundColor: props.modelValue && `${props.modelValue} !important`,
    }));

    return {
      shown,
      colorValue,
      style,
    };
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

    .v-color-picker {
      box-shadow: none;
      width: 100%;
    }
  }
}
</style>
