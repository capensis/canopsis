<template lang="pug">
  div
    settings(
    v-model="isSettingsOpen",
    :title="$t('settings.titles.alarmListSettings')",
    :fields="settingsProperties"
    )
    alarm-list(
    :alarmProperties="$mq | mq(alarmProperties)",
    @openSettings="openSettings"
    )
</template>

<script>
import AlarmList from '@/components/other/alarm-list/alarm-list.vue';
import viewMixin from '@/mixins/view';
import settingsMixin from '@/mixins/settings';

export default {
  components: {
    AlarmList,
  },
  mixins: [
    viewMixin,
    settingsMixin,
  ],
  data() {
    return {
      settingsProperties: [
        'title',
        'default-column-sort',
        'columns',
        'periodic-refresh',
        'default-elements-per-page',
        'opened-resolved-filter',
        'filters',
        'info-popup',
        'more-info',
      ],
      alarmProperties: {
        laptop: [
          {
            text: this.$t('tables.alarmGeneral.connector'),
            value: 'v.connector_name',
          },
          {
            text: this.$t('tables.alarmGeneral.component'),
            value: 'v.component',
          },
          {
            text: this.$t('tables.alarmGeneral.resource'),
            value: 'v.resource',
          },
          {
            text: this.$t('tables.alarmGeneral.output'),
            value: 'v.initial_output',
          },
          {
            text: this.$t('tables.alarmGeneral.state'),
            value: 'v.state.val',
          },
          {
            text: this.$t('tables.alarmGeneral.status'),
            value: 'v.status.val',
          },
          {
            text: this.$t('tables.alarmGeneral.lastUpdateDate'),
            value: 'v.last_update_date',
            filter: value => this.$d(new Date(value * 1000), 'long'),
          },
        ],
        tablet: [
          {
            text: this.$t('tables.alarmGeneral.connector'),
            value: 'v.connector',
          },
          {
            text: this.$t('tables.alarmGeneral.component'),
            value: 'v.component',
          },
          {
            text: this.$t('tables.alarmGeneral.resource'),
            value: 'v.resource',
          },
        ],
        mobile: [
          {
            text: this.$t('tables.alarmGeneral.connector'),
            value: 'v.connector',
          },
          {
            text: this.$t('tables.alarmGeneral.component'),
            value: 'v.component',
          },
        ],
      },
    };
  },
  mounted() {
    this.fetchView({ id: 'view.current_alarms' });
  },
};
</script>
