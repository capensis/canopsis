<template lang="pug">
  modal-wrapper(close)
    template(#title="")
      span {{ title }}
    template(#text="")
      healthcheck-engine-information(:engine="engine", :max-queue-length="maxQueueLength")
    template(#actions="")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.ok') }}
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';

import HealthcheckEngineInformation from '@/components/other/healthcheck/healthcheck-engine-information.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.healthcheckEngine,
  components: { HealthcheckEngineInformation, ModalWrapper },
  mixins: [modalInnerMixin, healthcheckNodesMixin],
  computed: {
    engine() {
      return this.config.engine;
    },

    maxQueueLength() {
      return this.config.maxQueueLength;
    },

    title() {
      return this.getNodeName(this.engine.name);
    },
  },
};
</script>
