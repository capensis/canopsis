import { mapActions, mapState, mapStores } from 'pinia';

import { usePlaylistStore } from '@/stores/modules/playlist';

export const entitiesPlaylistMixin = {
  computed: {
    ...mapStores(usePlaylistStore),
    ...mapState(usePlaylistStore, {
      playlists: 'items',
      playlistsPending: 'pending',
      playlistsMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions(usePlaylistStore, {
      fetchPlaylistsList: 'fetchList',
      createPlaylist: 'create',
      updatePlaylist: 'update',
      removePlaylist: 'remove',
      bulkEnablePlaylists: 'bulkEnable',
      bulkDisablePlaylists: 'bulkDisable',
      bulkRemovePlaylists: 'bulkRemove',
    }),
  },
};
