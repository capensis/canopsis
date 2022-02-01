<template lang="pug">
  div.healthcheck
    c-progress-overlay(:pending="pending")
    c-page-header
    div.mb-3.text-xs-center(v-if="hasAnyError")
      v-chip(color="error", dark)
        span.subheading {{ $t('healthcheck.systemStatusChipError') }}
    v-tabs(v-model="activeTab", centered, slider-color="primary")
      v-tab {{ $t('common.systemStatus') }}
      v-tab {{ $tc('common.graph', 2) }}
      v-tab {{ $tc('common.parameter', 2) }}
    v-tabs-items.white.healthcheck__tabs(v-model="activeTab")
      v-tab-item.healthcheck__graph-tab
        healthcheck-network-graph(
          v-if="!pending && !hasServerError",
          :services="services",
          :engines-graph="enginesGraph",
          :engines-parameters="enginesParameters",
          :has-invalid-engines-order="hasInvalidEnginesOrder",
          :max-queue-length="maxQueueLength",
          show-description
        )
        h2.my-4.headline.text-xs-center(v-else-if="hasServerError") {{ $t('healthcheck.systemStatusServerError') }}
      v-tab-item(lazy)
        healthcheck-graphs(:max-queue-length="maxQueueLength")
      v-tab-item(lazy)
        healthcheck-parameters
</template>

<script>
import { isEqual } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { SOCKET_ROOMS } from '@/config';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';
import HealthcheckGraphs from '@/components/other/healthcheck/exploitation/healthcheck-graphs.vue';
import HealthcheckParameters from '@/components/other/healthcheck/healthcheck-parameters.vue';

import entitiesEngineRunInfoMixin from '@/mixins/entities/engine-run-info';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  components: { HealthcheckParameters, HealthcheckNetworkGraph, HealthcheckGraphs },
  mixins: [entitiesEngineRunInfoMixin],
  data() {
    return {
      activeTab: 0,
      pending: true,
      services: [],
      enginesGraph: {},
      enginesParameters: {},
      hasInvalidEnginesOrder: false,
      maxQueueLength: 0,
      hasServerError: false,
    };
  },
  computed: {
    hasAnyError() {
      return this.hasServerError || this.hasInvalidEnginesOrder;
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
        max_queue_length: maxQueueLength = 0,
      } = data;

      const preparedData = {
        services,
        enginesGraph,
        enginesParameters,
        hasInvalidEnginesOrder,
        maxQueueLength,
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
.healthcheck {
  display: flex;
  flex-direction: column;
  height: 100%;

  &__tabs {
    flex: 1;

    /deep/ .v-window__container {
      height: 100%;
    }
  }

  &__graph-tab {
    height: 100%;
  }
}
</style>
