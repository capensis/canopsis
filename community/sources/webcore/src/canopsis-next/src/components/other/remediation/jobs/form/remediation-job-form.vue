<template lang="pug">
  v-layout(column)
    v-layout(row)
      c-enabled-field(v-field="form.enabled")
      c-enabled-field(
        v-field="form.multiple_executions",
        :label="$t('modals.createRemediationJob.fields.multipleExecutions')"
      )
    v-text-field(
      v-field="form.name",
      v-validate="'required'",
      :label="$t('common.name')",
      :error-messages="errors.collect('name')",
      name="name"
    )
    remediation-job-configuration-field(v-field="form.config")
    v-text-field(
      v-field="form.job_id",
      v-validate="'required'",
      :label="$t('modals.createRemediationJob.fields.jobId')",
      :error-messages="errors.collect('job_id')",
      name="job_id"
    )
    v-layout(row)
      v-btn.ml-0(
        v-if="!form.payload",
        color="primary",
        outline,
        @click="addPayload"
      ) {{ $t('modals.createRemediationJob.addPayload') }}
      template(v-else)
        c-json-field(
          v-field="form.payload",
          :label="$t('common.payload')",
          :help-text="$t('modals.createRemediationJob.payloadHelp')",
          name="payload",
          variables
        )
        c-action-btn(
          :tooltip="$t('modals.createRemediationJob.deletePayload')",
          icon="delete",
          color="error",
          bottom,
          @click="removePayload"
        )
    c-text-pairs-field(
      v-field="form.query",
      :title="$t('modals.createRemediationJob.fields.query')",
      :text-label="$t('common.field')",
      :value-label="$t('common.value')",
      name="query"
    )
</template>

<script>
import { formMixin } from '@/mixins/form';

import RemediationJobConfigurationField from './fields/remediation-job-configuration-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationJobConfigurationField,
  },
  mixins: [formMixin],
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
  methods: {
    addPayload() {
      this.updateField('payload', '{}');
    },

    removePayload() {
      this.updateField('payload', '');

      this.$validator.reset({ name: 'payload' });
    },
  },
};
</script>
