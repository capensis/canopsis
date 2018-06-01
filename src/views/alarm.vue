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
      alarmProperty: {
        Connector: 'v.connector',
        Component: 'v.component',
        Resource: 'v.resource',
        Output: 'v.initial_output',
        'Last Update Date': 'v.last_update_date',
      },
      alarmPropertyOnMobile: {
        Connector: 'v.connector',
        Component: 'v.component',
      },
      alarmPropertyOnTablet: {
        Connector: 'v.connector',
        Component: 'v.component',
        Resource: 'v.resource',
      },
    };
  },

};
</script>
