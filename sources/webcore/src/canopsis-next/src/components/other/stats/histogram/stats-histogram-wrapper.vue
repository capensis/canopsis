<template lang="pug">
  div.stats-wrapper
    stats-histogram(:labels="labels", :datasets="datasets")
</template>

<script>
import { get, omit, isString } from 'lodash';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import StatsHistogram from './stats-histogram.vue';

export default {
  components: {
    StatsHistogram,
  },
  mixins: [entitiesStatsMixin, widgetQueryMixin, entitiesUserPreferenceMixin],
  props: {
    widget: {
      type: Object,
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
    fetchList() {
      this.widget.parameters.groups.map(async (group) => {
        let filter = get(group, 'filter.filter', {});

        if (isString(filter)) {
          filter = JSON.parse(filter);
        }

        const stat = await this.fetchStatsListWithoutStore({
          params: {
            ...omit(this.widget.parameters, ['groups', 'statsColors']),

            mfilter: filter,
          },
          aggregate: ['sum'],
        });
        this.$set(this.stats, group.title, stat);
      });
    },
  },
};
</script>
