<template>
  <c-layout gap="3" column>
    <c-enabled-field
      v-model="multipleExecutions"
      :label="$t('remediation.job.multipleExecutions')"
    />

    <c-name-field v-model="name" autofocus required />

    <remediation-job-configuration-field v-model="config" />

    <c-id-field v-model="jobId" :label="$t('remediation.job.jobId')" name="job_id" required />

    <c-layout>
      <c-enabled-duration-field
        :duration="jobWaitInterval"
        :label="$t('remediation.job.jobWaitInterval')"
        :units="jobWaitIntervalUnits"
        name="job_wait_interval"
        @update:duration="updateJobWaitInterval"
      />
    </c-layout>

    <c-layout v-if="withPayload">
      <v-btn v-if="!payload" class="ml-0" color="primary" variant="outlined" @click="addPayload">
        {{ $t('remediation.job.addPayload') }}
      </v-btn>
      <template v-else>
        <c-payload-textarea-field
          :value="payload"
          :label="$t('common.payload')"
          :variables="templateVars.payload"
          name="payload"
          @update:value="updatePayload"
        >
          <template #append>
            <c-help-icon :text="$t('remediation.job.payloadHelp')" icon="help" left />
          </template>
        </c-payload-textarea-field>
        <c-action-btn
          :tooltip="$t('remediation.job.deletePayload')"
          icon="delete"
          color="error"
          left
          @click="removePayload"
        />
      </template>
    </c-layout>

    <c-text-pairs-field
      v-if="withQuery"
      :items="query"
      :title="$t('remediation.job.query')"
      :text-label="$t('common.field')"
      :value-label="$t('common.value')"
      :variables="templateVars.payload"
      name="query"
      text-required
      @update:items="updateQuery"
    />
  </c-layout>
</template>

<script>
import { AVAILABLE_TIME_UNITS } from '@/constants';

import { useFormModels } from '@/hooks/form';

import RemediationJobConfigurationField from './fields/remediation-job-configuration-field.vue';

export default {
  components: {
    RemediationJobConfigurationField,
  },

  props: {
    withPayload: {
      type: Boolean,
      default: false,
    },
    withQuery: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const {
      multiple_executions: multipleExecutions,
      name,
      config,
      job_id: jobId,
      job_wait_interval: jobWaitInterval,
      payload,
      query,
    } = useFormModels([
      'multiple_executions',
      'name',
      'config',
      'job_id',
      'job_wait_interval',
      'payload',
      'query',
    ]);

    const jobWaitIntervalUnits = [
      AVAILABLE_TIME_UNITS.second,
      AVAILABLE_TIME_UNITS.minute,
      AVAILABLE_TIME_UNITS.hour,
      AVAILABLE_TIME_UNITS.day,
    ];

    const updateJobWaitInterval = (duration) => (jobWaitInterval.value = duration);
    const updatePayload = (value) => (payload.value = value);
    const updateQuery = (items) => (query.value = items);
    const addPayload = () => (payload.value = '{}');
    const removePayload = () => (payload.value = '');

    return {
      multipleExecutions,
      name,
      config,
      jobId,
      jobWaitInterval,
      payload,
      query,
      jobWaitIntervalUnits,
      updateJobWaitInterval,
      updatePayload,
      updateQuery,
      addPayload,
      removePayload,
    };
  },
};
</script>
