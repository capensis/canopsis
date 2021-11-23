<template lang="pug">
  div
    c-page-header
    v-card.ma-4.mt-0
      v-tabs(v-model="activeTab", slider-color="primary", centered)
        v-tab(:href="`#${$constants.KPI_TABS.graphs}`") {{ $tc('common.graph', 2) }}
        v-tab(v-if="hasReadAnyKpiFiltersAccess", :href="`#${$constants.KPI_TABS.filters}`") {{ $t('common.filters') }}

      v-tabs-items(v-model="activeTab")
        v-card-text
          v-tab-item(:value="$constants.KPI_TABS.graphs")
            kpi-charts
          v-tab-item(:value="$constants.KPI_TABS.filters", lazy)
            kpi-filters

    v-slide-x-reverse-transition
      c-fab-btn(
        v-if="isFilterTab",
        @create="showCreateFilterModal",
        @refresh="fetchFiltersListWithPreviousParams",
        :has-access="hasCreateAnyKpiFiltersAccess"
      )
        span {{ $t('modals.filter.create.title') }}
</template>

<script>
import { KPI_TABS, MODALS } from '@/constants';

import { entitiesFilterMixin } from '@/mixins/entities/filter';
import { permissionsTechnicalKpiFiltersMixin } from '@/mixins/permissions/technical/kpi-filters';

import KpiCharts from '@/components/other/kpi/charts/kpi-charts.vue';
import KpiFilters from '@/components/other/kpi/filters/kpi-filters.vue';

export default {
  components: { KpiCharts, KpiFilters },
  mixins: [entitiesFilterMixin, permissionsTechnicalKpiFiltersMixin],
  data() {
    return {
      activeTab: KPI_TABS.graphs,
    };
  },
  computed: {
    isFilterTab() {
      return this.activeTab === KPI_TABS.filters;
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.patterns,
        config: {
          title: this.$t('modals.filter.create.title'),
          name: true,
          entity: true,
          patterns: {
            name: '',
            entity_patterns: [],
          },
          action: async (data) => {
            await this.createFilter({ data });

            this.fetchFiltersListWithPreviousParams();
          },
        },
      });
    },
  },
};
</script>
