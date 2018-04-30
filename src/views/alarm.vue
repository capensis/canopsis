<template lang="pug">
  div
    alarm-list(
    :alarmProperty=" $mq | mq({mobile: alarmPropertyOnMobile,tablet: alarmPropertyOnTablet,laptop: alarmProperty})"
    )
</template>

<script>
import Vue from 'vue';
import VueMq from 'vue-mq';
import AlarmList from '@/components/AlarmList/alarm-list.vue';

Vue.use(VueMq, {
  breakpoints: {
    mobile: 450,
    tablet: 900,
    laptop: Infinity,
  },
});

export default {
  name: 'Alarm',
  components: { AlarmList },
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

<style scoped>
</style>
