<template lang="pug">
  div
    v-layout(row)
      enabled-field(v-field="form.enabled")
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
import entitiesRemediationJobsMixin from '@/mixins/entities/remediation/jobs';

import EnabledField from '@/components/forms/fields/enabled-field.vue';
import RemediationInstructionStepsForm from './remediation-instruction-steps-form.vue';

export default {
  components: { RemediationInstructionStepsForm, EnabledField },
  inject: ['$validator'],
  mixins: [entitiesRemediationJobsMixin],
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
  mounted() {
    this.fetchRemediationJobsList();
  },
};
</script>
