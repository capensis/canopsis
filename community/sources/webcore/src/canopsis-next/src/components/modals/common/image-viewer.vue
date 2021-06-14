<template lang="pug">
  img(:src="config.src")
</template>

<script>
export default {};
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
