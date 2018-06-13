<template lang="pug">
  div(v-if="show")
    alarms-list-settings
    alarm-list(
    :alarmProperties="$mq | mq(alarmProperties)",
    )
</template>

<script>
import AlarmList from '@/components/alarm-list/alarm-list.vue';
import AlarmsListSettings from '@/components/alarms-list-settings.vue';
import viewMixin from '@/mixins/view';

export default {
  components: {
    AlarmList,
    AlarmsListSettings,
  },
  mixins: [
    viewMixin,
  ],
  data() {
    return {
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
            filter: value => this.$d(new Date(value * 1000), 'short'),
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
      show: false,
    };
  },
  watch: {
    viewPending() {
      if (!this.viewPending) {
        this.show = true;
      }
    },
  },
  mounted() {
    this.fetchView({ id: 'view.current_alarms' });
  },
};
</script>
