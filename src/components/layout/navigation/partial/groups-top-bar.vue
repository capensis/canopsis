<template lang="pug">
  v-layout.secondary.groups-wrapper
    v-tabs(ref="tabs", color="secondary", show-arrows, dark)
      v-menu(
      v-for="group in groups",
      :key="group._id",
      content-class="group-v-menu-content secondary",
      close-delay="0",
      open-on-hover,
      offset-y,
      bottom,
      dark
      )
        v-btn(slot="activator", flat, dark)
          span {{ group.name }}
          v-btn(
          v-show="isEditingMode",
          depressed,
          small,
          icon,
          @click.stop="showEditGroupModal(group)"
          )
            v-icon(small) edit
          v-icon(dark) arrow_drop_down
        v-list
          v-list-tile(
          v-for="view in getAvailableViewsForGroup(group)",
          :key="view._id",
          :to="{ name: 'view', params: { id: view._id } }",
          )
            v-list-tile-title
              span {{ view.title }}
              v-btn.edit-view-button(
              v-show="(checkUpdateViewAccessById(view._id) || checkDeleteViewAccessById(view._id)) && isEditingMode",
              color="grey darken-2",
              depressed,
              small,
              icon,
              @click.prevent="showEditViewModal(view)"
              )
                v-icon(small) edit
              v-btn.duplicate-view-button(
              v-show="isEditingMode",
              depressed,
              small,
              icon,
              color="grey darken-2",
              @click.prevent="showDuplicateViewModal(view)"
              )
                v-icon(small) file_copy
    groups-settings-button(
    :isEditingMode="isEditingMode",
    @toggleEditingMode="toggleEditingMode",
    :wrapperProps="{ direction: 'bottom', absolute: true, right: true, bottom: true }",
    :buttonProps="{ fab: true, dark: true, small: true }"
    )
</template>

<script>
import entitiesViewGroupMixin from '@/mixins/entities/view/group/index';
import layoutNavigationGroupMenuMixin from '@/mixins/layout/navigation/group-menu';

import GroupsSettingsButton from './groups-settings-button.vue';

export default {
  components: { GroupsSettingsButton },
  mixins: [entitiesViewGroupMixin, layoutNavigationGroupMenuMixin],
  watch: {
    groups() {
      this.$nextTick(this.callTabsOnResizeMethod);
    },
  },
  methods: {
    toggleEditingMode() {
      this.isEditingMode = !this.isEditingMode;

      this.$nextTick(this.callTabsOnResizeMethod);
    },

    callTabsOnResizeMethod() {
      this.$refs.tabs.onResize();
    },
  },
};
</script>

<style lang="scss">
  .groups-wrapper {
    height: 48px;

    .v-menu__activator .v-btn {
      text-transform: none;
    }

    .v-speed-dial--bottom.v-speed-dial--absolute {
      bottom: -10px;
    }

    .v-speed-dial--right.v-speed-dial--absolute {
      right: 25px;
    }
  }

  .group-v-menu-content {
    .v-list {
      background-color: inherit;

      .v-list__tile__title {
        height: 28px;
        line-height: 28px;
      }

      .edit-view-button, .duplicate-view-button {
        vertical-align: top;
        margin: 0 0 0 8px;
      }
    }
  }
</style>
