<template>
  <div class="pa-2">
    <v-fade-transition>
      <v-progress-linear
        v-if="countersPending"
        class="progress-linear-absolute--top"
        height="2"
        indeterminate
      />
    </v-fade-transition>

    <v-row class="fill-height" density="compact">
      <v-col v-if="hasNoData" cols="12">
        <v-alert type="info">
          {{ $t('common.noData') }}
        </v-alert>
      </v-col>

      <template v-else>
        <v-col
          v-for="counter in countersWithFilters"
          :key="counter.filter.title"
          :cols="columnSize"
        >
          <counter-item
            :counter="counter"
            :hide-actions="hideActions"
            :widget="widget"
            class="weatherItem"
          />
        </v-col>
      </template>
    </v-row>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { mapActions, mapState } from 'pinia';

import { useCounterStore } from '@/stores/modules/counter';

import CounterItem from './counter-item.vue';

export default {
  inject: ['$periodicRefresh'],

  components: {
    CounterItem,
  },

  props: {
    widget: {
      type: Object,
      required: true,
    },
    visible: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    ...mapState(useCounterStore, {
      getCountersListByWidgetId: 'getListByWidgetId',
      getCountersPendingByWidgetId: 'getPendingByWidgetId',
    }),

    columnsCount() {
      const displayName = this.$vuetify.display.name;

      if (['xs', 'sm'].includes(displayName)) {
        return this.widget.parameters.columnMobile;
      }

      if (displayName === 'md') {
        return this.widget.parameters.columnTablet;
      }

      return this.widget.parameters.columnDesktop;
    },

    columnSize() {
      return 12 / this.columnsCount;
    },

    counters() {
      return this.getCountersListByWidgetId(this.widget._id);
    },

    countersPending() {
      return this.getCountersPendingByWidgetId(this.widget._id);
    },

    hasNoData() {
      return this.counters.length === 0;
    },

    countersWithFilters() {
      const { filters = [] } = this.widget;

      return this.counters
        .filter((counter, index) => filters[index])
        .map((counter, index) => ({ ...counter, filter: filters[index] }));
    },

    query() {
      const { filters = [], parameters = {} } = this.widget;

      return {
        opened: parameters.opened,
        correlation: parameters.isCorrelationEnabled,
        filters: filters.map(({ _id: id }) => id),
      };
    },

    queryWithoutFilters() {
      return omit(this.query, ['filters']);
    },
  },

  watch: {
    visible: {
      immediate: true,
      handler(visible) {
        if (visible) {
          this.fetchList();
        }
      },
    },

    query() {
      if (this.visible) {
        this.fetchList();
      }
    },
  },

  created() {
    this.$periodicRefresh.registerChild(this.refresh);
  },

  beforeUnmount() {
    this.$periodicRefresh.unregisterChild(this.refresh);
  },

  methods: {
    ...mapActions(useCounterStore, {
    fetchCountersList: 'fetchList',
  }),

    refresh() {
      if (this.visible) {
        return this.fetchList();
      }

      return Promise.resolve();
    },

    fetchList() {
      return this.fetchCountersList({
        widgetId: this.widget._id,
        filters: this.query.filters,
        params: this.queryWithoutFilters,
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
