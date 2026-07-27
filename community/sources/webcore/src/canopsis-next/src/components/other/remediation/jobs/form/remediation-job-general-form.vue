<template>
  <v-layout class="gap-3" column>
    <c-name-field
      v-field="form.name"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row
        :label="$t('remediation.job.multipleExecutions')"
        align-center
      >
        <c-enabled-field
          v-field="form.multiple_executions"
          :label="$t('remediation.job.multipleExecutions')"
          hide-details
          no-margin
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.job.configuration')">
        <remediation-job-configuration-field v-field="form.config" />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.job.jobId')">
        <c-id-field
          v-field="form.job_id"
          :label="$t('remediation.job.jobId')"
          name="job_id"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('remediation.job.jobWaitInterval')">
        <c-enabled-duration-field
          v-field="form.job_wait_interval"
          :label="$t('remediation.job.jobWaitInterval')"
          :units="jobWaitIntervalUnits"
          name="job_wait_interval"
          switcher
        />
      </c-form-block-row>

      <c-form-block-row
        v-if="withPayload"
        :label="$t('common.payload')"
      >
        <v-btn
          v-if="!form.payload"
          class="ml-0"
          color="primary"
          outlined
          @click="addPayload"
        >
          {{ $t('remediation.job.addPayload') }}
        </v-btn>
        <v-layout v-else>
          <c-payload-textarea-field
            v-field="form.payload"
            :label="$t('common.payload')"
            :variables="templateVars.payload"
            name="payload"
          >
            <template #append="">
              <c-help-icon
                :text="$t('remediation.job.payloadHelp')"
                icon="help"
                left
              />
            </template>
          </c-payload-textarea-field>
          <c-action-btn
            :tooltip="$t('remediation.job.deletePayload')"
            icon="delete"
            color="error"
            left
            @click="removePayload"
          />
        </v-layout>
      </c-form-block-row>

      <c-form-block-row
        v-if="withQuery"
        :label="$t('remediation.job.query')"
      >
        <c-text-pairs-field
          v-field="form.query"
          :text-label="$t('common.field')"
          :value-label="$t('common.value')"
          :variables="templateVars.payload"
          name="query"
          text-required
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed, watch } from 'vue';

import { AVAILABLE_TIME_UNITS } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useInfo } from '@/hooks/store/modules/info';
import { useValidator } from '@/hooks/validator/validator';

import RemediationJobConfigurationField from './fields/remediation-job-configuration-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationJobConfigurationField,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const validator = useValidator();
    const { updateField } = useModelField(props, emit);
    const { remediationJobConfigTypes } = useInfo();

    const jobWaitIntervalUnits = [
      AVAILABLE_TIME_UNITS.second,
      AVAILABLE_TIME_UNITS.minute,
      AVAILABLE_TIME_UNITS.hour,
      AVAILABLE_TIME_UNITS.day,
    ];

    const remediationJobConfigType = computed(() => remediationJobConfigTypes.value.find(
      ({ name }) => name === props.form.config?.type,
    ));

    const withPayload = computed(() => remediationJobConfigType.value?.with_body);
    const withQuery = computed(() => remediationJobConfigType.value?.with_query);

    const addPayload = () => updateField('payload', '{}');

    const removePayload = () => {
      updateField('payload', '');

      validator.reset({ name: 'payload' });
    };

    watch(remediationJobConfigType, newConfigType => updateField('configType', newConfigType), { immediate: true });

    return {
      jobWaitIntervalUnits,
      withPayload,
      withQuery,
      addPayload,
      removePayload,
    };
  },
};
</script>
