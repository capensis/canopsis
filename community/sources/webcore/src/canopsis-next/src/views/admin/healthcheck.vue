<template>
  <div class="healthcheck">
    <c-progress-overlay :pending="pending" />
    <c-page-header />
    <div
      v-if="hasAnyError"
      class="mb-3 text-center"
    >
      <v-chip
        color="error"
        dark
      >
        <span class="text-subtitle-1">{{ $t('healthcheck.systemStatusChipError') }}</span>
      </v-chip>
    </div>
    <v-sheet>
      <v-tabs
        v-model="activeTab"
        slider-color="primary"
        centered
      >
        <v-tab>{{ $t('common.systemStatus') }}</v-tab>
        <v-tab>{{ $tc('common.graph', 2) }}</v-tab>
        <v-tab>{{ $tc('common.parameter', 2) }}</v-tab>
        <v-tab v-if="hasAccessToTechMetrics">
          {{ $t('common.enginesMetrics') }}
        </v-tab>
      </v-tabs>
      <v-tabs-items
        v-model="activeTab"
        class="healthcheck__tabs"
      >
        <v-tab-item class="healthcheck__graph-tab">
          <healthcheck-network-graph
            v-if="!pending && !hasServerError"
            :services="services"
            :engines-graph="enginesGraph"
            :engines-parameters="enginesParameters"
            :has-invalid-engines-order="hasInvalidEnginesOrder"
            :max-queue-length="maxQueueLength"
            class="healthcheck__graph"
            show-description
          />
          <h2
            v-else-if="hasServerError"
            class="my-4 text-h5 text-center"
          >
            {{ $t('healthcheck.systemStatusServerError') }}
          </h2>
        </v-tab-item>
        <v-tab-item>
          <healthcheck-graphs :max-messages-per-minute="maxMessagesLength" />
        </v-tab-item>
        <v-tab-item>
          <healthcheck-parameters />
        </v-tab-item>
        <v-tab-item>
          <healthcheck-engines-metrics />
        </v-tab-item>
      </v-tabs-items>
    </v-sheet>
  </div>
</template>

<script>
import { isEqual } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { SOCKET_ROOMS } from '@/config';
import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/healthcheck-network-graph.vue';
import HealthcheckGraphs from '@/components/other/healthcheck/healthcheck-graphs.vue';
import HealthcheckParameters from '@/components/other/healthcheck/healthcheck-parameters.vue';
import HealthcheckEnginesMetrics from '@/components/other/healthcheck/healthcheck-engines-metrics.vue';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  components: {
    HealthcheckParameters,
    HealthcheckNetworkGraph,
    HealthcheckGraphs,
    HealthcheckEnginesMetrics,
  },
  mixins: [authMixin],
  data() {
    return {
      activeTab: 0,
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
    hasAnyError() {
      return this.hasServerError || this.hasInvalidEnginesOrder;
    },

    hasAccessToTechMetrics() {
      return this.checkAccess(USERS_PERMISSIONS.technical.techmetrics);
    },
  },
  mounted() {
    this.fetchList();

    this.$socket
      .join(SOCKET_ROOMS.healthcheck)
      .addListener(this.setData);
  },
  beforeDestroy() {
    this.$socket
      .leave(SOCKET_ROOMS.healthcheck)
      .removeListener(this.setData);
  },
  methods: {
    ...mapActions({
      fetchHealthcheckEnginesWithoutStore: 'fetchEnginesWithoutStore',
    }),

    setData(data) {
      const {
        services = [],
        engines: {
          graph: enginesGraph = {},
          parameters: enginesParameters = {},
        },
        has_invalid_engines_order: hasInvalidEnginesOrder = false,
        max_queue_length: maxQueueLength,
        max_messages_length: maxMessagesLength,
      } = data;

      const preparedData = {
        services,
        enginesGraph,
        enginesParameters,
        hasInvalidEnginesOrder,
        maxQueueLength: maxQueueLength || Infinity,
        maxMessagesLength: maxMessagesLength || Infinity,
      };

      Object.entries(preparedData).forEach(([key, value]) => {
        if (!isEqual(this[key], value)) {
          this[key] = value;
        }
      });
    },

    async fetchList() {
      try {
        this.hasServerError = false;
        this.pending = true;

        const data = await this.fetchHealthcheckEnginesWithoutStore();

        this.setData(data);
      } catch (err) {
        this.hasServerError = true;
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
$networkGraphHeight: calc(100vh - 230px); // 230px was calculated by imperative method
$networkGraphMinHeight: 100px;

.healthcheck {
  display: flex;
  flex-direction: column;

  &__tabs {
    flex: 1;

    ::v-deep .v-window__container {
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
}
</style>
