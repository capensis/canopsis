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
        0: {
          color: 'grey',
          text: 'ongoing',
        },
        1: {
          color: 'yellow darken-1',
          text: 'minor',
        },
        2: {
          color: 'orange',
          text: 'major',
        },
        3: {
          color: 'red',
          text: 'critical',
        },
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
