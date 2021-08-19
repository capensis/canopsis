<template lang="pug">
  div.healthcheck
    c-progress-overlay(:pending="pending")
    c-page-header

    v-tabs(v-model="activeTab", centered, slider-color="primary")
      v-tab {{ $t('common.systemStatus') }}
      v-tab {{ $tc('common.graph', 2) }}
      v-tab {{ $tc('common.parameter', 2) }}
    v-tabs-items.white.healthcheck__tabs(v-model="activeTab")
      v-tab-item.healthcheck__graph-tab
        healthcheck-network-graph(
          v-if="!pending",
          :services="services",
          :engines="engines",
          :has-invalid-engines-order="hasInvalidEnginesOrder",
          show-description,
          @click="showNodeModal"
        )
      v-tab-item(lazy)
        healthcheck-graphs(:max-queue-length="maxQueueLength")
      v-tab-item(lazy)
        healthcheck-parameters

    c-fab-btn(@refresh="fetchList")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { HEALTHCHECK_SERVICES_NAMES, MODALS } from '@/constants';

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
      maxQueueLength: 0,
      pending: true,
      services: [],
      engines: {},
      hasInvalidEnginesOrder: false,
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchHealthcheckStatusWithoutStore: 'fetchStatusWithoutStore',
    }),

    showEngineModal(engine) {
      const excludedServices = [
        HEALTHCHECK_SERVICES_NAMES.api,
        HEALTHCHECK_SERVICES_NAMES.healthcheck,
        HEALTHCHECK_SERVICES_NAMES.events,
      ];

      if (excludedServices.includes(engine.name)) {
        return;
      }

      const hasError = !engine.is_running
        || engine.is_queue_overflown
        || engine.is_too_few_instances
        || engine.is_diff_instances_config;

      if (hasError) {
        this.$modals.show({
          name: MODALS.healthcheckEngine,
          config: {
            engine,
          },
        });
      }
    },

    showEngineChainReferenceModal() {
      this.$modals.show({
        name: MODALS.healthcheckEngineChainReference,
      });
    },

    showNodeModal(engine) {
      if (engine.name !== HEALTHCHECK_SERVICES_NAMES.enginesChain) {
        this.showEngineModal(engine);
      } else if (this.hasInvalidEnginesOrder) {
        this.showEngineChainReferenceModal(engine);
      }
    },

    async fetchList() {
      this.pending = true;

      const {
        services = [],
        engines = {},
        max_queue_length: maxQueueLength,
        has_invalid_engines_order: hasInvalidEnginesOrder,
      } = await this.fetchHealthcheckStatusWithoutStore();

      this.services = services;
      this.hasInvalidEnginesOrder = hasInvalidEnginesOrder;
      this.maxQueueLength = maxQueueLength;
      this.engines = {
        edges: engines.edges,
        nodes: engines.nodes.map(node => ({ ...node, max_queue_length: maxQueueLength })),
      };
      this.pending = false;
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

#cy {
  background: white;
  position: relative;
  width: 100%;
  height: 100vh;
  cursor: grabbing;
}
</style>
