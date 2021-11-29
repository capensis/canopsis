<template lang="pug">
  modal-wrapper(data-test="colorPickerModal", close)
    template(slot="title")
      span {{ config.title }}
    template(slot="text")
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
