<template>
  <modal-wrapper
    class="fill-height"
    close
  >
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <div class="pre-wrap">
        {{ $t('healthcheck.chainConfigurationInvalid') }}
      </div>
      <v-fade-transition>
        <v-layout
          v-if="pending"
          justify-center
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </v-layout>
        <div
          v-else
          class="healthcheck-engine-chain-reference"
        >
          <healthcheck-network-graph
            :engines-graph="enginesGraph"
            :engines-parameters="enginesParameters"
          />
        </div>
      </v-fade-transition>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.ok') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, HEALTHCHECK_SERVICES_NAMES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/healthcheck-network-graph.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  name: MODALS.healthcheckEnginesChainReference,
  components: { HealthcheckNetworkGraph, ModalWrapper },
  mixins: [modalInnerMixin, healthcheckNodesMixin, entitiesInfoMixin],
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
    ...mapActions({
      fetchEnginesOrderWithoutStore: 'fetchEnginesOrderWithoutStore',
    }),

    async fetchList() {
      try {
        this.pending = true;
        this.enginesGraph = await this.fetchEnginesOrderWithoutStore();
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
