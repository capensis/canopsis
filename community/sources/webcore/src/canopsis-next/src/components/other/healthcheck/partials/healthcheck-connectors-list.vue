<template>
  <section class="healthcheck-connectors-list">
    <h3 class="text-h6 mb-3">
      {{ title }}
    </h3>
    <p
      v-if="!connectors.length"
      class="text-body-2 text--secondary ma-0"
    >
      {{ $t('healthcheck.connectorsBlocks.noConnectors') }}
    </p>
    <v-layout
      v-else
      class="gap-4"
      wrap
    >
      <healthcheck-connector-tile
        v-for="connector in connectors"
        :key="connector.id"
        :connector="connector"
        @refresh="refresh"
      />
    </v-layout>
  </section>
</template>

<script>
import HealthcheckConnectorTile from '@/components/other/healthcheck/partials/healthcheck-connector-tile.vue';

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
  setup(props, { emit }) {
    const refresh = () => emit('refresh');

    return {
      refresh,
    };
  },
};
</script>
