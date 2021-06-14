import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('playlist');

export default {
  computed: {
    ...mapGetters({
      playlists: 'items',
      playlistsPending: 'pending',
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
