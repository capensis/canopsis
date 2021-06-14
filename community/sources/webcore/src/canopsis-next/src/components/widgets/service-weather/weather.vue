<template lang="pug">
  div.pa-2
    v-fade-transition(v-if="watchersPending", key="progress", mode="out-in")
      v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
    v-layout.fill-height(key="content", wrap)
      v-alert(v-if="hasNoData && watchersError", :value="true", type="error")
        v-layout(align-center)
          div.mr-4 {{ $t('errors.default') }}
          v-tooltip(top)
            v-icon(slot="activator") help
            div(v-if="watchersError.name") {{ $t('common.name') }}: {{ watchersError.name }}
            div(v-if="watchersError.description") {{ $t('common.description') }}: {{ watchersError.description }}
      v-alert(v-else-if="hasNoData", :value="true", type="info") {{ $t('tables.noData') }}
      template(v-else)
        v-flex(v-for="item in watchers", :key="item._id", :class="flexSize")
          weather-item.weatherItem(
            :data-test="item._id",
            :watcher="item",
            :widget="widget",
            :template="widget.parameters.blockTemplate"
          )
</template>

<script>
import { omit } from 'lodash';

import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';

import WeatherItem from './weather-item.vue';

export default {
  components: {
    WeatherItem,
  },
  mixins: [
    widgetPeriodicRefreshMixin,
    entitiesWatcherMixin,
    widgetFetchQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    flexSize() {
      return [
        `xs${this.widget.parameters.columnSM}`,
        `md${this.widget.parameters.columnMD}`,
        `lg${this.widget.parameters.columnLG}`,
      ];
    },
    hasNoData() {
      return this.watchers.length === 0;
    },
  },
  methods: {
    getQuery() {
      const query = omit(this.query, [
        'page',
        'sortKey',
        'sortDir',
      ]);

      if (this.query.sortKey) {
        query.sort_by = this.query.sortKey;
        query.sort = this.query.sortDir.toLowerCase();
      }

      return query;
    },

    fetchList() {
      this.fetchWatchersList({
        filter: this.widget.parameters.mfilter.filter,
        params: this.getQuery(),
        widgetId: this.widget._id,
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .weatherItem {
    height: 100%;
  }
</style>
