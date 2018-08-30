<template lang="pug">
  transition(name="fade", :duration="duration")
    slot
</template>

<script>
export default {
  props: {
    duration: {
      type: Number,
      default: 1000,
    },
  },
};
</script>

<style scoped>
  .fade-enter-active, .fade-leave-active {
    transition: opacity .5s;
  }
  .fade-enter, .fade-leave-to {
    opacity: 0;
  }
</style>
