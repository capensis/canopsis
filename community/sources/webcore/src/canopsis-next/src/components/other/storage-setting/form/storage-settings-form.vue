<template>
  <c-layout column>
    <storage-settings-alarm-form v-model:form="alarm" :history="history.alarm" />
    <storage-settings-entity-unlinked-form
      v-model:form="entity"
      :history="history.entity_unlinked"
      @archive="$emit('archive:unlinked')"
    />
    <storage-settings-remediation-form v-model:form="remediation" :history="history.remediation" />
    <storage-settings-pbehavior-form v-model:form="pbehavior" :history="history.pbehavior" />
    <storage-settings-junit-form v-model:form="junit" :history="history.junit" />
    <storage-settings-health-check-form
      v-model:form="healthCheck"
      :history="history.health_check"
    />
    <storage-settings-webhook-form v-model:form="webhook" :history="history.webhook" />
    <storage-settings-metrics-form v-model:form="metrics" />
    <storage-settings-perf-data-metrics-form v-model:form="perfDataMetrics" />
    <storage-settings-event-filter-failure-form
      v-model:form="eventFilterFailure"
      :history="history.event_filter_failure"
    />
    <storage-settings-entity-infos-log-form v-model:form="entityInfosLog" />
    <storage-settings-alarm-external-tag-form
      v-model:form="alarmExternalTag"
      :history="history.alarm_external_tag"
    />
    <storage-settings-events-records-form
      v-model:form="eventRecords"
      :history="history.event_records"
    />
    <storage-settings-llm-chat-form v-model:form="llmChat" />
    <storage-settings-connector-anomalies-form v-model:form="connectorAnomalies" />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import StorageSettingsAlarmExternalTagForm from './storage-settings-alarm-external-tag-form.vue';
import StorageSettingsAlarmForm from './storage-settings-alarm-form.vue';
import StorageSettingsConnectorAnomaliesForm from './storage-settings-connector-anomalies-form.vue';
import StorageSettingsEntityInfosLogForm from './storage-settings-entity-infos-log-form.vue';
import StorageSettingsEntityUnlinkedForm from './storage-settings-entity-unlinked-form.vue';
import StorageSettingsEventFilterFailureForm from './storage-settings-event-filter-failure-form.vue';
import StorageSettingsEventsRecordsForm from './storage-settings-events-records-form.vue';
import StorageSettingsHealthCheckForm from './storage-settings-health-check-form.vue';
import StorageSettingsJunitForm from './storage-settings-junit-form.vue';
import StorageSettingsLlmChatForm from './storage-settings-llm-chat-form.vue';
import StorageSettingsMetricsForm from './storage-settings-metrics-form.vue';
import StorageSettingsPbehaviorForm from './storage-settings-pbehavior-form.vue';
import StorageSettingsPerfDataMetricsForm from './storage-settings-perf-data-metrics-form.vue';
import StorageSettingsRemediationForm from './storage-settings-remediation-form.vue';
import StorageSettingsWebhookForm from './storage-settings-webhook-form.vue';

export default {
  components: {
    StorageSettingsAlarmExternalTagForm,
    StorageSettingsEntityUnlinkedForm,
    StorageSettingsEventFilterFailureForm,
    StorageSettingsPerfDataMetricsForm,
    StorageSettingsMetricsForm,
    StorageSettingsWebhookForm,
    StorageSettingsHealthCheckForm,
    StorageSettingsJunitForm,
    StorageSettingsPbehaviorForm,
    StorageSettingsRemediationForm,
    StorageSettingsAlarmForm,
    StorageSettingsEventsRecordsForm,
    StorageSettingsEntityInfosLogForm,
    StorageSettingsLlmChatForm,
    StorageSettingsConnectorAnomaliesForm,
  },

  inheritAttrs: false,

  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    history: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:modelValue', 'archive:disabled', 'archive:unlinked', 'clean:archive'],

  setup(props, { emit }) {
    const section = path => computed({
      get: () => props.modelValue?.[path],
      set: value => emit('update:modelValue', {
        ...props.modelValue,
        [path]: value,
      }),
    });

    return {
      alarm: section('alarm'),
      entity: section('entity'),
      remediation: section('remediation'),
      pbehavior: section('pbehavior'),
      junit: section('junit'),
      healthCheck: section('health_check'),
      webhook: section('webhook'),
      metrics: section('metrics'),
      perfDataMetrics: section('perf_data_metrics'),
      eventFilterFailure: section('event_filter_failure'),
      entityInfosLog: section('entity_infos_log'),
      alarmExternalTag: section('alarm_external_tag'),
      eventRecords: section('event_records'),
      llmChat: section('llm_chat'),
      connectorAnomalies: section('connector_anomalies'),
    };
  },
};
</script>
