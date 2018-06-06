<template lang="pug">
  transition(name="fade")
    div.loader(v-show="pending")
      slot
</template>

<script>
export default {
  props: {
    pending: {
      type: Boolean,
      required: true,
    },
  },
};
</script>

<style scoped>
  .fade-enter-active, .fade-leave-active {
    transition: opacity .001ms;
  }
  .fade-enter, .fade-leave-to {
    opacity: 0;
  }
  .loader{
    position: absolute;
    width: 100%;
    top: 15%;
    z-index: 999;
  }
</style>
