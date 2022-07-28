<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createMap.title') }}
      template(#text="")
        | Create map
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMap,
  components: {
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
  ],
};
</script>
