<template>
  <c-layout class="healthcheck-connectors position-relative pa-4 ga-4" column>
    <c-progress-overlay :pending="pending" />

    <c-layout class="mb-4" justify="center">
      <v-btn :to="manageConnectorsRoute" color="primary" variant="outlined">
        <v-icon start>
          settings
        </v-icon>
        {{ $t('healthcheck.connectorsBlocks.manageMonitoredConnectors') }}
      </v-btn>
    </c-layout>

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
      @update:pagination="updatePaginationOptions"
    />
  </c-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import { ANOMALY_MONITORED_CONNECTOR_STATUSES, ROUTES_NAMES } from '@/constants';

import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useAnomalyMonitoredConnector } from '@/hooks/store/modules/anomaly-monitored-connector';

import { useHealthcheckConnectorsSocket } from './hooks/healthcheck-connectors-socket';
import HealthcheckConnectorsList from './partials/healthcheck-connectors-list.vue';

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

    const filteredConnectors = computed(() => connectors.value.reduce(
        (acc, connector) => {
          if (connector.status === ANOMALY_MONITORED_CONNECTOR_STATUSES.disabled) {
            acc.disabled.push(connector);
          } else {
            acc.enabled.push(connector);
          }

          return acc;
        },
        { enabled: [], disabled: [] },
      ),
    );

    const updatePage = page => updateOptions({ ...options.value, page });

    const updateItemsPerPage = itemsPerPage => updateOptions({ ...options.value, itemsPerPage });

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
