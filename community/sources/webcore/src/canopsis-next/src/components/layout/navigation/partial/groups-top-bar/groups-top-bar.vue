<template lang="pug">
  v-layout.secondary.groups-wrapper
    v-tabs(ref="tabs", color="secondary", show-arrows, dark)
      template(v-if="hasReadAnyViewAccess")
        groups-top-bar-group(
          v-for="group in availableGroups",
          :key="group._id",
          :group="group"
        )
      groups-top-bar-playlists
    groups-settings-button(
      tooltipLeft,
      :wrapperProps="{ direction: 'bottom', absolute: true, right: true, bottom: true }",
      :buttonProps="{ fab: true, dark: true, small: true }",
      @toggleEditingMode="toggleEditingMode"
    )
</template>

<script>
import { groupSchema } from '@/store/schemas';

import { vuetifyTabsMixin } from '@/mixins/vuetify/tabs';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import layoutNavigationGroupsBarMixin from '@/mixins/layout/navigation/groups-bar';
import { registrableMixin } from '@/mixins/registrable';

import GroupsSettingsButton from '../groups-settings-button.vue';

import GroupsTopBarGroup from './groups-top-bar-group.vue';
import GroupsTopBarPlaylists from './groups-top-bar-playlists.vue';

export default {
  components: {
    GroupsSettingsButton,
    GroupsTopBarGroup,
    GroupsTopBarPlaylists,
  },
  mixins: [
    vuetifyTabsMixin,
    entitiesViewGroupMixin,
    layoutNavigationGroupsBarMixin,

    registrableMixin([groupSchema], 'groups'),
  ],
  watch: {
    filledGroups() {
      this.$nextTick(this.callTabsOnResizeMethod);
    },
  },
  methods: {
    toggleEditingMode() {
      this.toggleNavigationEditingMode();

      this.$nextTick(this.callTabsOnResizeMethod);
    },
  },
};
</script>

<style lang="scss">
  .groups-wrapper {
    height: 48px;

    .v-speed-dial--absolute {
      &.v-speed-dial--bottom {
        bottom: -10px;
      }

      &.v-speed-dial--right {
        right: 25px;
      }
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

  .group-item .v-menu__activator .v-btn {
    text-transform: none;
  }
</style>
