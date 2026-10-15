<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ title }}</span>
    </template>

    <template #text>
      <healthcheck-engine-information :engine="engine" :max-queue-length="maxQueueLength" />
    </template>

    <template #actions>
      <v-btn variant="text" @click="closeModal">
        {{ $t('common.ok') }}
      </v-btn>
    </template>
  </c-modal-wrapper>
</template>

<script>
import HealthcheckEngineInformation from '@/components/other/healthcheck/healthcheck-engine-information.vue';

export default {
  components: { HealthcheckEngineInformation },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },

    engine() {
      return this.config.engine ?? {};
    },

    maxQueueLength() {
      return this.config.maxQueueLength;
    },

    title() {
      return this.getNodeName(this.engine.name);
    },
  },

  methods: {
    closeModal() {
      this.$modals.hide({ id: this.modal.id });
    },

    getNodeName(name) {
      const nodeLabelKey = `healthcheck.nodes.${name}.name`;

      return this.$te(nodeLabelKey) ? this.$t(nodeLabelKey) : name;
    },
  },
};
</script>
