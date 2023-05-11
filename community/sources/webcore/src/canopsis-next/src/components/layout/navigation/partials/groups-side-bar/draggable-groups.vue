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
      v-for="(group, groupIndex) in groups",
      :group="group",
      :key="group._id"
    )
      draggable-group-views(
        :views="group.views",
        :put="viewPut",
        :pull="viewPull",
        @change="changeViewsHandler(groupIndex, $event)"
      )
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

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
    put: {
      type: Boolean,
      default: false,
    },
    pull: {
      type: [Boolean, String],
      default: false,
    },
    viewPut: {
      type: Boolean,
      default: false,
    },
    viewPull: {
      type: [Boolean, String],
      default: false,
    },
  },
  computed: {
    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
        group: { name: 'groups', put: this.put, pull: this.pull },
      };
    },

    isGroupsEmpty() {
      return this.groups.length === 0;
    },
  },
  methods: {
    changeGroupsOrdering(event) {
      this.$emit('change', dragDropChangePositionHandler(this.groups, event));
    },

    changeViewsHandler(groupIndex, views) {
      const group = this.groups[groupIndex];

      this.$emit('change:group', groupIndex, { ...group, views });
    },
  },
};
</script>

<style lang="scss" scoped>
  .groups-panel {
    cursor: move;

    & ::v-deep .v-expansion-panel__header {
      cursor: move;
    }

    &.empty {
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
  }
</style>
