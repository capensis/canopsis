<template lang="pug">
  div
    c-page-header

    v-card.ma-4.mt-0
      v-layout.pa-4(v-if="pending", justify-center)
        v-progress-circular(color="primary", indeterminate)
      template(v-else)
        v-tabs(v-model="activeTab", slider-color="primary", centered)
          v-tab(:href="`#${$constants.KPI_TABS.graphs}`") {{ $tc('common.graph', 2) }}
          v-tab(
            v-if="hasReadAnyKpiFiltersAccess",
            :href="`#${$constants.KPI_TABS.filters}`"
          ) {{ $t('common.filters') }}
          v-tab(
            v-if="hasReadAnyKpiRatingSettingsAccess",
            :href="`#${$constants.KPI_TABS.ratingSettings}`"
          ) {{ $t('kpi.tabs.ratingSettings') }}
          v-tab(
            v-if="hasReadAnyKpiCollectionSettingsAccess",
            :href="`#${$constants.KPI_TABS.collectionSettings}`"
          ) {{ $t('kpi.tabs.collectionSettings') }}

        v-tabs-items(v-model="activeTab")
          v-card-text
            v-tab-item(:value="$constants.KPI_TABS.graphs")
              div.error--text.text-xs-center(v-if="!timescaleAvailable") {{ $t('kpi.metricsNotAvailable') }}
              kpi-charts(:unavailable="!timescaleAvailable")
            v-tab-item(:value="$constants.KPI_TABS.filters", lazy)
              kpi-filters
            v-tab-item(:value="$constants.KPI_TABS.ratingSettings", lazy)
              kpi-rating-settings
            v-tab-item(:value="$constants.KPI_TABS.collectionSettings", lazy)
              v-layout(row)
                v-flex(xs12, offset-md1, md10, offset-lg2, lg8)
                  kpi-collection-settings

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
import { createNamespacedHelpers } from 'vuex';

import { SOCKET_ROOMS } from '@/config';

import { HEALTHCHECK_SERVICES_NAMES, KPI_TABS, MODALS } from '@/constants';

import { entitiesFilterMixin } from '@/mixins/entities/filter';
import { entitiesRatingSettingsMixin } from '@/mixins/entities/rating-settings';
import { permissionsTechnicalKpiFiltersMixin } from '@/mixins/permissions/technical/kpi-filters';
import { permissionsTechnicalKpiRatingSettingsMixin } from '@/mixins/permissions/technical/kpi-rating-settings';
import { permissionsTechnicalKpiCollectionSettingsMixin } from '@/mixins/permissions/technical/kpi-collection-settings';

import KpiCharts from '@/components/other/kpi/charts/kpi-charts.vue';
import KpiFilters from '@/components/other/kpi/filters/kpi-filters.vue';
import KpiRatingSettings from '@/components/other/kpi/rating-settings/kpi-rating-settings.vue';
import KpiCollectionSettings from '@/components/other/kpi/collection-settings/kpi-collection-settings.vue';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  components: { KpiRatingSettings, KpiCharts, KpiFilters, KpiCollectionSettings },
  mixins: [
    entitiesFilterMixin,
    entitiesRatingSettingsMixin,
    permissionsTechnicalKpiFiltersMixin,
    permissionsTechnicalKpiRatingSettingsMixin,
    permissionsTechnicalKpiCollectionSettingsMixin,
  ],
  data() {
    return {
      activeTab: KPI_TABS.graphs,
      pending: false,
      timescaleAvailable: false,
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
  mounted() {
    this.fetchHealthcheck();

    this.$socket
      .join(SOCKET_ROOMS.healthcheck)
      .addListener(this.setTimescaleIsAvailable);
  },
  beforeDestroy() {
    this.$socket
      .leave(SOCKET_ROOMS.healthcheck)
      .removeListener(this.setTimescaleIsAvailable);
  },
  methods: {
    ...mapActions({
      fetchHealthcheckEnginesWithoutStore: 'fetchEnginesWithoutStore',
    }),

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

    setTimescaleIsAvailable({ services }) {
      const timeScaleService = services.find(({ name }) => name === HEALTHCHECK_SERVICES_NAMES.timescaleDB);

      this.timescaleAvailable = !!timeScaleService?.is_running;
    },

    async fetchHealthcheck() {
      try {
        this.pending = true;

        const data = await this.fetchHealthcheckEnginesWithoutStore();

        this.setTimescaleIsAvailable(data);
      } catch (err) {
        this.timescaleAvailable = false;
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
