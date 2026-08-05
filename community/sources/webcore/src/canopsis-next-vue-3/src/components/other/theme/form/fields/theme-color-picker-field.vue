<template>
  <c-layout class="theme-color-picker-field" justify="space-between" align="center">
    <c-layout align="center">
      <label
        :class="{ 'theme-color-picker-field__label--disabled': disabled }"
        class="theme-color-picker-field__label mr-2"
      >
        {{ label }}
      </label>
      <c-help-icon v-if="helpText" :text="helpText" top />
    </c-layout>
    <c-color-picker-menu-field
      :model-value="modelValue"
      :disabled="disabled"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </c-layout>
</template>

<script>
export default {
  emits: ['update:modelValue'],
  props: {
    modelValue: {
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
