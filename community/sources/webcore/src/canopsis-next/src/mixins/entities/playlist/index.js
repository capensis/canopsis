import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('playlist');

export const entitiesPlaylistMixin = {
  computed: {
    ...mapGetters({
      playlists: 'items',
      playlistsPending: 'pending',
      playlistsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchPlaylistsList: 'fetchList',
      createPlaylist: 'create',
      updatePlaylist: 'update',
      removePlaylist: 'remove',
    }),
  },
};
