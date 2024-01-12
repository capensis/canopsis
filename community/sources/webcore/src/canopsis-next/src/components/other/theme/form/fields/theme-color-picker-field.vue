<template>
  <v-layout
    class="theme-color-picker-field"
    justify-space-between
    align-center
  >
    <v-layout align-center>
      <label
        class="theme-color-picker-field__label mr-2"
        :class="{ 'theme-color-picker-field__label--disabled': disabled }"
      >
        {{ label }}
      </label>
      <c-help-icon
        v-if="helpText"
        :text="helpText"
        top
      />
    </v-layout>
    <c-color-picker-menu-field
      v-field="value"
      :disabled="disabled"
    />
  </v-layout>
</template>

<script>
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.theme-color-picker-field {
  &__label {
    font-size: 16px;

    &--disabled {
      opacity: 0.5;
    }
  }
}
</style>
