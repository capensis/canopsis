<template lang="pug">
  group-panel(
    :group="group",
    :orderChanged="isGroupsOrderChanged",
    :isEditing="isNavigationEditingMode",
    @change="showEditGroupModal"
  )
    draggable.views-panel.secondary.lighten-1(
      :class="{ empty: isGroupEmpty }",
      :value="group.views",
      :options="draggableOptions",
      @change="changeViewsOrdering"
    )
      groups-side-bar-group-view(
        v-for="view in group.views",
        :key="view._id",
        :view="view",
        :isGroupsOrderChanged="isGroupsOrderChanged"
      )
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { dragDropChangePositionHandler } from '@/helpers/dragdrop';

import layoutNavigationGroupsBarGroupMixin from '@/mixins/layout/navigation/groups-bar-group';

import GroupsSideBarGroupView from './groups-side-bar-group-view.vue';
import GroupPanel from './group-panel.vue';

export default {
  components: { GroupPanel, Draggable, GroupsSideBarGroupView },
  mixins: [layoutNavigationGroupsBarGroupMixin],
  props: {
    isGroupsOrderChanged: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    draggableOptions() {
      return {
        disabled: !this.isNavigationEditingMode,
        animation: VUETIFY_ANIMATION_DELAY,
        group: 'views',
      };
    },

    isGroupEmpty() {
      return this.group.views && this.group.views.length === 0;
    },
  },
  methods: {
    changeViewsOrdering(event) {
      this.$emit('update:group', {
        ...this.group,
        views: dragDropChangePositionHandler(this.group.views, event),
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .group-item {
    & ::v-deep .v-expansion-panel__header {
      height: 48px;
    }

    &.editing {
      & ::v-deep .v-expansion-panel__header {
        cursor: move;
      }

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
    }

    & ::v-deep .v-expansion-panel__body  .v-card {
      border-radius: 0;
      box-shadow: 0 0 0 0 rgba(0,0,0,.2),0 0 0 0 rgba(0,0,0,.14),0 0 0 0 rgba(0,0,0,.12)!important;
    }
  }
</style>
