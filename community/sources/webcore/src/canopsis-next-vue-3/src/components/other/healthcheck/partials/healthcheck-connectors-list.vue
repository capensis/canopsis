<template>
  <section class="healthcheck-connectors-list">
    <h3 class="text-h6 mb-3">
      {{ title }}
    </h3>
    <p v-if="!connectors.length" class="text-body-2 text-medium-emphasis ma-0">
      {{ $t('healthcheck.connectorsBlocks.noConnectors') }}
    </p>
    <c-layout v-else class="ga-4" wrap>
      <healthcheck-connector-tile
        v-for="connector in connectors"
        :key="connector.id"
        :connector="connector"
        @refresh="$emit('refresh')"
      />
    </c-layout>
  </section>
</template>

<script>
import HealthcheckConnectorTile from './healthcheck-connector-tile.vue';

export default {
  components: { HealthcheckConnectorTile },

  props: {
    title: {
      type: String,
      required: true,
    },
    connectors: {
      type: Array,
      required: true,
    },
  },

  emits: ['refresh'],
};
</script>
