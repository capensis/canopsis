<template lang="pug">
  v-expansion-panel(expand, focusable, dark)
    v-expansion-panel-content.secondary.lighten-1.white--text
      div.panel-header(slot="header")
        span {{ $t('common.playlists') }}
      v-card(
        v-for="playlist in availablePlaylists",
        :key="playlist._id"
      )
        router-link(
          :title="playlist.name",
          :to="{ name: 'playlist', params: { id: playlist._id } }"
        )
          v-card-text.secondary.lighten-2
            v-layout(align-center, justify-space-between)
              v-flex
                v-layout(align-center)
                  span.pl-3 {{ playlist.name }}
</template>

<script>
import authMixin from '@/mixins/auth';

import { playlistSchema } from '@/store/schemas';

import entitiesPlaylistMixin from '@/mixins/entities/playlist';
import registrableMixin from '@/mixins/registrable';

export default {
  mixins: [
    authMixin,
    entitiesPlaylistMixin,
    registrableMixin([playlistSchema], 'playlists'),
  ],
  computed: {
    availablePlaylists() {
      return this.playlists.filter(playlist => this.checkReadAccess(playlist._id));
    },
  },
  mounted() {
    this.fetchPlaylistsList();
  },
};
</script>

<style scoped>
  a {
    color: inherit;
    text-decoration: none;
  }
</style>
