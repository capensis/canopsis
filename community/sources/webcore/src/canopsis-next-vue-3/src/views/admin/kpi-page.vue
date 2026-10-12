<template>
  <div>
    <c-progress-overlay :pending="pending" />
    <c-page-header :name="USER_PERMISSIONS.technical.kpi" />

    <v-card class="ma-4 mt-0">
      <v-tabs v-model="activeTab" color="primary" align-tabs="center">
        <v-tab :value="KPI_TABS.graphs">
          {{ $t('common.graph', 2) }}
        </v-tab>
        <v-tab v-if="hasReadAnyKpiFiltersAccess" :value="KPI_TABS.filters">
          {{ $t('common.filters') }}
        </v-tab>
        <v-tab v-if="hasReadAnyKpiRatingSettingsAccess" :value="KPI_TABS.ratingSettings">
          {{ $t('kpi.tabs.ratingSettings') }}
        </v-tab>
        <v-tab v-if="hasReadAnyKpiCollectionSettingsAccess" :value="KPI_TABS.collectionSettings">
          {{ $t('kpi.tabs.collectionSettings') }}
        </v-tab>
      </v-tabs>

      <v-card-text>
        <v-window v-model="activeTab">
          <v-window-item :value="KPI_TABS.graphs">
            <div v-if="!timescaleAvailable" class="text-error text-center">
              {{ $t('kpi.metricsNotAvailable') }}
            </div>
            <kpi-charts :unavailable="!timescaleAvailable" />
          </v-window-item>

          <v-window-item v-if="hasReadAnyKpiFiltersAccess" :value="KPI_TABS.filters">
            <kpi-filters />
          </v-window-item>

          <v-window-item v-if="hasReadAnyKpiRatingSettingsAccess" :value="KPI_TABS.ratingSettings">
            <kpi-rating-settings />
          </v-window-item>

          <v-window-item
            v-if="hasReadAnyKpiCollectionSettingsAccess"
            :value="KPI_TABS.collectionSettings"
          >
            <v-row>
              <v-col
                cols="12"
                md="10"
                offset-md="1"
                lg="8"
                offset-lg="2"
              >
                <kpi-collection-settings />
              </v-col>
            </v-row>
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <c-fab-btn
      v-if="hasFabButton"
      :key="activeTab"
      :has-access="hasAccessToCreate"
      v-bind="fabBind"
    >
      <span>{{ $t('modals.createFilter.create.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { SOCKET_ROOMS } from '@/config';
import {
  HEALTHCHECK_SERVICES_NAMES,
  HEALTHCHECK_STATUSES,
  KPI_TABS,
  MODALS,
  USER_PERMISSIONS,
} from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useFilter } from '@/hooks/store/modules/filter';
import { useHealthcheck } from '@/hooks/store/modules/healthcheck';
import { useRatingSettings } from '@/hooks/store/modules/rating-settings';

import KpiCharts from '@/components/other/kpi/charts/kpi-charts.vue';
import KpiCollectionSettings from '@/components/other/kpi/collection-settings/kpi-collection-settings.vue';
import KpiFilters from '@/components/other/kpi/filters/kpi-filters.vue';
import KpiRatingSettings from '@/components/other/kpi/rating-settings/kpi-rating-settings.vue';

export default {
  components: {
    KpiCharts,
    KpiCollectionSettings,
    KpiFilters,
    KpiRatingSettings,
  },

  setup() {
    const {
      hasCreateAccess: hasCreateAnyKpiFiltersAccess,
      hasReadAccess: hasReadAnyKpiFiltersAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.kpiFilters);
    const { hasReadAccess: hasReadAnyKpiRatingSettingsAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.kpiRatingSettings,
    );
    const { hasReadAccess: hasReadAnyKpiCollectionSettingsAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.kpiCollectionSettings,
    );

    const { fetchFiltersListWithPreviousParams } = useFilter();
    const { createFilter } = useFilter();
    const { fetchRatingSettingsListWithPreviousParams } = useRatingSettings();
    const { fetchHealthcheckEnginesWithoutStore } = useHealthcheck();

    return {
      fetchFiltersListWithPreviousParams,
      createFilter,
      fetchHealthcheckEnginesWithoutStore,
      fetchRatingSettingsListWithPreviousParams,
      hasCreateAnyKpiFiltersAccess,
      hasReadAnyKpiCollectionSettingsAccess,
      hasReadAnyKpiFiltersAccess,
      hasReadAnyKpiRatingSettingsAccess,
    };
  },

  data() {
    return {
      KPI_TABS,
      USER_PERMISSIONS,
      activeTab: KPI_TABS.graphs,
      pending: false,
      timescaleAvailable: false,
    };
  },

  computed: {
    availableTabs() {
      return [
        KPI_TABS.graphs,
        this.hasReadAnyKpiFiltersAccess && KPI_TABS.filters,
        this.hasReadAnyKpiRatingSettingsAccess && KPI_TABS.ratingSettings,
        this.hasReadAnyKpiCollectionSettingsAccess && KPI_TABS.collectionSettings,
      ].filter(Boolean);
    },

    hasFabButton() {
      return [KPI_TABS.filters, KPI_TABS.ratingSettings].includes(this.activeTab);
    },

    hasAccessToCreate() {
      return this.activeTab === KPI_TABS.filters && this.hasCreateAnyKpiFiltersAccess;
    },

    fabBind() {
      return {
        onRefresh: this.refresh,
        ...(this.activeTab === KPI_TABS.filters ? { onCreate: this.create } : {}),
      };
    },
  },

  watch: {
    availableTabs(tabs) {
      if (!tabs.includes(this.activeTab)) {
        this.activeTab = tabs[0] || KPI_TABS.graphs;
      }
    },
  },

  mounted() {
    this.fetchHealthcheck();

    this.$socket.join(SOCKET_ROOMS.healthcheck).addListener(this.setTimescaleIsAvailable);
  },

  beforeUnmount() {
    this.$socket.leave(SOCKET_ROOMS.healthcheck).removeListener(this.setTimescaleIsAvailable);
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
      if (this.activeTab !== KPI_TABS.filters) {
        return;
      }

      this.showCreateFilterModal();
    },

    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createKpiFilter,
        config: {
          action: async (data) => {
            const result = await this.createFilter({ data });

            this.fetchFiltersListWithPreviousParams();

            return result;
          },
        },
      });
    },

    setTimescaleIsAvailable({ services = [] } = {}) {
      const timeScaleService = services.find(
        ({ name }) => name === HEALTHCHECK_SERVICES_NAMES.timescaleDB,
      );

      this.timescaleAvailable = [
        HEALTHCHECK_STATUSES.healthy,
        HEALTHCHECK_STATUSES.degraded,
      ].includes(timeScaleService?.status);
    },

    async fetchHealthcheck() {
      try {
        this.pending = true;

        const data = await this.fetchHealthcheckEnginesWithoutStore();

        this.setTimescaleIsAvailable(data);
      } catch (err) {
        console.error(err);

        this.timescaleAvailable = false;
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.kpi-page__placeholder {
  padding: 24px;
  text-align: center;
}
</style>
