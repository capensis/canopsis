<template lang="pug">
  div
    v-btn(icon, @click="showSettings")
      v-icon settings
    stats-histogram(:labels="labels", :datasets="datasets")
</template>

<script>
import omit from 'lodash/omit';
import entitiesStatsMixin from '@/mixins/entities/stats';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import StatsHistogram from './stats-histogram.vue';

export default {
  components: {
    StatsHistogram,
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
      return Object.keys(this.stats);
    },
    datasets() {
      return Object.keys(this.widget.parameters.stats).map((stat) => {
        const data = Object.values(this.stats).reduce((acc, group) => {
          if (group.aggregations) {
            acc.push(group.aggregations[stat].sum);
          }

          return acc;
        }, []);

        return {
          data,
          label: stat,
          backgroundColor: this.widget.parameters.statsColors ? this.widget.parameters.statsColors[stat] : '#DDDDDD',
        };
      });
    },

  },
  methods: {
    showSettings() {
      this.showSideBar({
        name: this.$constants.SIDE_BARS.statsHistogramSettings,
        config: {
          widget: this.widget,
          rowId: this.rowId,
        },
      });
    },
    fetchList() {
      this.widget.parameters.groups.map(async (group) => {
        const stat = await this.fetchStatsListWithoutStore({
          params: {
            ...omit(this.widget.parameters, ['groups', 'statsColors']),
            mfilter: group.filter || {},
          },
          aggregate: ['sum'],
        });
        this.$set(this.stats, group.title, stat);
      });
    },
  },
};
</script>

