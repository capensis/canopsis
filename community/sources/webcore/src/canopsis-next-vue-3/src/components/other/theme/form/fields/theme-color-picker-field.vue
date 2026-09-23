<template>
  <v-row class="theme-color-picker-field ma-0" align="center" no-gutters>
    <v-col class="theme-color-picker-field__label-wrapper pa-0" cols="auto">
      <label
        :class="{ 'theme-color-picker-field__label--disabled': disabled }"
        class="theme-color-picker-field__label mr-2"
      >
        {{ label }}
      </label>
      <c-help-icon v-if="helpText" :text="helpText" top />
    </v-col>
    <v-col class="theme-color-picker-field__picker pa-0 d-flex align-center">
      <c-color-picker-field
        :model-value="modelValue"
        :label="colorLabel"
        :hide-label="!colorLabel"
        :disabled="disabled"
        justify-end
        hide-details
        @update:model-value="$emit('update:modelValue', $event)"
      />
    </v-col>
  </v-row>
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
    colorLabel: {
      type: String,
      default: '',
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
  &__picker {
    display: flex;
    flex: 0 0 auto;
    align-items: center;
    justify-content: flex-end !important;
    width: auto !important;
    max-width: 100%;
    margin-left: auto;

    > .v-input {
      flex: 0 0 auto;
      width: auto;
    }
  }

  &__label-wrapper {
    flex: 0 0 auto;
    white-space: nowrap;
  }

  &__label {
    font-size: 16px;

    &--disabled {
      opacity: 0.5;
    }
  }
}
</style>
