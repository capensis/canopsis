<template lang="pug">
  div
    alarms-list(
    :alarmProperties="$mq | mq(alarmProperties)",
    @openSettings="openSettings"
    )
</template>

<script>
import AlarmsList from '@/components/other/alarm/alarms-list.vue';
import viewMixin from '@/mixins/entities/view';

export default {
  components: {
    AlarmsList,
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
            filter: value => this.$t(`tables.alarmStatus.${value}`),
          },
          {
            text: this.$t('tables.alarmGeneral.lastUpdateDate'),
            value: 'v.last_update_date',
            filter: value => this.$d(new Date(value * 1000), 'long'),
          },
          {
            text: this.$t('tables.alarmGeneral.extraDetails'),
            value: 'extra_details',
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
