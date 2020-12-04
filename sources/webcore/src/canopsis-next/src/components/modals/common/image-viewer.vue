<template lang="pug">
  img(:src="config.src")
</template>

<script>
import modalInnerMixin from '@/mixins/modal/inner';

export default {
  mixins: [modalInnerMixin],
};
</script>

<style lang="scss">
.v-dialog__image-viewer {
  width: auto;
  overflow: hidden;

  img {
    display: block;
  }
}
</style>
