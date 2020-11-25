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
      payload-field(v-field="form.payload", :help-text="$t('modals.createRemediationJob.payloadHelp')")
</template>

<script>
import EnabledField from '@/components/forms/fields/enabled-field.vue';
import PayloadField from '@/components/forms/fields/payload-field.vue';

import RemediationJobConfigurationField from './fields/remediation-job-configuration-field.vue';

export default {
  components: {
    EnabledField,
    PayloadField,
    RemediationJobConfigurationField,
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
