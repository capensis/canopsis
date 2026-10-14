<template>
  <v-layout
    class="position-relative pa-4 gap-4"
    column
  >
    <c-progress-overlay :pending="pending" />
    <v-layout
      class="mb-4"
      justify-center
    >
      <v-btn
        :to="manageConnectorsRoute"
        color="primary"
        outlined
      >
        <v-icon left>
          settings
        </v-icon>
        {{ $t('healthcheck.connectorsBlocks.manageMonitoredConnectors') }}
      </v-btn>
    </v-layout>
    <healthcheck-connectors-list
      :title="$t('healthcheck.connectorsBlocks.enabledConnectors')"
      :connectors="filteredConnectors.enabled"
      @refresh="fetchList"
    />
    <healthcheck-connectors-list
      :title="$t('healthcheck.connectorsBlocks.disabledConnectors')"
      :connectors="filteredConnectors.disabled"
      @refresh="fetchList"
    />
    <c-table-pagination
      :total-items="meta.total_count"
      :items-per-page="options.itemsPerPage"
      :page="options.page"
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
      @input="updatePaginationOptions"
    />
  </v-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import { ANOMALY_MONITORED_CONNECTOR_STATUSES, ROUTES_NAMES } from '@/constants';

import { useAnomalyMonitoredConnector } from '@/hooks/store/modules/anomaly-monitored-connector';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

import HealthcheckConnectorsList from '@/components/other/healthcheck/partials/healthcheck-connectors-list.vue';

import { useHealthcheckConnectorsSocket } from './hooks/healthcheck-connectors-socket';

export default {
  components: { HealthcheckConnectorsList },
  setup() {
    const { fetchAnomalyMonitoredConnectorStatesListWithoutStore } = useAnomalyMonitoredConnector();

    const {
      data: connectors,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchAnomalyMonitoredConnectorStatesListWithoutStore,
    });

    const manageConnectorsRoute = { name: ROUTES_NAMES.adminAnomalyMonitoredConnectors };

    const filteredConnectors = computed(() => connectors.value.reduce((acc, connector) => {
      if (connector.status === ANOMALY_MONITORED_CONNECTOR_STATUSES.disabled) {
        acc.disabled.push(connector);
      } else {
        acc.enabled.push(connector);
      }

      return acc;
    }, { enabled: [], disabled: [] }));

    /**
     * Updates the page option and fetches the list.
     *
     * @param {number} page - The page number to update.
     */
    const updatePage = page => updateOptions({ ...options.value, page });

    /**
     * Updates items per page; current page is adjusted in `c-table-pagination` when needed.
     *
     * @param {number} itemsPerPage - The new page size.
     */
    const updateItemsPerPage = itemsPerPage => updateOptions({ ...options.value, itemsPerPage });

    /**
     * Updates both page and items per page (e.g. when page is recomputed for a new page size).
     *
     * @param {Object} payload
     * @param {number} payload.page
     * @param {number} payload.itemsPerPage
     */
    const updatePaginationOptions = ({ page, itemsPerPage }) => updateOptions({
      ...options.value,
      page,
      itemsPerPage,
    });

    useHealthcheckConnectorsSocket({ connectors, meta });

    onMounted(fetchList);

    return {
      meta,
      pending,
      options,
      filteredConnectors,
      manageConnectorsRoute,

      updatePage,
      updateItemsPerPage,
      updatePaginationOptions,
      fetchList,
    };
  },
};
</script>
