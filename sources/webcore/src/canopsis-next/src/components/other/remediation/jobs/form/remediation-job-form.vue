<template lang="pug">
  v-layout(column)
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
      remediation-job-configuration-field(v-field="form.config")
    v-layout(row)
      v-text-field(
        v-field="form.job_id",
        v-validate="'required'",
        :label="$t('modals.createRemediationJob.fields.jobId')",
        :error-messages="errors.collect('job_id')",
        name="job_id"
      )
    v-layout(row)
      remediation-job-payload-field(v-field="form.payload")
</template>

<script>
import EnabledField from '@/components/forms/fields/enabled-field.vue';

import RemediationJobConfigurationField from './fields/remediation-job-configuration-field.vue';
import RemediationJobPayloadField from './fields/remediation-job-payload-field.vue';

export default {
  components: {
    EnabledField,
    RemediationJobConfigurationField,
    RemediationJobPayloadField,
  },
  inject: ['$validator'],
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
