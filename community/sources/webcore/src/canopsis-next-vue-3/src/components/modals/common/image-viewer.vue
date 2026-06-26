<template>
  <img
    :src="config.src"
    alt=""
  >
</template>

<script>
export default {
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
  },
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
