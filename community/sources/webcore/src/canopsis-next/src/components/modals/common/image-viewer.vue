<template>
  <img
    :src="config.src"
    alt=""
  >
</template>

<script>
import { modalInnerMixin } from '@/mixins/modal/inner';

export default {
  mixins: [modalInnerMixin],
};
</script>

<style lang="scss">
.v-dialog__image-viewer {
  width: auto;
  overflow: hidden;

  img {
    max-height: 90vh;
    max-width: 90vw;
    display: block;
  }
}
</style>
