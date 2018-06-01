<template lang="pug">
  div
    alarms-list-settings
    alarm-list(
    :alarmProperty="$mq | mq({mobile: alarmPropertyOnMobile,tablet: alarmPropertyOnTablet,laptop: alarmProperty})",
    )
</template>

<script>
import AlarmList from '@/components/alarm-list/alarm-list.vue';
import AlarmsListSettings from '@/components/alarms-list-settings.vue';

export default {
  name: 'Alarm',

  components: {
    AlarmList,
    AlarmsListSettings,
  },
  data() {
    return {
      alarmProperty: [
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
          text: this.$t('tables.alarmGeneral.lastUpdateDate'),
          value: 'v.last_update_date',
        },
      ],
      alarmPropertyOnMobile: [
        {
          text: this.$t('tables.alarmGeneral.connector'),
          value: 'v.connector',
        },
        {
          text: this.$t('tables.alarmGeneral.component'),
          value: 'v.component',
        },
      ],
      alarmPropertyOnTablet: [
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
    };
  },

};
</script>
