import authMixin from '@/mixins/auth';

import entitiesPlaylistMixin from '@/mixins/entities/playlist';

export default {
  mixins: [
    authMixin,
    entitiesPlaylistMixin,
  ],
  computed: {
    availablePlaylists() {
      return this.playlists.filter(playlist => this.checkAccess(playlist._id) && playlist.enabled);
    },
  },
  mounted() {
    this.fetchPlaylistsList();
  },
};
