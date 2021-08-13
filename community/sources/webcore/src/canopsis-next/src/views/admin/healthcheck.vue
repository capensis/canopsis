<template lang="pug">
  div.position-relative.fill-height
    v-layout(column, fill-height)
      c-page-header
      c-progress-overlay(:pending="pending")
      v-flex
        healthcheck-network-graph(
          v-if="!pending",
          :services="services",
          :engines="engines",
          :has-invalid-engines-order="hasInvalidEnginesOrder",
          @click="showNodeModal"
        )
      c-fab-btn(@refresh="fetchList")
</template>

<script>
import { HEALTHCHECK_SERVICES_NAMES, MODALS } from '@/constants';

import { createNamespacedHelpers } from 'vuex';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';

import entitiesEngineRunInfoMixin from '@/mixins/entities/engine-run-info';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  components: { HealthcheckNetworkGraph },
  mixins: [entitiesEngineRunInfoMixin],
  data() {
    return {
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

    showNodeModal(engine) {
      const excludedServices = [
        HEALTHCHECK_SERVICES_NAMES.api,
        HEALTHCHECK_SERVICES_NAMES.healthcheck,
        HEALTHCHECK_SERVICES_NAMES.events,
      ];

      if (excludedServices.includes(engine.id)) {
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
#cy {
  background: white;
  position: relative;
  width: 100%;
  height: 100vh;
  cursor: grabbing;
}
</style>
