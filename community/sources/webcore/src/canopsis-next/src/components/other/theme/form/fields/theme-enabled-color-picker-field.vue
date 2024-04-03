<template>
  <v-layout
    :class="{ 'theme-enabled-color-picker-field--disabled': !value.enabled }"
    class="theme-enabled-color-picker-field"
    column
  >
    <c-enabled-field
      v-field="value.enabled"
      :label="enableLabel"
      class="mt-0"
    >
      <template #append="">
        <c-help-icon
          v-if="enableHelpText"
          :text="enableHelpText"
          top
        />
      </template>
    </c-enabled-field>
    <theme-color-picker-field
      v-field="value.color"
      v-bind="$attrs"
      :disabled="!value.enabled"
    />
  </v-layout>
</template>

<script>
import ThemeColorPickerField from './theme-color-picker-field.vue';

export default {
  components: { ThemeColorPickerField },
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
};
</script>

<style lang="scss">
.theme-enabled-color-picker-field {
  .theme--light.v-label {
    color: unset;
  }

  &--disabled .v-label {
    opacity: 0.5;
  }
}
</style>
