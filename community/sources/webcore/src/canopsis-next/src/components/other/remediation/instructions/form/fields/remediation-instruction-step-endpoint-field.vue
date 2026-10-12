<template>
  <v-text-field
    v-field="value"
    v-validate="'required'"
    :label="$t('remediation.instruction.endpoint')"
    :name="name"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
  >
    <template #append="">
      <c-help-icon
        :text="$t('remediation.instruction.tooltips.endpoint')"
        icon="help"
        left
      />
    </template>
  </v-text-field>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
};
</script>
