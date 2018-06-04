<template lang="pug">
  div
    alarms-list-settings
    alarm-list(
    :alarmProperties="$mq | mq(alarmProperties)",
    )
</template>

<script>
import AlarmList from '@/components/alarm-list/alarm-list.vue';
import AlarmsListSettings from '@/components/alarms-list-settings.vue';

export default {
  components: {
    AlarmList,
    AlarmsListSettings,
  },
  data() {
    return {
      alarmProperties: {
        laptop: {
          Connector: 'v.connector',
          Component: 'v.component',
          Resource: 'v.resource',
          Output: 'v.initial_output',
          'Last Update Date': 'v.last_update_date',
        },
        tablet: {
          Connector: 'v.connector',
          Component: 'v.component',
          Resource: 'v.resource',
        },
        mobile: {
          Connector: 'v.connector',
          Component: 'v.component',
        },
      },
    };
  },
};
</script>
