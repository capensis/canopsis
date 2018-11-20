<template lang="pug">
  div
    v-btn(icon, @click="showSettings")
      v-icon settings
    stats-curves(:labels="labels", :datasets="datasets")
</template>

<script>
import omit from 'lodash/omit';
import entitiesStatsMixin from '@/mixins/entities/stats';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
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
      const labels = [];
      if (this.stats.aggregations) {
        const stats = Object.keys(this.stats.aggregations);
        const values = { ...this.stats.aggregations };
        /*
        'start' correspond to the beginning timestamp.
        It's the same for all stats, that's why we can just take the first.
        We then give it to the date filter, to display it with a date format
        */
        values[stats[0]].sum.map(value => labels.push(this.$options.filters.date(value.start, 'long')));
        return labels;
      }
      return labels;
    },
    datasets() {
      return Object.keys(this.widget.parameters.stats).map((stat) => {
        let data = [];
        if (this.stats.aggregations) {
          data = this.stats.aggregations[stat].sum.map(value => value.value);
        }

        return {
          data,
          label: stat,
          borderColor: this.widget.parameters.statsColors ? this.widget.parameters.statsColors[stat] : '#DDDDDD',
          backgroundColor: 'transparent',
        };
      });
    },
  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: this.$constants.SIDE_BARS.statsCurvesSettings,
        config: {
          widget: this.widget,
          rowId: this.rowId,
        },
      });
    },
    async fetchList() {
      const stats = await this.fetchStatsEvolutionWithoutStore({
        params: omit(this.widget.parameters, ['statsColors']),
        aggregate: ['sum'],
      });
      this.stats = stats;
    },
  },
};
</script>
