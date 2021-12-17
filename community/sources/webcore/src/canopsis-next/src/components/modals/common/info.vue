<template lang="pug">
  modal-wrapper(close)
    template(#title)
      span {{ config.title }}
    template(#text)
      div(v-html="config.text")
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.info,
  components: {
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
  ],
};
</script>
