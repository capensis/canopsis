<template lang="pug">
  div
    c-page-header
    v-card.ma-4.mt-0
      v-tabs(v-model="activeTab", slider-color="primary", centered)
        v-tab(:href="`#${$constants.KPI_TABS.graphs}`") {{ $tc('common.graph', 2) }}
        v-tab(
          v-if="hasReadAnyKpiFiltersAccess",
          :href="`#${$constants.KPI_TABS.filters}`"
        ) {{ $t('common.filters') }}
        v-tab(
          v-if="hasReadAnyKpiRatingSettingsAccess",
          :href="`#${$constants.KPI_TABS.ratingSettings}`"
        ) {{ $t('common.ratingSettings') }}

      v-tabs-items(v-model="activeTab")
        v-card-text
          v-tab-item(:value="$constants.KPI_TABS.graphs")
            kpi-charts
          v-tab-item(:value="$constants.KPI_TABS.filters", lazy)
            kpi-filters
          v-tab-item(:value="$constants.KPI_TABS.ratingSettings", lazy)
            kpi-rating-settings

    v-slide-x-reverse-transition
      c-fab-btn(
        v-if="hasFabButton",
        :has-access="hasAccessToCreate",
        @create="create",
        @refresh="refresh"
      )
        span {{ $t('modals.createFilter.create.title') }}
</template>

<script>
import { KPI_TABS, MODALS } from '@/constants';

import { entitiesFilterMixin } from '@/mixins/entities/filter';
import { entitiesRatingSettingsMixin } from '@/mixins/entities/rating-settings';
import { permissionsTechnicalKpiFiltersMixin } from '@/mixins/permissions/technical/kpi-filters';
import { permissionsTechnicalKpiRatingSettingsMixin } from '@/mixins/permissions/technical/kpi-rating-settings';

import KpiCharts from '@/components/other/kpi/charts/kpi-charts.vue';
import KpiFilters from '@/components/other/kpi/filters/kpi-filters.vue';
import KpiRatingSettings from '@/components/other/kpi/rating-settings/kpi-rating-settings.vue';

export default {
  components: { KpiRatingSettings, KpiCharts, KpiFilters },
  mixins: [
    entitiesFilterMixin,
    entitiesRatingSettingsMixin,
    permissionsTechnicalKpiFiltersMixin,
    permissionsTechnicalKpiRatingSettingsMixin,
  ],
  data() {
    return {
      activeTab: KPI_TABS.graphs,
    };
  },
  computed: {
    hasFabButton() {
      return [KPI_TABS.filters, KPI_TABS.ratingSettings].includes(this.activeTab);
    },

    hasAccessToCreate() {
      return this.activeTab === KPI_TABS.filters && this.hasCreateAnyKpiFiltersAccess;
    },
  },
  methods: {
    refresh() {
      switch (this.activeTab) {
        case KPI_TABS.filters:
          this.fetchFiltersListWithPreviousParams();
          break;
        case KPI_TABS.ratingSettings:
          this.fetchRatingSettingsListWithPreviousParams();
          break;
      }
    },

    create() {
      if (this.activeTab === KPI_TABS.filters) {
        this.showCreateFilterModal();
      }
    },

    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createKpiFilter,
        config: {
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
