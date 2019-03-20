<template lang="pug">
  div
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { STATS_DURATION_UNITS } from '@/constants';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetStatsQueryMixin from '@/mixins/widget/stats/stats-query';

import StatsTextStat from './stats-text-stat.vue';

export default {
  components: { VRuntimeTemplate, StatsTextStat },
  mixins: [entitiesStatsMixin, widgetStatsQueryMixin],
  props: {
    template: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      query: {
        mfilter: {},
        dateInterval: {
          periodValue: 1,
          periodUnit: STATS_DURATION_UNITS.day,
          tstart: 'now/d',
          tstop: 'now/d',
        },
        stats: {},
        statsColors: {},
      },
      stats: [],
      compiledTemplate: '',
    };
  },
  watch: {
    template: {
      immediate: true,
      handler() {
        this.stats = [];
        this.compiledTemplate = `<div>${compile(this.template)}</div>`;
        this.fetchList();
      },
    },
  },
  beforeCreate() {
    registerHelper('stat', ({ hash }) => {
      const statName = hash.name;

      this.stats.push(statName);

      return new Handlebars.SafeString(`<stats-text-stat name="${statName}"></stats-text-stat>`);
    });
  },
  beforeDestroy() {
    unregisterHelper('entities');
  },
  methods: {
    getQuery() {
      const {
        mfilter,
        tstop,
        duration,
      } = this.getStatsQuery();

      return {
        duration,
        mfilter,

        tstop: tstop.startOf('h').unix(),
        stats: this.stats.reduce((acc, key) => {
          acc[key] = {
            aggregate: ['sum'],
          };

          return acc;
        }, {}),
      };
    },

    async fetchList() {
      const response = await this.fetchStatsListWithoutStore({
        params: this.getQuery(),
      });

      console.warn('fetch LIST', [...this.stats], response);
    },
  },
};
</script>
