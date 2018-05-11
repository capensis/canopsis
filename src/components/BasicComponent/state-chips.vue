<template lang="pug">
  v-chip(:color="color" text-color="white") {{ text }}
</template>

<script>
import getProp from 'lodash/get';

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
      return getProp(this.correspondingColorAndText, `${this.val}.color`);
    },
    text() {
      return getProp(this.correspondingColorAndText, `${this.val}.text`);
    },
  },
};
</script>

<style scoped>

</style>
