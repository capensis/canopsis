<template lang="pug">
  div
    v-btn(@click="openModal") Modal
    div#brand Canopsis Next
</template>

<script>
export default {
  methods: {
    openModal() {
      this.$store.dispatch('modal/show', {
        name: 'confirmation',
        config: {
          number: 0,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  #brand {
    text-align: center;
    position: relative;
    top: 25%;
    max-width: 50%;
    max-height: 5em;
    margin: auto;
    font-weight: bold;
    font-size: 2em;
  }
</style>
