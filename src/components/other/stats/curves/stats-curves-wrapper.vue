<template lang="pug">
  v-container(fluid)
    v-btn(icon, @click="showSettings")
      v-icon settings
    stats-curves(:labels="labels", :datasets="datasets")
</template>

<script>
import entitiesStatsMixin from '@/mixins/entities/stats';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { SIDE_BARS } from '@/constants';
import StatsCurves from './stats-curves.vue';

export default {
  components: {
    StatsCurves,
  },
  mixins: [entitiesStatsMixin, sideBarMixin, widgetQueryMixin, entitiesUserPreferenceMixin],
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
  data() {
    return {
      stats: {},
    };
  },
  computed: {
    labels() {
      return [];
    },
    datasets() {
      return [];
    },
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.statsCurvesSettings,
        config: {
          widget: this.widget,
          rowId: this.rowId,
        },
      });
    },
    async fetchList() {
      const stats = await this.fetchStatsEvolutionWithoutStore({
        params: { ...this.widget.parameters },
        aggregate: ['sum'],
      });
      this.stats = stats.aggregations;
    },
  },
};
</script>
