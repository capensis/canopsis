<template>
  <div class="groups-wrapper bg-secondary">
    <div class="groups-wrapper__items">
      <groups-top-bar-group v-for="group in availableGroups" :key="group._id" :group="group" />
      <groups-top-bar-playlists />
    </div>
    <groups-settings-button
      top-bar
      tooltip-left
      @toggle-editing-mode="toggleNavigationEditingMode"
    />
  </div>
</template>

<script>
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

  mixins: [entitiesViewGroupMixin, layoutNavigationGroupsBarMixin],
};
</script>

<style lang="scss" scoped>
.groups-wrapper {
  display: flex;
  align-items: center;
  flex: 1 1 100%;
  height: 48px;
  min-width: 0;
  position: relative;
  width: 100%;
  padding: 0 56px 0 4px;

  &__items {
    display: flex;
    align-items: center;
    min-width: 0;
    overflow-x: auto;
    scrollbar-width: none;

    &::-webkit-scrollbar {
      display: none;
    }
  }

  :deep(.groups-settings-button) {
    position: absolute;
    right: 16px;
    top: 50%;
    bottom: auto;
    transform: translateY(-50%);
  }
}
</style>
