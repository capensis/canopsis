<template lang="pug">
  v-layout.mb-3(align-center)
    v-flex(xs5)
      c-enabled-field(
        v-field="value.enabled",
        hide-details
      )
    v-flex(xs7)
      v-text-field(
        v-field.number="value.limit",
        v-validate="'required|numeric|min_value:1'",
        :error-messages="errors.collect('enabledLimitValue')",
        :label="label",
        :disabled="!value.enabled",
        :min="1",
        name="enabledLimitValue",
        type="number",
        hide-details
      )
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
      type: Object,
      default: () => ({
        enabled: false,
        limit: 1,
      }),
    },
    label: {
      type: String,
      required: false,
    },
  },
};
</script>
