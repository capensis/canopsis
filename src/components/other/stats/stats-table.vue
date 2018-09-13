<template lang="pug">
  v-container(fluid)
    div Stats table
</template>

<script>
export default {
  data() {
    return {
      /*
      stats: [
        {
          text: 'Alarms created',
          value: 'alarms_created',
        },
        {
          text: 'Alarms resolved',
          value: 'alarms_resolved',
        },
        {
          text: 'Alarms canceled',
          value: 'alarms_canceled',
        },
      ],
      items: [
        {
          entity: {
            _id: 'resource1/component1',
            type: 'resource',
            impact: [
              'service',
            ],
          },
          'Alarmes critiques': {
            value: 117,
            trend: 76,
            sla: false,
          },
          'Alarmes majeures': {
            value: 37,
            trend: 10,
          },
        },
        {
          entity: {
            _id: 'resource2/component1',
            type: 'resource',
            impact: [
              'service',
            ],
          },
          'Alarmes critiques': {
            value: 2,
            trend: -1,
            sla: true,
          },
          'Alarmes majeures': {
            value: 3,
            trend: -1,
          },
        },
      ],
      */
    };
  },
};
</script>
