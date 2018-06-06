<template lang="pug">
  transition(name="fade")
    div(v-show="pending")
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
    transition: opacity .5s;
  }
  .fade-enter, .fade-leave-to {
    opacity: 0;
  }
</style>
