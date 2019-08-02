<template lang="pug">
  div
    v-layout
    v-fade-transition
      v-layout.fill-height(v-show="!watchersPending", wrap)
        v-flex(v-for="item in watchers", :key="item._id", :class="flexSize")
          weather-item.weatherItem(
          :watcher="item",
          :widget="widget",
          :template="widget.parameters.blockTemplate",
          :isEditingMode="isEditingMode",
        )
    v-fade-transition
      v-layout(v-show="watchersPending", column)
        v-flex(xs12)
          v-layout(justify-center)
            v-progress-circular(indeterminate, color="primary")
</template>

<script>
import { omit } from 'lodash';

import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import widgetQueryMixin from '@/mixins/widget/query';

import WeatherItem from './weather-item.vue';

export default {
  components: {
    WeatherItem,
  },
  mixins: [
    widgetPeriodicRefreshMixin,
    entitiesWatcherMixin,
    widgetQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
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
  },
  methods: {
    getQuery() {
      const query = omit(this.query, [
        'page',
        'sortKey',
        'sortDir',
      ]);

      if (this.query.sortKey) {
        query.orderby = this.query.sortKey;
        query.direction = this.query.sortDir;
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
