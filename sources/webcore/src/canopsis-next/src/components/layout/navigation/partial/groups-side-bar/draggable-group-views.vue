<template lang="pug">
  draggable.views-panel.secondary.lighten-1(
    v-field="views",
    :class="{ empty: isViewsEmpty }",
    :options="draggableOptions"
  )
    group-view-panel(
      v-for="view in views",
      :key="view._id",
      :view="view"
    )
</template>

<script>
import Draggable from 'vuedraggable';
import { VUETIFY_ANIMATION_DELAY } from '@/config';

import GroupViewPanel from './group-view-panel.vue';

export default {
  components: { GroupViewPanel, Draggable },
  model: {
    prop: 'views',
    event: 'change',
  },
  props: {
    views: {
      type: Array,
      required: true,
    },
  },
  computed: {
    isViewsEmpty() {
      return this.views.length === 0;
    },

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        group: 'views',
      };
    },
  },
};
</script>

<style lang="scss" scoped>
  .views-panel .empty {
    &:after {
      content: '';
      display: block;
      height: 48px;
      border: 4px dashed #4f6479;
      border-radius: 5px;
      position: relative;
    }
  }
</style>
