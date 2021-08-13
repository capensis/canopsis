<template lang="pug">
  div.position-relative.fill-height
    v-layout(column, fill-height)
      c-page-header
      c-progress-overlay(:pending="healthcheckPending")
      v-flex
        healthcheck-network-graph(
          v-if="!healthcheckPending",
          :services="services",
          :engines="engines",
          :has-invalid-engines-order="hasInvalidEnginesOrder",
          @click="showNodeModal"
        )
      c-fab-btn(@refresh="fetchList")
</template>

<script>
import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';

import { entitiesHealthcheckMixin } from '@/mixins/entities/healthcheck';

export default {
  components: { HealthcheckNetworkGraph },
  mixins: [entitiesHealthcheckMixin],
  mounted() {
    this.fetchList();
  },
  methods: {
    showNodeModal() {},

    fetchList() {
      this.fetchHealthcheckStatus();
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
