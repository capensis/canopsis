<template lang="pug">
  v-icon(large :color="color") assistant_photo
</template>

<script>
export default {
  name: 'state-flag',
  data() {
    return {
      correspondingColorAndText: {
        0: 'grey',
        1: 'yellow darken-1',
        2: 'orange',
        3: 'red',
      },
    };
  },
  props: {
    val: Number,
    default: 0,
  },
  computed: {
    color() {
      return this.correspondingColorAndText[this.val];
    },
  },
};
</script>

<style scoped>

</style>
