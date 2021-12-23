<template lang="pug">
  v-menu.group-item(
    v-if="availablePlaylists.length",
    content-class="group-v-menu-content secondary",
    close-delay="0",
    open-on-hover,
    offset-y,
    bottom,
    dark
  )
    div.v-btn.v-btn--flat.theme--dark.secondary.lighten-1(slot="activator")
      span {{ $t(`pageHeaders.${$constants.USERS_PERMISSIONS.technical.playlist}.title`) }}
      v-icon(dark) arrow_drop_down
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
  a {
    color: inherit;
    text-decoration: none;
  }
</style>
