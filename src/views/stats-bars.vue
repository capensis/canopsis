<template lang="pug">
  v-container
    v-layout
      v-flex(xs4)
        h2 {{ widget.title }}
        v-btn(icon, @click="showSettings")
          v-icon settings
        component(
        is="statsBars",
        :widget="widget",
        )
</template>

<script>
import StatsBars from '@/components/other/stats/stats-bars.vue';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import { SIDE_BARS } from '@/constants';

export default {
  components: {
    StatsBars,
  },
  mixins: [sideBarMixin],
  data() {
    return {
      widget: {
        _id: 'testStatsBars',
        title: 'Stats bars',
        parameters: {
          data: {
            labels: ['CLI', 'Services', 'Hôtes'],
            datasets: [
              {
                label: 'alarms_created',
                data: [12, 19, 3],
              },
              {
                label: 'alarms_resolved',
                data: [4, 7, 3],
              },
              {
                label: 'alarms_canceled',
                data: [1, 2, 4],
              },
            ],
          },
        },
      },
    };
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.statsHistogramSettings,
        config: {
          widget: this.widget,
        },
      });
    },
  },
};
</script>
