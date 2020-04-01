<template lang="pug">
  draggable.groups-panel.secondary(
    :value="groups",
    :class="{ empty: isGroupsEmpty }",
    :component-data="{ props: { expand: true, dark: true, focusable: true } }",
    :options="draggableOptions",
    element="v-expansion-panel",
    @change="changeGroupsOrdering"
  )
    group-panel(
      v-for="group in groups",
      :group="group",
      :key="group._id"
    )
      draggable-group-views(v-model="group.views")
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { getDuplicateCountItems } from '@/helpers/searching';

import DraggableGroupViews from './draggable-group-views.vue';
import GroupPanel from './group-panel.vue';

export default {
  components: { DraggableGroupViews, Draggable, GroupPanel },
  model: {
    prop: 'groups',
    event: 'change',
  },
  props: {
    groups: {
      type: Array,
      required: true,
    },
  },
  computed: {
    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        group: 'groups',
      };
    },

    isGroupsEmpty() {
      return this.groups.length === 0;
    },
  },
  methods: {
    changeGroupsOrdering({ moved, added, removed }) {
      const groups = [...this.groups];

      if (moved) {
        const [item] = groups.splice(moved.oldIndex, 1);

        groups.splice(moved.newIndex, 0, item);
      } else if (added) {
        const duplicateGroupCount = getDuplicateCountItems(groups, added.element);

        const group = duplicateGroupCount !== 0
          ? { ...added.element, name: `${added.element.name} (${duplicateGroupCount})` }
          : added.element;

        groups.splice(added.newIndex, 0, group);
      } else if (removed) {
        groups.splice(removed.oldIndex, 1);
      }

      if (groups) {
        this.$emit('change', groups);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .groups-panel.empty {
    &:after {
      content: '';
      display: block;
      height: 48px;
      width: 100%;
      border: 4px dashed #3c5365;
      border-radius: 5px;
      position: relative;
    }
  }
</style>

