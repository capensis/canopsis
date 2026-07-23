<template>
  <div class="healthcheck">
    <c-progress-overlay :pending="pending" />
    <c-page-header :name="USER_PERMISSIONS.technical.healthcheck" />

    <div v-if="hasAnyError" class="mb-3 text-center">
      <v-chip color="error">
        <span class="text-subtitle-1">
          {{ $t('healthcheck.systemStatusChipError') }}
        </span>
      </v-chip>
    </div>

    <v-sheet>
      <v-tabs v-model="activeTab" color="primary" align-tabs="center">
        <v-tab :value="HEALTHCHECK_TABS.systemStatus">
          {{ $t('common.systemStatus') }}
        </v-tab>
        <v-tab v-if="hasAccessToAnomalyMonitoredConnectors" :value="HEALTHCHECK_TABS.connectors">
          {{ $t('healthcheck.connectors') }}
        </v-tab>
        <v-tab :value="HEALTHCHECK_TABS.graphs">
          {{ $t('common.graph', 2) }}
        </v-tab>
        <v-tab :value="HEALTHCHECK_TABS.parameters">
          {{ $t('common.parameter', 2) }}
        </v-tab>
        <v-tab v-if="hasAccessToTechMetrics" :value="HEALTHCHECK_TABS.enginesMetrics">
          {{ $t('common.enginesMetrics') }}
        </v-tab>
      </v-tabs>

      <v-window v-model="activeTab" class="healthcheck__tabs">
        <v-window-item :value="HEALTHCHECK_TABS.systemStatus" class="healthcheck__graph-tab">
          <div v-if="!pending && !hasServerError" class="healthcheck__graph">
            <healthcheck-network-graph
              :services="services"
              :engines-graph="enginesGraph"
              :engines-parameters="enginesParameters"
              :has-invalid-engines-order="hasInvalidEnginesOrder"
              :max-queue-length="maxQueueLength"
              show-description
            />
          </div>
          <h2 v-else-if="hasServerError" class="my-4 text-h5 text-center">
            {{ $t('healthcheck.systemStatusServerError') }}
          </h2>
        </v-window-item>

        <v-window-item
          v-if="hasAccessToAnomalyMonitoredConnectors"
          :value="HEALTHCHECK_TABS.connectors"
        >
          <div class="healthcheck__placeholder">
            {{ $t('healthcheck.connectors') }}
          </div>
        </v-window-item>

        <v-window-item :value="HEALTHCHECK_TABS.graphs">
          <healthcheck-graphs :max-messages-per-minute="maxMessagesLength" />
        </v-window-item>

        <v-window-item :value="HEALTHCHECK_TABS.parameters">
          <healthcheck-parameters />
        </v-window-item>

        <v-window-item v-if="hasAccessToTechMetrics" :value="HEALTHCHECK_TABS.enginesMetrics">
          <healthcheck-engines-metrics />
        </v-window-item>
      </v-window>
    </v-sheet>
  </div>
</template>

<script>
import isEqual from 'lodash/isEqual';

import { SOCKET_ROOMS } from '@/config';
import { CRUD_ACTIONS, USER_PERMISSIONS } from '@/constants';

import { checkUserAccess } from '@/helpers/permission';

import { useAuthStore } from '@/stores/modules/auth';
import { useHealthcheckStore } from '@/stores/modules/healthcheck';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/healthcheck-network-graph.vue';
import HealthcheckEnginesMetrics from '@/components/other/healthcheck/healthcheck-engines-metrics.vue';
import HealthcheckGraphs from '@/components/other/healthcheck/healthcheck-graphs.vue';
import HealthcheckParameters from '@/components/other/healthcheck/healthcheck-parameters.vue';

const HEALTHCHECK_TABS = {
  systemStatus: 'system-status',
  connectors: 'connectors',
  graphs: 'graphs',
  parameters: 'parameters',
  enginesMetrics: 'engines-metrics',
};

export default {
  components: {
    HealthcheckEnginesMetrics,
    HealthcheckGraphs,
    HealthcheckNetworkGraph,
    HealthcheckParameters,
  },

  data() {
    return {
      HEALTHCHECK_TABS,
      USER_PERMISSIONS,
      activeTab: HEALTHCHECK_TABS.systemStatus,
      pending: true,
      services: [],
      enginesGraph: {},
      enginesParameters: {},
      hasInvalidEnginesOrder: false,
      maxQueueLength: Infinity,
      maxMessagesLength: Infinity,
      hasServerError: false,
    };
  },

  computed: {
    currentUserPermissionsById() {
      return useAuthStore().currentUserPermissionsById;
    },

    hasAccessToTechMetrics() {
      return checkUserAccess(
        this.currentUserPermissionsById[USER_PERMISSIONS.technical.techmetrics],
        CRUD_ACTIONS.can,
      );
    },

    hasAccessToAnomalyMonitoredConnectors() {
      return checkUserAccess(
        this.currentUserPermissionsById[USER_PERMISSIONS.technical.anomalyMonitoredConnector],
        CRUD_ACTIONS.can,
      );
    },

    availableTabs() {
      return [
        HEALTHCHECK_TABS.systemStatus,
        this.hasAccessToAnomalyMonitoredConnectors && HEALTHCHECK_TABS.connectors,
        HEALTHCHECK_TABS.graphs,
        HEALTHCHECK_TABS.parameters,
        this.hasAccessToTechMetrics && HEALTHCHECK_TABS.enginesMetrics,
      ].filter(Boolean);
    },

    hasAnyError() {
      return this.hasServerError || this.hasInvalidEnginesOrder;
    },
  },

  watch: {
    availableTabs(tabs) {
      if (!tabs.includes(this.activeTab)) {
        this.activeTab = tabs[0] || HEALTHCHECK_TABS.systemStatus;
      }
    },
  },

  mounted() {
    this.fetchList();
    this.$socket.join(SOCKET_ROOMS.healthcheck).addListener(this.setData);
  },

  beforeUnmount() {
    this.$socket.leave(SOCKET_ROOMS.healthcheck).removeListener(this.setData);
  },

  methods: {
    updateIfChanged(currentValue, updateValue, nextValue) {
      if (!isEqual(currentValue, nextValue)) {
        updateValue(nextValue);
      }
    },

    setData(data = {}) {
      const {
        services: nextServices = [],
        engines: { graph: nextEnginesGraph = {}, parameters: nextEnginesParameters = {} } = {},
        has_invalid_engines_order: nextHasInvalidEnginesOrder = false,
        max_queue_length: nextMaxQueueLength,
        max_messages_length: nextMaxMessagesLength,
      } = data;

      this.updateIfChanged(
        this.services,
        (value) => {
          this.services = value;
        },
        nextServices,
      );
      this.updateIfChanged(
        this.enginesGraph,
        (value) => {
          this.enginesGraph = value;
        },
        nextEnginesGraph,
      );
      this.updateIfChanged(
        this.enginesParameters,
        (value) => {
          this.enginesParameters = value;
        },
        nextEnginesParameters,
      );
      this.updateIfChanged(
        this.hasInvalidEnginesOrder,
        (value) => {
          this.hasInvalidEnginesOrder = value;
        },
        nextHasInvalidEnginesOrder,
      );
      this.updateIfChanged(
        this.maxQueueLength,
        (value) => {
          this.maxQueueLength = value;
        },
        nextMaxQueueLength || Infinity,
      );
      this.updateIfChanged(
        this.maxMessagesLength,
        (value) => {
          this.maxMessagesLength = value;
        },
        nextMaxMessagesLength || Infinity,
      );
    },

    async fetchList() {
      try {
        this.hasServerError = false;
        this.pending = true;

        const data = await useHealthcheckStore().fetchHealthcheckEnginesWithoutStore();

        this.setData(data);
      } catch {
        this.hasServerError = true;
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
$networkGraphHeight: calc(100vh - 230px);
$networkGraphMinHeight: 100px;

.healthcheck {
  display: flex;
  flex-direction: column;

  &__tabs {
    flex: 1;

    :deep(.v-window__container) {
      height: 100%;
    }
  }

  &__graph-tab {
    height: 100%;
  }

  &__graph {
    height: $networkGraphHeight;
    min-height: $networkGraphMinHeight;
  }

  &__placeholder {
    min-height: $networkGraphMinHeight;
    padding: 16px;
  }
}
</style>
