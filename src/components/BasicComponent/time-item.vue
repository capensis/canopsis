<template lang="pug">
  .container
    p {{ dateToString }}
</template>

<script>
export default {
  name: 'time-item',
  props: {
    time: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    dateToString() {
      const date = new Date(this.time);
      return `${date.toLocaleString('fr-FR', {
        weekday: 'long', year: 'numeric', month: 'long', day: 'numeric',
      })
      } - ${
        date.toLocaleTimeString('fr-FR')}`;
    },
  },
};
</script>

<style scoped>
  .container {
    padding: 0;
  }
</style>
