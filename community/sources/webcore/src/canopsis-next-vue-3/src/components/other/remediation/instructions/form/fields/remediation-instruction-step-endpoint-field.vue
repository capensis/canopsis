<template>
  <v-text-field
    v-validate="'required'"
    :model-value="modelValue"
    :label="$t('remediation.instruction.endpoint')"
    :name="name"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #append="">
      <c-help-icon :text="$t('remediation.instruction.tooltips.endpoint')" icon="help" left />
    </template>
  </v-text-field>
</template>

<script>
export default {
  inject: ['$validator'],
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'endpoint',
    },
  },
  emits: ['update:modelValue'],
};
</script>
