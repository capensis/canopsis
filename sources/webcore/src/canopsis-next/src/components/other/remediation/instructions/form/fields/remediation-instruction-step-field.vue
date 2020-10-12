<template lang="pug">
  v-layout(row)
    v-flex(xs11)
      v-layout
        v-text-field(
          v-field="step.name",
          v-validate="'required'",
          :label="$t('common.name')",
          :error-messages="nameErrorMessages",
          :name="name",
          box
        )
      v-layout
        remediation-instruction-steps-workflow-field(v-field="step.stop_on_fail")
    v-flex.mt-3(xs1)
      v-layout(justify-center)
        v-btn.ma-0(icon, small, @click.prevent="$emit('remove')")
          v-icon(color="error") delete
</template>

<script>
import formMixin from '@/mixins/form';

import RemediationInstructionStepsWorkflowField from './remediation-instruction-steps-workflow-field.vue';

export default {
  inject: ['$validator'],
  components: { RemediationInstructionStepsWorkflowField },
  mixins: [formMixin],
  model: {
    prop: 'step',
    event: 'input',
  },
  props: {
    step: {
      type: Object,
      required: true,
    },
  },
  computed: {
    fieldSuffix() {
      return this.step.key ? `-${this.step.key}` : '';
    },

    name() {
      return `name${this.fieldSuffix}`;
    },

    nameErrorMessages() {
      return this.errors.collect(this.name).map(error => error.replace(this.fieldSuffix, ''));
    },
  },
};
</script>
