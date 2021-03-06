<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    c-alert-overlay(
      :value="hasError",
      :message="serverErrorMessage"
    )
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import { isEmpty } from 'lodash';
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { CANOPSIS_EDITION } from '@/constants';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';
import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetStatsQueryMixin from '@/mixins/widget/stats/stats-query';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';

import TextStatTemplate from './text-stat-template.vue';

export default {
  components: {
    VRuntimeTemplate,
    TextStatTemplate,
  },
  mixins: [
    entitiesStatsMixin,
    widgetStatsQueryMixin,
    widgetFetchQueryMixin,
    widgetStatsWrapperMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      stats: {},
    };
  },
  asyncComputed: {
    compiledTemplate: {
      async get() {
        const compiledTemplate = await compile(this.widget.parameters.template);

        return `<div>${compiledTemplate}</div>`;
      },
      default: '',
    },
  },
  computed: {
    /**
     * Check if there are 'stats' associated with the widget. As stats are only available with 'cat' edition
     * Override editionError computed prop from widgetStatsWrapperMixin
     */
    editionError() {
      return Object.keys(this.stats).length && this.edition === CANOPSIS_EDITION.core;
    },
  },
  beforeCreate() {
    registerHelper('stat', ({ hash }) => {
      const statName = hash.name;

      return new Handlebars.SafeString(`
        <text-stat-template v-if="editionError" name="${statName}" :stats="stats"></text-stat-template>
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
        stats = {},
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
      try {
        this.pending = true;
        this.serverErrorMessage = null;

        if (!isEmpty(this.widget.parameters.stats)) {
          const { aggregations } = await this.fetchStatsListWithoutStore({
            params: this.getQuery(),
          });

          this.stats = aggregations;
        }
      } catch (err) {
        this.serverErrorMessage = err.description || null;
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
