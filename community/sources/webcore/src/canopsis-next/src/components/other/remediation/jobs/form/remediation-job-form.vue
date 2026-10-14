<template>
  <v-row density="comfortable">
    <v-col cols="12">
      <c-name-field v-model="name" autofocus required />
    </v-col>

    <v-col cols="12">
      <c-form-block>
        <c-form-block-row
          :label="t('remediation.job.multipleExecutions')"
          align-center
        >
          <c-enabled-field
            v-model="multipleExecutions"
            :label="t('remediation.job.multipleExecutions')"
            hide-details
            no-margin
          />
        </c-form-block-row>

        <c-form-block-row :label="t('remediation.job.configuration')">
          <remediation-job-configuration-field v-model="config" />
        </c-form-block-row>

        <c-form-block-row :label="t('remediation.job.jobId')">
          <c-id-field
            v-model="jobId"
            :label="t('remediation.job.jobId')"
            name="job_id"
            required
          />
        </c-form-block-row>

        <c-form-block-row :label="t('remediation.job.jobWaitInterval')" align-center>
          <c-enabled-duration-field
            :duration="jobWaitInterval"
            :label="t('remediation.job.jobWaitInterval')"
            :units="jobWaitIntervalUnits"
            name="job_wait_interval"
            switcher
            @update:duration="updateJobWaitInterval"
          />
        </c-form-block-row>

        <c-form-block-row
          v-if="withPayload"
          :label="t('common.payload')"
          :align-center="!payload"
          indented
        >
          <v-btn
            v-if="!payload"
            class="ml-0"
            color="primary"
            variant="outlined"
            @click="addPayload"
          >
            {{ t('remediation.job.addPayload') }}
          </v-btn>
          <v-row v-else>
            <v-col>
              <c-payload-textarea-field
                :value="payload"
                :label="t('common.payload')"
                :variables="templateVars.payload"
                name="payload"
                @update:value="updatePayload"
              >
                <template #append>
                  <c-help-icon :text="t('remediation.job.payloadHelp')" icon="help" left />
                </template>
              </c-payload-textarea-field>
            </v-col>
            <v-col cols="auto">
              <c-action-btn
                :tooltip="t('remediation.job.deletePayload')"
                icon="delete"
                color="error"
                left
                @click="removePayload"
              />
            </v-col>
          </v-row>
        </c-form-block-row>

        <c-form-block-row
          v-if="withQuery"
          :label="t('remediation.job.query')"
        >
          <c-text-pairs-field
            :items="query"
            :text-label="t('common.field')"
            :value-label="t('common.value')"
            :variables="templateVars.payload"
            name="query"
            text-required
            @update:items="updateQuery"
          />
        </c-form-block-row>
      </c-form-block>
    </v-col>
  </v-row>
</template>

<script>
import { AVAILABLE_TIME_UNITS } from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

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
    const { t } = useI18n();
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

    const updateJobWaitInterval = duration => (jobWaitInterval.value = duration);
    const updatePayload = value => (payload.value = value);
    const updateQuery = items => (query.value = items);
    const addPayload = () => (payload.value = '{}');
    const removePayload = () => (payload.value = '');

    return {
      t,
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
