<template lang="pug">
  v-menu.group-item.groups-top-bar-playlist(
    v-if="availablePlaylists.length",
    content-class="group-v-menu-content secondary",
    close-delay="0",
    open-on-hover,
    offset-y,
    bottom,
    dark
  )
    template(#activator="{ on }")
      v-btn.groups-top-bar-playlist__dropdown-btn(v-on="on", color="secondary lighten-1")
        span {{ $t(`pageHeaders.${$constants.USERS_PERMISSIONS.technical.playlist}.title`) }}
        v-icon.ml-0(right, dark) arrow_drop_down
    v-list
      v-list-tile(
        v-for="playlist in availablePlaylists",
        :key="playlist._id",
        :to="{ name: $constants.ROUTES_NAMES.playlist, params: { id: playlist._id } }"
      )
        v-list-tile-title
          span {{ playlist.name }}
</template>

<script>
import { playlistSchema } from '@/store/schemas';

import layoutNavigationGroupsBarPlaylistsMixin from '@/mixins/layout/navigation/groups-bar-playlists';
import { registrableMixin } from '@/mixins/registrable';

export default {
  mixins: [
    layoutNavigationGroupsBarPlaylistsMixin,

    registrableMixin([playlistSchema], 'playlists'),
  ],
};
</script>

<style lang="scss" scoped>
.groups-top-bar-playlist {
  &__dropdown-btn {
    text-transform: none;
  }
}
</style>
