<template>
  <v-layout class="secondary groups-wrapper">
    <v-tabs
      ref="tabs"
      background-color="secondary"
      show-arrows
      dark
    >
      <groups-top-bar-group
        v-for="group in availableGroups"
        :key="group._id"
        :group="group"
      />
      <groups-top-bar-playlists />
    </v-tabs>
    <groups-settings-button
      tooltip-left
      :wrapper-props="{ direction: 'bottom', absolute: true, right: true, bottom: true }"
      :button-props="{ fab: true, dark: true, small: true }"
      @toggleEditingMode="toggleEditingMode"
    />
  </v-layout>
</template>

<script>
import { vuetifyTabsMixin } from '@/mixins/vuetify/tabs';
import { entitiesViewGroupMixin } from '@/mixins/entities/view/group';
import { layoutNavigationGroupsBarMixin } from '@/mixins/layout/navigation/groups-bar';

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

      .v-list-item__title {
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
