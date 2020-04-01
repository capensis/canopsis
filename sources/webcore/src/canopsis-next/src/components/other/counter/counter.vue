<template lang="pug">
  div.pa-2
    v-fade-transition(mode="out-in")
      v-layout(v-if="watchersPending", key="progress", column)
        v-flex(xs12)
          v-layout(justify-center)
            v-progress-circular(indeterminate, color="primary")
      v-layout.fill-height(v-else, key="content", wrap)
        v-alert(type="error", :value="true", v-if="hasNoData && watchersError")
          v-layout(align-center)
            div.mr-4 {{ $t('errors.default') }}
            v-tooltip(top)
              v-icon(slot="activator") help
              div(v-if="watchersError.name") {{ $t('common.name') }}: {{ watchersError.name }}
              div(v-if="watchersError.description") {{ $t('common.description') }}: {{ watchersError.description }}
        v-alert(v-else-if="hasNoData", type="info", :value="true") {{ $t('tables.noData') }}
        v-flex(v-else, v-for="item in watchers", :key="item._id", :class="flexSize")
          counter-item.weatherItem(
            :data-test="item.entity_id",
            :watcher="item",
            :widget="widget",
            :template="widget.parameters.blockTemplate"
          )
</template>

<script>
import { omit } from 'lodash';

import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import widgetQueryMixin from '@/mixins/widget/query';

import CounterItem from './counter-item.vue';

export default {
  components: {
    CounterItem,
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
