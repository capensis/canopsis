<template lang="pug">
  v-layout(column)
    c-enabled-field(v-field="form.multiple_executions", :label="$t('remediation.job.multipleExecutions')")
    c-name-field(v-field="form.name", required)
    remediation-job-configuration-field(v-field="form.config")
    c-id-field(
      v-field="form.job_id",
      :label="$t('remediation.job.jobId')",
      name="job_id",
      required
    )
    v-layout(row)
      v-flex.pr-3(xs6)
        c-number-field(v-field="form.retry_amount", :label="$t('remediation.job.retryAmount')")
      v-flex(xs6)
        c-duration-field(
          v-field="form.retry_interval",
          :label="$t('remediation.job.retryInterval')",
          clearable
        )
    v-layout(row)
      v-btn.ml-0(
        v-if="!form.payload",
        color="primary",
        outline,
        @click="addPayload"
      ) {{ $t('remediation.job.addPayload') }}
      template(v-else)
        c-json-field(
          v-field="form.payload",
          :label="$t('common.payload')",
          :help-text="$t('remediation.job.payloadHelp')",
          name="payload",
          variables
        )
        c-action-btn(
          :tooltip="$t('remediation.job.deletePayload')",
          icon="delete",
          color="error",
          bottom,
          @click="removePayload"
        )
    c-text-pairs-field(
      v-field="form.query",
      :title="$t('remediation.job.query')",
      :text-label="$t('common.field')",
      :value-label="$t('common.value')",
      name="query",
      text-required
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
