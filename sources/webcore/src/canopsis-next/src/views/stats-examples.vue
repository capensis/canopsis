<template lang="pug">
  v-container
    div
      h2 {{ widget.title }}
      component(
      is="statsTable",
      :widget="widget",
      )
</template>

<script>
import StatsTable from '@/components/other/stats/stats-table.vue';

export default {
  components: {
    StatsTable,
  },
  data() {
    return {
      widget: {
        _id: 'testStatsTable',
        title: 'Stats table',
        parameters: {
          mfilter: {
            type: 'component',
          },
          tstop: 1534716000,
          duration: '2d',
          stats: {
            'Taux de disponibilité': {
              stat: 'state_rate',
              parameters: {
                states: [0, 1, 2],
              },
              trend: true,
              sla: '>= 0.99',
            },
            'Taux de performance': {
              stat: 'state_rate',
              parameters: {
                states: [0],
              },
              trend: true,
              sla: '>= 0.95',
            },
            'Durée d\'indisponibilité': {
              stat: 'time_in_state',
              parameters: {
                states: [3],
              },
              trend: true,
            },
            'Nombre d\'indisponibilités': {
              stat: 'alarms_created',
              parameters: {
                states: [3],
              },
              trend: true,
              sla: '<= 100',
            },
          },
        },
      },
    };
  },
};
</script>

