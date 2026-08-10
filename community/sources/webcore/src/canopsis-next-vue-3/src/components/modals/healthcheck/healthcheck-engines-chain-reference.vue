<template>
  <c-modal-wrapper fill-height close>
    <template #title>
      <span>{{ title }}</span>
    </template>

    <template #text>
      <div class="pre-wrap">
        {{ $t('healthcheck.chainConfigurationInvalid') }}
      </div>
      <v-fade-transition>
        <div v-if="pending" class="d-flex justify-center">
          <v-progress-circular color="primary" indeterminate />
        </div>
        <div v-else class="healthcheck-engine-chain-reference">
          <healthcheck-network-graph
            :engines-graph="enginesGraph"
            :engines-parameters="enginesParameters"
          />
        </div>
      </v-fade-transition>
    </template>

    <template #actions>
      <v-btn variant="text" @click="closeModal">
        {{ $t('common.ok') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import HealthcheckNetworkGraph from '@/components/other/healthcheck/healthcheck-network-graph.vue';
import { HEALTHCHECK_SERVICES_NAMES } from '@/constants';
import { useHealthcheckStore } from '@/stores/modules/healthcheck';

export default {
  components: { HealthcheckNetworkGraph },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      pending: false,
      enginesGraph: {
        edges: [],
        nodes: [],
      },
    };
  },

  computed: {
    enginesParameters() {
      return this.enginesGraph.nodes.reduce((acc, name) => {
        acc[name] = { name, is_running: true };

        return acc;
      }, {});
    },

    title() {
      return this.getNodeName(HEALTHCHECK_SERVICES_NAMES.enginesChain);
    },
  },

  mounted() {
    this.fetchList();
  },

  methods: {
    closeModal() {
      this.$modals.hide({ id: this.modal.id });
    },

    getNodeName(name) {
      const nodeLabelKey = `healthcheck.nodes.${name}.name`;

      return this.$te(nodeLabelKey) ? this.$t(nodeLabelKey) : name;
    },

    async fetchList() {
      try {
        this.pending = true;
        this.enginesGraph = await useHealthcheckStore().fetchHealthcheckEnginesOrderWithoutStore();
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss">
.healthcheck-engine-chain-reference {
  height: 65vh;
}
</style>
