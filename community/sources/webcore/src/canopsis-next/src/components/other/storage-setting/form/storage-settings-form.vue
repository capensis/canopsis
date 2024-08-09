<template>
  <v-layout column>
    <storage-settings-alarm-form
      v-field="form.alarm"
      :history="history.alarm"
    />
    <storage-settings-entity-disabled-form
      :history="history.entity_disabled"
      @archive="$emit('archive:disabled')"
    />
    <storage-settings-entity-unlinked-form
      v-field="form.entity_unlinked"
      :history="history.entity_unlinked"
      @archive="$emit('archive:unlinked')"
    />
    <storage-settings-entity-clean-form
      :history="history.entity_cleaned"
      @clean="$emit('clean:archive')"
    />
    <storage-settings-remediation-form
      v-field="form.remediation"
      :history="history.remediation"
    />
    <storage-settings-pbehavior-form
      v-field="form.pbehavior"
      :history="history.pbehavior"
    />
    <storage-settings-junit-form
      v-field="form.junit"
      :history="history.junit"
    />
    <storage-settings-health-check-form
      v-field="form.health_check"
      :history="history.health_check"
    />
    <storage-settings-webhook-form
      v-field="form.webhook"
      :history="history.webhook"
    />
    <storage-settings-metrics-form v-field="form.metrics" />
    <storage-settings-perf-data-metrics-form v-field="form.perf_data_metrics" />
    <storage-settings-event-filter-failure-form v-field="form.event_filter_failure" />
    <storage-settings-alarm-external-tag-form
      v-field="form.alarm_external_tag"
      :history="history.alarm_external_tag"
    />
    <storage-settings-events-records-form
      v-field="form.event_records"
      :history="history.event_records"
    />
  </v-layout>
</template>

<script>
import StorageSettingsAlarmForm from './storage-settings-alarm-form.vue';
import StorageSettingsEntityDisabledForm from './storage-settings-entity-disabled-form.vue';
import StorageSettingsRemediationForm from './storage-settings-remediation-form.vue';
import StorageSettingsPbehaviorForm from './storage-settings-pbehavior-form.vue';
import StorageSettingsJunitForm from './storage-settings-junit-form.vue';
import StorageSettingsHealthCheckForm from './storage-settings-health-check-form.vue';
import StorageSettingsWebhookForm from './storage-settings-webhook-form.vue';
import StorageSettingsMetricsForm from './storage-settings-metrics-form.vue';
import StorageSettingsPerfDataMetricsForm from './storage-settings-perf-data-metrics-form.vue';
import StorageSettingsEventFilterFailureForm from './storage-settings-event-filter-failure-form.vue';
import StorageSettingsEntityCleanForm from './storage-settings-entity-clean-form.vue';
import StorageSettingsAlarmExternalTagForm from './storage-settings-alarm-external-tag-form.vue';
import StorageSettingsEntityUnlinkedForm from './storage-settings-entity-unlinked-form.vue';
import StorageSettingsEventsRecordsForm from './storage-settings-events-records-form.vue';

export default {
  inject: ['$validator'],
  components: {
    StorageSettingsAlarmExternalTagForm,
    StorageSettingsEntityCleanForm,
    StorageSettingsEntityUnlinkedForm,
    StorageSettingsEntityDisabledForm,
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
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    history: {
      type: Object,
      required: true,
    },
  },
};
</script>
