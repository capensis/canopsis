<template>
  <v-layout
    class="position-relative pa-4 gap-4"
    column
  >
    <c-progress-overlay :pending="pending" />
    <v-layout
      v-if="hasAccessToAnomalyMonitoredConnectors"
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
    <c-pagination
      :page="options.page"
      :limit="options.itemsPerPage"
      :total="meta.total_count"
      @input="updatePage"
    />
  </v-layout>
</template>

<script>
import { computed, onMounted } from 'vue';

import { ANOMALY_MONITORED_CONNECTOR_STATUSES, ROUTES_NAMES, USER_PERMISSIONS } from '@/constants';

import { useCanPermission } from '@/hooks/auth';
import { useAnomalyMonitoredConnectors } from '@/hooks/store/modules/anomaly-monitored-connector';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

import HealthcheckConnectorsList from '@/components/other/healthcheck/partials/healthcheck-connectors-list.vue';

import { useHealthcheckConnectorsSocket } from './hooks/healthcheck-connectors-socket';

export default {
  components: { HealthcheckConnectorsList },
  setup() {
    const { hasAccess: hasAccessToAnomalyMonitoredConnectors } = useCanPermission(
      USER_PERMISSIONS.technical.anomalyMonitoredConnector,
    );

    const { fetchAnomalyMonitoredConnectorStatesListWithoutStore } = useAnomalyMonitoredConnectors();

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

    useHealthcheckConnectorsSocket({ connectors, meta });

    onMounted(fetchList);

    return {
      meta,
      pending,
      options,
      filteredConnectors,
      hasAccessToAnomalyMonitoredConnectors,
      manageConnectorsRoute,

      updatePage,
      fetchList,
    };
  },
};
</script>
