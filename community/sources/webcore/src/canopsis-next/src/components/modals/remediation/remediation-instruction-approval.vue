<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.remediationInstructionApproval.title') }}
      template(slot="text")
        span something
</template>

<script>
import { MODALS } from '@/constants';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.remediationInstructionApproval,
  components: { ModalWrapper },
};
</script>
