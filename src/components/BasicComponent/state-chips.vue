<template lang="pug">
  v-chip(:color="color" text-color="white") {{ text }}
</template>

<script>
export default {
  name: 'state-chips',
  props: {
    val: Number,
    default: 0,
  },
  data() {
    return {
      correspondingColorAndText: {
        0: ['grey', 'ongoing'],
        1: ['yellow darken-1', 'minor'],
        2: ['orange', 'major'],
        3: ['red', 'critical'],
      },
    };
  },
  computed: {
    color() {
      return this.correspondingColorAndText[this.val][0];
    },
    text() {
      return this.correspondingColorAndText[this.val][1];
    },
  },
};
</script>

<style scoped>

</style>
