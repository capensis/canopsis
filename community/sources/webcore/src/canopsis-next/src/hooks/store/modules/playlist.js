import { storeToRefs } from 'pinia';

import { usePlaylistStore } from '@/stores/modules/playlist';

/**
 * Playlist store facade.
 *
 * @returns {Object}
 */
export const usePlaylist = () => {
  const playlistStore = usePlaylistStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(playlistStore);

  return {
    playlists: items,
    playlistsMeta: meta,
    playlistsPending: pending,
    playlistsFetchingParams: fetchingParams,

    fetchPlaylistsList: playlistStore.fetchList,
    fetchPlaylistsListWithPreviousParams: playlistStore.fetchListWithPreviousParams,
    fetchPlaylistWithoutStore: playlistStore.fetchItemWithoutStore,
    createPlaylist: playlistStore.create,
    updatePlaylist: playlistStore.update,
    removePlaylist: playlistStore.remove,
    bulkEnablePlaylist: playlistStore.bulkEnable,
    bulkDisablePlaylist: playlistStore.bulkDisable,
    bulkRemovePlaylist: playlistStore.bulkRemove,
  };
};
