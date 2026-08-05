<template>
  <c-layout
    :class="{ 'theme-enabled-color-picker-field--disabled': !modelValue.enabled }"
    class="theme-enabled-color-picker-field"
    column
  >
    <c-enabled-field
      :model-value="modelValue.enabled"
      :label="enableLabel"
      class="mt-0"
      @update:model-value="updateField('enabled', $event)"
    >
      <template #append>
        <c-help-icon v-if="enableHelpText" :text="enableHelpText" top />
      </template>
    </c-enabled-field>
    <v-expand-transition>
      <theme-color-picker-field
        v-if="modelValue.enabled"
        v-bind="$attrs"
        :model-value="modelValue.color"
        :disabled="!modelValue.enabled"
        @update:model-value="updateField('color', $event)"
      />
    </v-expand-transition>
  </c-layout>
</template>

<script>
import ThemeColorPickerField from './theme-color-picker-field.vue';

export default {
  components: { ThemeColorPickerField },
  inheritAttrs: false,
  emits: ['update:modelValue'],
  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    enableHelpText: {
      type: String,
      required: false,
    },
    enableLabel: {
      type: String,
      required: false,
    },
  },
  methods: {
    updateField(field, value) {
      this.$emit('update:modelValue', {
        ...this.modelValue,
        [field]: value,
      });
    },
  },
};
</script>

<style lang="scss">
.theme-enabled-color-picker-field {
  &--disabled .v-label {
    opacity: 0.5;
  }
}
</style>
