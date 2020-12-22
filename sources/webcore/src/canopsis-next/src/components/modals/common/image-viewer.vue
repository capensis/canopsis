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
    display: block;
  }
}
</style>
