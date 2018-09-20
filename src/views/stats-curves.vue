<template lang="pug">
  v-container
    v-layout
      v-flex(xs4)
        h2 {{ widget.title }}
        component(
        is="statsCurves",
        :widget="widget",
        )
</template>

<script>
import StatsCurves from '@/components/other/stats/stats-curves.vue';

export default {
  components: {
    StatsCurves,
  },
  data() {
    return {
      widget: {
        _id: 'testStatsCurves',
        title: 'Stats curves',
        parameters: {
          data: {
            labels: ['CLI', 'Services', 'Hôtes'],
            datasets: [{
              label: 'alarms_created',
              data: [12, 19, 3],
            }],
          },
        },
      },
    };
  },
};
</script>

