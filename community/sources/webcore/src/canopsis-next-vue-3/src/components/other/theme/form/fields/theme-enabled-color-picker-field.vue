<template>
  <v-row
    :class="{ 'theme-enabled-color-picker-field--disabled': !modelValue.enabled }"
    class="theme-enabled-color-picker-field"
    align="center"
    no-gutters
  >
    <v-col class="pa-0" cols="auto">
      <c-enabled-field
        :model-value="modelValue.enabled"
        :label="enableLabel"
        class="py-2"
        hide-details
        @update:model-value="updateField('enabled', $event)"
      >
        <template #append>
          <c-help-icon v-if="enableHelpText" :text="enableHelpText" top />
        </template>
      </c-enabled-field>
    </v-col>
    <v-col class="theme-enabled-color-picker-field__color pa-0 d-flex align-center">
      <slot name="append" />
      <theme-color-picker-field
        v-if="modelValue.enabled"
        v-bind="$attrs"
        :model-value="modelValue.color"
        :disabled="!modelValue.enabled"
        @update:model-value="updateField('color', $event)"
      >
        <template v-if="$slots.label" #label>
          <slot name="label" />
        </template>
      </theme-color-picker-field>
    </v-col>
  </v-row>
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
  &__color {
    display: flex;
    flex: 0 0 auto;
    align-items: center;
    justify-content: flex-end;
    width: auto;
    max-width: 100%;
    margin-left: auto;
    gap: 8px;
  }

  .c-enabled-field .v-input--switch {
    margin-top: 0;
    padding-top: 0;
  }

  .c-enabled-field,
  .c-enabled-field .v-label {
    flex: 0 0 auto;
    white-space: nowrap;
  }

  &--disabled .v-label {
    opacity: 0.5;
  }
}
</style>
