<template lang="pug">
  v-expansion-panel-content.secondary.white--text.group-item(:data-test="`panel-group-${group._id}`")
    div.panel-header(slot="header")
      span(:data-test="`groupsSideBar-group-${group._id}`") {{ group.name }}
      v-btn(
        :data-test="`editGroupButton-group-${group._id}`",
        v-show="isEditingMode",
        depressed,
        small,
        icon,
        @click.stop="showEditGroupModal(group)"
      )
        v-icon(small) edit
    draggable.panel(
      :value="group.views",
      :options="draggableOptions",
      @update="updateViewsOrdering"
    )
      groups-side-bar-group-view(
        v-for="view in group.views",
        :key="view._id",
        :view="view",
        :isEditingMode="isEditingMode"
      )
</template>

<script>
import Draggable from 'vuedraggable';
import arrayMove from 'array-move';

import GroupsSideBarGroupView from './groups-side-bar-group-view.vue';

export default {
  components: { Draggable, GroupsSideBarGroupView },
  props: {
    group: {
      type: Object,
      required: true,
    },
    draggableOptions: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    updateViewsOrdering({ oldIndex, newIndex }) {
      this.$emit('update:group', {
        ...this.group,

        views: arrayMove(this.group.views, oldIndex, newIndex),
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .group-item /deep/ .v-expansion-panel__header {
    height: 48px;
  }

  .panel-header {
    max-width: 88%;

    span {
      max-width: 100%;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      display: inline-block;
      vertical-align: middle;

      .editing & {
        max-width: 73%;
      }
    }
  }
</style>
