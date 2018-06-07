<template lang="pug">
    vue-content-loading(:width="width", :height="height")
      template(v-for="i in elementsCount")
        rect(:y="(i - 1) * listItemLoaderHeight", width="100%", rx="1", ry="1", :height="listItemLoaderHeight")
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
  data() {
    return {
      listItemLoaderHeight: 15,
      height: 100,
      width: 300,
    };
  },
};
</script>

<style scoped>

</style>
