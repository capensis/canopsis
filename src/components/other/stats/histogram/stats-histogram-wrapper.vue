<template lang="pug">
  v-container(fluid)
    v-btn(icon, @click="showSettings")
      v-icon settings
    stats-histogram
</template>

<script>
import entitiesStatsMixin from '@/mixins/entities/stats';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import { SIDE_BARS } from '@/constants';
import StatsHistogram from './stats-histogram.vue';

export default {
  components: {
    StatsHistogram,
  },
  mixins: [entitiesStatsMixin, sideBarMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    rowId: {
      type: String,
      required: true,
    },
  },
  mounted() {
    this.widget.parameters.groups.map((group) => {
      const params = {
        mfilter: group.filter,
        tstop: this.widget.parameters.tstop,
        duration: this.widget.parameters.duration,
        stats: this.widget.parameters.stats,
      };
      return this.fetchStats({ params });
    });
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.statsHistogramSettings,
        config: {
          widget: this.widget,
          rowId: this.rowId,
        },
      });
    },
  },
};
</script>

