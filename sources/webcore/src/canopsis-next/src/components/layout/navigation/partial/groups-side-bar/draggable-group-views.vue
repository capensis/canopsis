<template lang="pug">
  draggable.views-panel.secondary.lighten-1(
    :value="views",
    :class="{ empty: isViewsEmpty }",
    :options="draggableOptions",
    @change="changeViewsOrdering"
  )
    group-view-panel(
      v-for="view in views",
      :key="view._id",
      :view="view"
    )
</template>

<script>
import Draggable from 'vuedraggable';

import { getDuplicateCountItems } from '@/helpers/searching';
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
        group: { name: 'views' },
      };
    },
  },
  methods: {
    changeViewsOrdering({ moved, added, removed }) {
      const views = [...this.views];

      if (moved) {
        const [item] = views.splice(moved.oldIndex, 1);

        views.splice(moved.newIndex, 0, item);
      } else if (added) {
        const duplicateViewsCount = getDuplicateCountItems(views, added.element);

        const view = duplicateViewsCount !== 0
          ? { ...added.element, name: `${added.element.name}_${duplicateViewsCount}` }
          : added.element;

        views.splice(added.newIndex, 0, view);
      } else if (removed) {
        views.splice(removed.oldIndex, 1);
      }

      if (views) {
        this.$emit('change', views);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .views-panel.empty {
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
