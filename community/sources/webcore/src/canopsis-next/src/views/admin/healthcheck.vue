<template lang="pug">
  div.position-relative
    c-page-header
    c-progress-overlay(:pending="pending")
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

    showNodeModal() {},

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
