<template>
  <v-menu
    class="group-item groups-top-bar-playlist"
    v-if="availablePlaylists.length"
    content-class="group-v-menu-content secondary"
    close-delay="0"
    open-on-hover
    offset-y
    bottom
    dark
  >
    <template #activator="{ on }">
      <v-btn
        class="groups-top-bar-playlist__dropdown-btn"
        v-on="on"
        color="secondary lighten-1"
      >
        <span>{{ $t(`pageHeaders.${$constants.USERS_PERMISSIONS.technical.playlist}.title`) }}</span>
        <v-icon
          class="ml-0"
          right
          dark
        >
          arrow_drop_down
        </v-icon>
      </v-btn>
    </template>
    <v-list>
      <v-list-item
        v-for="playlist in availablePlaylists"
        :key="playlist._id"
        :to="{ name: $constants.ROUTES_NAMES.playlist, params: { id: playlist._id } }"
      >
        <v-list-item-title>
          <span>{{ playlist.name }}</span>
        </v-list-item-title>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import layoutNavigationGroupsBarPlaylistsMixin from '@/mixins/layout/navigation/groups-bar-playlists';

export default {
  mixins: [
    layoutNavigationGroupsBarPlaylistsMixin,
  ],
};
</script>

<style lang="scss" scoped>
.groups-top-bar-playlist {
  &__dropdown-btn {
    margin: 6px 0;
    text-transform: none;
  }
}
</style>
