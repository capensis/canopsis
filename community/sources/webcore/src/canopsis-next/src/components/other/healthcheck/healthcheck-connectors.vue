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
      :connectors="enabledConnectors"
      @refresh="fetchList"
    />
    <healthcheck-connectors-list
      :title="$t('healthcheck.connectorsBlocks.disabledConnectors')"
      :connectors="disabledConnectors"
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
import { isEqual } from 'lodash';
import {
  computed,
  watch,
  set,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { SOCKET_ROOMS } from '@/config';
import { ANOMALY_MONITORED_CONNECTOR_STATUSES, ROUTES_NAMES, USER_PERMISSIONS } from '@/constants';

import { mapIds } from '@/helpers/array';

import { useSocket } from '@/hooks/socket';
import { useCanPermission } from '@/hooks/auth';
import { useAnomalyMonitoredConnectors } from '@/hooks/store/modules/anomaly-monitored-connector';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

import HealthcheckConnectorsList from '@/components/other/healthcheck/partials/healthcheck-connectors-list.vue';

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

    const enabledConnectors = computed(() => connectors.value.filter(
      ({ status }) => status !== ANOMALY_MONITORED_CONNECTOR_STATUSES.disabled,
    ));

    const disabledConnectors = computed(() => connectors.value.filter(
      ({ status }) => status === ANOMALY_MONITORED_CONNECTOR_STATUSES.disabled,
    ));

    const updatePage = page => updateOptions({ ...options.value, page });

    const socket = useSocket();

    const socketListener = (data) => {
      data.data.forEach((connector, index) => {
        if (isEqual(connector, connectors.value[index])) {
          return;
        }

        set(connectors.value, index, connector);
      });

      if (meta.value.total_count !== data.total_count) {
        set(meta.value, 'total_count', data.total_count);
      }
    };

    const joinToSocketRoom = (ids = []) => socket
      .join(SOCKET_ROOMS.anomalyMonitoredConnectorStates, { ids }, true)
      .addListener(socketListener);

    const leaveFromSocketRoom = () => socket
      .leave(SOCKET_ROOMS.anomalyMonitoredConnectorStates)
      .removeListener(socketListener);

    const reconnectToSocketRoom = (ids = []) => {
      leaveFromSocketRoom();
      joinToSocketRoom(ids);
    };

    watch(connectors, (items, prevItems) => {
      const newIds = mapIds(items, 'id');
      const oldIds = mapIds(prevItems ?? [], 'id');

      if (isEqual(newIds, oldIds)) {
        return;
      }

      reconnectToSocketRoom(newIds);
    }, { immediate: true });

    onMounted(fetchList);
    onBeforeUnmount(leaveFromSocketRoom);

    return {
      meta,
      pending,
      options,
      enabledConnectors,
      disabledConnectors,
      hasAccessToAnomalyMonitoredConnectors,
      manageConnectorsRoute,

      updatePage,
      fetchList,
    };
  },
};
</script>
