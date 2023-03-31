<template lang="pug">
  v-expansion-panel.sidebar-playlists(v-if="availablePlaylists.length", expand, focusable, dark)
    v-expansion-panel-content.secondary.lighten-1.white--text
      template(#header="")
        div.panel-header
          span {{ $t(`pageHeaders.${$constants.USERS_PERMISSIONS.technical.playlist}.title`) }}
      router-link(
        v-for="playlist in availablePlaylists",
        :key="playlist._id",
        :title="playlist.name",
        :to="{ name: $constants.ROUTES_NAMES.playlist, params: { id: playlist._id } }"
      )
        v-card.secondary.lighten-2.sidebar-playlists__button
          v-card-text
            v-layout(align-center, justify-space-between)
              v-flex
                v-layout(align-center)
                  span.ellipsis.pl-3 {{ playlist.name }}
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
  .sidebar-playlists {
    padding: 10px;
    box-shadow: none;

    &__button {
      border-radius: 0;
    }
  }

  a {
    color: inherit;
    text-decoration: none;

  &.router-link-active ::v-deep .v-card {
     background-color: #73879a !important;
     border-color: #73879a !important;
   }
  }
</style>
