import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the playlist Vuex store module
 *
 * @returns {Object} Store module hooks for playlist namespace
 * @property {import('vuex').Store} store - Vuex store instance
 * @property {import('vuex').Module} module - Playlist module instance
 * @property {Function} useGetters - Hook for accessing playlist module getters
 * @property {Function} useActions - Hook for accessing playlist module actions
 */
const usePlaylistStoreModule = () => useStoreModuleHooks('playlist');

/**
 * Hook for managing playlist-related operations and state
 * Replaces the entitiesPlaylistMixin functionality with Composition API
 *
 * @returns {Object} An object containing playlist getters, actions, and methods
 * @property {import('vue').ComputedRef} items - Playlist items
 * @property {import('vue').ComputedRef} pending - Pending state
 * @property {import('vue').ComputedRef} meta - Metadata
 * @property {Function} fetchList - Fetches the list of playlists
 * @property {Function} createPlaylist - Creates a new playlist
 * @property {Function} updatePlaylist - Updates an existing playlist
 * @property {Function} removePlaylist - Removes a playlist
 * @property {Function} fetchItemWithoutStore - Fetches a playlist without storing it
 * @property {Function} createPlaylistWithPopup - Creates playlist with success popup
 * @property {Function} updatePlaylistWithPopup - Updates playlist with success popup
 * @property {Function} removePlaylistWithPopup - Removes playlist with success popup
 */
export const usePlaylist = () => {
  const { t } = useI18n();
  const popups = usePopups();

  const { useGetters, useActions } = usePlaylistStoreModule();

  const getters = useGetters({
    items: 'items',
    pending: 'pending',
    meta: 'meta',
  });

  const actions = useActions({
    fetchList: 'fetchList',
    createPlaylist: 'create',
    updatePlaylist: 'update',
    removePlaylist: 'remove',
    fetchItemWithoutStore: 'fetchItemWithoutStore',
  });

  /**
   * Creates a new playlist and shows a success popup message
   *
   * @param {Object} params - The parameters object
   * @param {Object} params.data - The playlist data to create
   */
  const createPlaylistWithPopup = async ({ data }) => {
    try {
      await actions.createPlaylist({ data });

      popups.success({ text: t('modals.createPlaylist.success.create') });
    } catch (err) {
      popups.error({ text: t('modals.createPlaylist.fail.create') });

      throw err;
    }
  };

  /**
   * Updates an existing playlist and shows a success popup message
   *
   * @param {Object} params - The parameters object
   * @param {string|number} params.id - The ID of the playlist to update
   * @param {Object} params.data - The playlist data to update
   */
  const updatePlaylistWithPopup = async ({ id, data }) => {
    try {
      await actions.updatePlaylist({ id, data });

      popups.success({ text: t('modals.createPlaylist.success.edit') });
    } catch (err) {
      popups.error({ text: t('modals.createPlaylist.fail.edit') });

      throw err;
    }
  };

  /**
   * Removes a playlist and shows a success popup message
   *
   * @param {Object} params - The parameters object
   * @param {string|number} params.id - The ID of the playlist to remove
   */
  const removePlaylistWithPopup = async ({ id }) => {
    try {
      await actions.removePlaylist({ id });

      popups.success({ text: t('modals.createPlaylist.success.delete') });
    } catch (err) {
      popups.error({ text: t('modals.createPlaylist.fail.delete') });

      throw err;
    }
  };

  return {
    ...getters,
    ...actions,
    createPlaylistWithPopup,
    updatePlaylistWithPopup,
    removePlaylistWithPopup,
  };
};
