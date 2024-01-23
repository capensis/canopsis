<template>
  <v-layout column>
    <c-enabled-field
      v-field="form.multiple_executions"
      :label="$t('remediation.job.multipleExecutions')"
    />
    <c-name-field
      v-field="form.name"
      required
    />
    <remediation-job-configuration-field v-field="form.config" />
    <c-id-field
      v-field="form.job_id"
      :label="$t('remediation.job.jobId')"
      name="job_id"
      required
    />
    <v-layout>
      <c-enabled-duration-field
        v-field="form.job_wait_interval"
        :label="$t('remediation.job.jobWaitInterval')"
        :units="jobWaitIntervalUnits"
        name="job_wait_interval"
      />
    </v-layout>
    <v-layout v-if="withPayload">
      <v-btn
        v-if="!form.payload"
        class="ml-0"
        color="primary"
        outlined
        @click="addPayload"
      >
        {{ $t('remediation.job.addPayload') }}
      </v-btn>
      <template v-else>
        <c-json-field
          v-field="form.payload"
          :label="$t('common.payload')"
          :help-text="$t('remediation.job.payloadHelp')"
          name="payload"
          variables
        />
        <c-action-btn
          :tooltip="$t('remediation.job.deletePayload')"
          icon="delete"
          color="error"
          bottom
          @click="removePayload"
        />
      </template>
    </v-layout>
    <c-text-pairs-field
      v-if="withQuery"
      v-field="form.query"
      :title="$t('remediation.job.query')"
      :text-label="$t('common.field')"
      :value-label="$t('common.value')"
      name="query"
      text-required
    />
  </v-layout>
</template>

<script>
import { AVAILABLE_TIME_UNITS } from '@/constants';

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
    withPayload: {
      type: Boolean,
      default: false,
    },
    withQuery: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    jobWaitIntervalUnits() {
      return [
        AVAILABLE_TIME_UNITS.second,
        AVAILABLE_TIME_UNITS.minute,
        AVAILABLE_TIME_UNITS.hour,
        AVAILABLE_TIME_UNITS.day,
      ];
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
