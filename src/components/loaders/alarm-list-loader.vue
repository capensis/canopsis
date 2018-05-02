<template lang="pug">
    vue-content-loading( :width="300" :height="100")
      template(v-for="i in elementsCount")
        rect(:y="(i - 1) * 15.2" width="100%"  rx="1" ry="1" height="15")

</template>

<script>
import VueContentLoading from 'vue-content-loading';

export default {
  name: 'AlarmListLoader',
  components: {
    VueContentLoading,
  },
  props: {
    elementsCount: {
      type: Number,
      default: 6,
    },
  },
};
</script>

<style scoped>

</style>
