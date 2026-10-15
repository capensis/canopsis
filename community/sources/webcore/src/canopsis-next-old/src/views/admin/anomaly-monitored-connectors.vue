<template>
  <c-page
    :creatable="hasCreateAnyAnomalyMonitoredConnectorAccess"
    :create-tooltip="$t('modals.createAnomalyMonitoredConnector.create.title')"
    @refresh="fetchList"
    @create="showCreateAnomalyMonitoredConnectorModal"
  >
    <anomaly-monitored-connectors-list
      :anomaly-monitored-connectors="anomalyMonitoredConnectors"
      :pending="pending"
      :options="options"
      :total-items="meta.total_count"
      :updatable="hasUpdateAnyAnomalyMonitoredConnectorAccess"
      :removable="hasDeleteAnyAnomalyMonitoredConnectorAccess"
      @edit="showEditAnomalyMonitoredConnectorModal"
      @remove="showRemoveAnomalyMonitoredConnectorModal"
      @update:options="updateOptions"
      @refresh="fetchList"
    />
  </c-page>
</template>

<script>
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useAnomalyMonitoredConnector } from '@/hooks/store/modules/anomaly-monitored-connector';
import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useCRUDPermissions } from '@/hooks/auth';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

import AnomalyMonitoredConnectorsList from '@/components/other/anomaly-monitored-connector/anomaly-monitored-connectors-list.vue';

export default {
  components: { AnomalyMonitoredConnectorsList },
  setup() {
    const { t } = useI18n();
    const modals = useModals();

    const {
      hasCreateAccess: hasCreateAnyAnomalyMonitoredConnectorAccess,
      hasUpdateAccess: hasUpdateAnyAnomalyMonitoredConnectorAccess,
      hasDeleteAccess: hasDeleteAnyAnomalyMonitoredConnectorAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.anomalyMonitoredConnector);

    const {
      createAnomalyMonitoredConnector,
      updateAnomalyMonitoredConnector,
      removeAnomalyMonitoredConnector,
      fetchAnomalyMonitoredConnectorsListWithoutStore,
    } = useAnomalyMonitoredConnector();

    const { callActionWithPopup } = useCallActionWithPopup();

    const {
      data: anomalyMonitoredConnectors,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchAnomalyMonitoredConnectorsListWithoutStore,
    });

    /**
     * Show modal for creating a new anomaly monitored connector
     */
    const showCreateAnomalyMonitoredConnectorModal = () => {
      modals.show({
        name: MODALS.createAnomalyMonitoredConnector,
        config: {
          action: newConnector => callActionWithPopup(
            () => createAnomalyMonitoredConnector({ data: newConnector }),
            fetchList,
          ),
        },
      });
    };

    /**
     * Show modal for editing an existing anomaly monitored connector
     *
     * @param {Object} connector - The anomaly monitored connector to edit
     */
    const showEditAnomalyMonitoredConnectorModal = (connector) => {
      modals.show({
        name: MODALS.createAnomalyMonitoredConnector,
        config: {
          connector,
          title: t('modals.createAnomalyMonitoredConnector.edit.title'),

          action: newConnector => callActionWithPopup(
            () => updateAnomalyMonitoredConnector({ id: connector.id, data: newConnector }),
            fetchList,
          ),
        },
      });
    };

    /**
     * Show confirmation modal for removing an anomaly monitored connector
     *
     * @param {Object} connector - The anomaly monitored connector to remove
     */
    const showRemoveAnomalyMonitoredConnectorModal = (connector) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('modals.createAnomalyMonitoredConnector.remove.text', { connectorName: connector.name }),
          action: () => callActionWithPopup(
            () => removeAnomalyMonitoredConnector({ id: connector.id }),
            fetchList,
          ),
        },
      });
    };

    onMounted(fetchList);

    return {
      hasCreateAnyAnomalyMonitoredConnectorAccess,
      hasUpdateAnyAnomalyMonitoredConnectorAccess,
      hasDeleteAnyAnomalyMonitoredConnectorAccess,
      anomalyMonitoredConnectors,
      meta,
      pending,
      options,

      updateOptions,
      showCreateAnomalyMonitoredConnectorModal,
      showEditAnomalyMonitoredConnectorModal,
      showRemoveAnomalyMonitoredConnectorModal,
      fetchList,
    };
  },
};
</script>
