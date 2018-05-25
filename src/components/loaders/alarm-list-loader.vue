<template lang="pug">
    vue-content-loading(:width="width", :height="height")
      template(v-for="i in elementsCount")
        rect(:y="(i - 1) * listItemLoaderHeight", width="100%", rx="1", ry="1", :height="listItemLoaderHeight")
</template>

<script>
import { ALARM_LIST_LOADER_HEIGHT, ALARM_LIST_LOADER_WIDTH } from '@/config';
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
      height: ALARM_LIST_LOADER_HEIGHT,
      width: ALARM_LIST_LOADER_WIDTH,
    };
  },
};
</script>

<style scoped>

</style>
