<template lang="pug">
  div
    v-layout(row)
      c-enabled-field(v-field="form.enabled")
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row)
      v-text-field(
        v-field="form.description",
        v-validate="'required'",
        :label="$t('common.description')",
        :error-messages="errors.collect('description')",
        name="description"
      )
    v-layout(row)
      remediation-instruction-steps-form(v-field="form.steps")
</template>

<script>
import RemediationInstructionStepsForm from './remediation-instruction-steps-form.vue';

export default {
  inject: ['$validator'],
  components: { RemediationInstructionStepsForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
