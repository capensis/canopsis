<template lang="pug">
  div
    progress-overlay(:pending="pending")
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { compile, registerHelper, unregisterHelper } from '../../../helpers/handlebars';

import widgetQueryMixin from '../../../mixins/widget/query';
import entitiesStatsMixin from '../../../mixins/entities/stats';
import widgetStatsQueryMixin from '../../../mixins/widget/stats/stats-query';

import ProgressOverlay from '../../layout/progress/progress-overlay.vue';

import TextStatTemplate from './text-stat-template.vue';

export default {
  components: { VRuntimeTemplate, ProgressOverlay, TextStatTemplate },
  mixins: [
    widgetQueryMixin,
    entitiesStatsMixin,
    widgetStatsQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
      stats: {},
    };
  },
  computed: {
    compiledTemplate() {
      return `<div>${compile(this.widget.parameters.template)}</div>`;
    },
  },
  beforeCreate() {
    registerHelper('stat', ({ hash }) => {
      const statName = hash.name;

      return new Handlebars.SafeString(`
        <text-stat-template name="${statName}" :stats="stats"></text-stat-template>
      `);
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
        stats,
        duration,
      } = this.getStatsQuery();

      return {
        duration,
        mfilter,
        tstop: tstop.startOf('h').unix(),
        stats: Object.entries(stats).reduce((acc, [key, value]) => {
          acc[key] = {
            ...value,

            aggregate: ['sum'],
          };

          return acc;
        }, {}),
      };
    },

    async fetchList() {
      this.pending = true;

      const { aggregations } = await this.fetchStatsListWithoutStore({
        params: this.getQuery(),
      });

      this.stats = aggregations;
      this.pending = false;
    },
  },
};
</script>
