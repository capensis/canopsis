<template lang="pug">
  div.position-relative
    progress-overlay(:pending="pending")
    alert-overlay(
      :value="hasError",
      :message="errorMessage"
    )
    v-runtime-template(:template="compiledTemplate")
</template>

<script>
import { isEmpty } from 'lodash';
import Handlebars from 'handlebars';
import VRuntimeTemplate from 'v-runtime-template';

import { CANOPSIS_EDITION } from '@/constants';

import { compile, registerHelper, unregisterHelper } from '@/helpers/handlebars';

import widgetQueryMixin from '@/mixins/widget/query';
import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetStatsQueryMixin from '@/mixins/widget/stats/stats-query';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';
import AlertOverlay from '@/components/layout/alert/alert-overlay.vue';

import TextStatTemplate from './text-stat-template.vue';

export default {
  components: {
    VRuntimeTemplate,
    ProgressOverlay,
    AlertOverlay,
    TextStatTemplate,
  },
  mixins: [
    widgetQueryMixin,
    entitiesStatsMixin,
    widgetStatsQueryMixin,
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
      pending: true,
      serverErrorMessage: null,
      stats: {},
    };
  },
  computed: {
    compiledTemplate() {
      return `<div>${compile(this.widget.parameters.template)}</div>`;
    },
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
