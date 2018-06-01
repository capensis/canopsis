<template lang="pug">
  div
    alarm-list(
    :alarmProperties="$mq | mq(alarmProperties)",
    )
</template>

<script>
import AlarmList from '@/components/alarm-list/alarm-list.vue';

export default {
  components: { AlarmList },
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
