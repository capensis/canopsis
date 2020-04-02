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
      draggable-group-views(
        v-model="group.views",
        :prepareView="mapViewEntity",
        :put="viewPut",
        :pull="viewPull"
      )
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { getDuplicateEntityName } from '@/helpers/entities';
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

    allViewsList() {
      return this.groups.reduce((ids, { views }) => {
        ids.push(...views);

        return ids;
      }, []);
    },
  },
  methods: {
    mapViewEntity(view) {
      return { ...view, name: getDuplicateEntityName(view, this.allViewsList) };
    },

    mapGroupEntity(group) {
      return {
        ...group,
        name: getDuplicateEntityName(group, this.groups),
        views: group.views.map(this.mapViewEntity),
      };
    },

    changeGroupsOrdering(event) {
      this.$emit('change', dragDropChangePositionHandler(this.groups, event, this.mapGroupEntity));
    },
  },
};
</script>

<style lang="scss" scoped>
  .groups-panel {
    cursor: move;

    & /deep/ .v-expansion-panel__header {
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

