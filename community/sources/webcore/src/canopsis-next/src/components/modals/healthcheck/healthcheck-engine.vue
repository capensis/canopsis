<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ title }}
    template(slot="text")
      healthcheck-engine-information(:engine="engine")
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
</template>

<script>
import { MODALS } from '@/constants';

import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';

import HealthcheckEngineInformation from '@/components/other/healthcheck/healthcheck-engine-information.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.healthcheckEngine,
  components: { HealthcheckEngineInformation, ModalWrapper },
  mixins: [healthcheckNodesMixin],
  computed: {
    engine() {
      return this.config.engine;
    },

    title() {
      return this.getNodeName(this.engine.name);
    },
  },
};
</script>
