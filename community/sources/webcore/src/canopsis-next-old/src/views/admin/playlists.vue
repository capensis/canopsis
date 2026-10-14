<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <playlists-list
        :playlists="playlists"
        :pending="playlistsPending"
        :options.sync="options"
        :total-items="playlistsMeta.total_count"
        :creatable="hasCreateAccess"
        :updatable="hasUpdateAccess"
        :removable="hasDeleteAccess"
        @edit="showEditPlaylistModal"
        @remove="showRemovePlaylistModal"
        @duplicate="showDuplicatePlaylistModal"
        @refresh="fetchList"
      />
    </v-card>
    <c-fab-btn
      :has-access="hasCreateAccess"
      @refresh="fetchList"
      @create="showCreatePlaylistModal"
    >
      <span>{{ $t('modals.createPlaylist.create.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions, useAuth } from '@/hooks/auth';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { usePlaylist } from '@/hooks/store/modules/playlist';

import PlaylistsList from '@/components/other/playlists/playlists-list.vue';

export default {
  components: {
    PlaylistsList,
  },
  setup() {
    const {
      items: playlists,
      pending: playlistsPending,
      meta: playlistsMeta,
      fetchList: fetchPlaylistsList,
      createPlaylist,
      updatePlaylist,
      removePlaylist,
    } = usePlaylist();

    const {
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.playlist);

    const { fetchCurrentUser } = useAuth();
    const modals = useModals();
    const { t } = useI18n();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery => fetchPlaylistsList({
        params: convertQueryToRequest(fetchQuery),
      }),
    });

    /**
     * Shows the modal for creating a new playlist.
     * After successful creation, refreshes the current user and playlists list.
     */
    const showCreatePlaylistModal = () => {
      modals.show({
        name: MODALS.createPlaylist,
        config: {
          action: async (newPlaylist) => {
            await createPlaylist({ data: newPlaylist });

            await Promise.all([
              fetchCurrentUser(),
              fetchList(),
            ]);
          },
        },
      });
    };

    /**
     * Shows the modal for editing an existing playlist.
     *
     * @param {Object} playlist - The playlist object to edit.
     * @param {string} playlist._id - The unique identifier of the playlist.
     */
    const showEditPlaylistModal = (playlist) => {
      modals.show({
        name: MODALS.createPlaylist,
        config: {
          playlist,
          title: t('modals.createPlaylist.edit.title'),
          action: async (newPlaylist) => {
            await updatePlaylist({ id: playlist._id, data: newPlaylist });

            await Promise.all([
              fetchCurrentUser(),
              fetchList(),
            ]);
          },
        },
      });
    };

    /**
     * Shows the confirmation modal for deleting a playlist.
     * After successful deletion, refreshes the current user and playlists list.
     *
     * @param {string} id - The unique identifier of the playlist to delete.
     */
    const showRemovePlaylistModal = (id) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removePlaylist({ id });

            await Promise.all([
              fetchCurrentUser(),
              fetchList(),
            ]);
          },
        },
      });
    };

    /**
     * Shows the modal for duplicating an existing playlist.
     * Creates a new playlist based on the provided playlist.
     *
     * @param {Object} playlist - The playlist object to duplicate.
     */
    const showDuplicatePlaylistModal = (playlist) => {
      modals.show({
        name: MODALS.createPlaylist,
        config: {
          title: t('modals.createPlaylist.duplicate.title'),
          playlist: omit(playlist, ['_id']),
          action: async (newPlaylist) => {
            await createPlaylist({ data: newPlaylist });

            await Promise.all([
              fetchCurrentUser(),
              fetchList(),
            ]);
          },
        },
      });
    };

    onMounted(fetchList);

    return {
      playlists,
      playlistsPending,
      playlistsMeta,
      options,
      updateOptions,
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
      fetchList,
      showCreatePlaylistModal,
      showEditPlaylistModal,
      showRemovePlaylistModal,
      showDuplicatePlaylistModal,
    };
  },
};
</script>
