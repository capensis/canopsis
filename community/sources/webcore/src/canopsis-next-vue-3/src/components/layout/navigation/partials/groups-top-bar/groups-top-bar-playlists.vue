<template>
  <v-menu
    v-if="availablePlaylists.length"
    class="group-item groups-top-bar-playlist"
    content-class="group-v-menu-content bg-secondary"
    location="bottom"
    open-on-hover
  >
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        class="groups-top-bar-playlist__dropdown-btn"
        color="secondary-lighten-1"
      >
        <span class="text-none">{{ $t(`pageHeaders.${USER_PERMISSIONS.technical.playlist}.title`) }}</span>
        <c-icon class="ml-1">arrow_drop_down</c-icon>
      </v-btn>
    </template>
    <v-list class="bg-secondary">
      <v-list-item
        v-for="playlist in availablePlaylists"
        :key="playlist._id"
        :to="{ name: ROUTES_NAMES.playlist, params: { id: playlist._id } }"
      >
        <v-list-item-title>{{ playlist.name }}</v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { ROUTES_NAMES, USER_PERMISSIONS } from '@/constants';

import layoutNavigationGroupsBarPlaylistsMixin from '@/mixins/layout/navigation/groups-bar-playlists';

export default {
  mixins: [layoutNavigationGroupsBarPlaylistsMixin],

  setup() {
    return { ROUTES_NAMES, USER_PERMISSIONS };
  },
};
</script>

<style lang="scss" scoped>
.groups-top-bar-playlist {
  flex: 0 0 auto;

  &__dropdown-btn {
    margin: 6px 0;
    text-transform: none;
  }
}
</style>
