<template>
  <div>
    <c-page-header />

    <v-card class="ma-4 mt-0">
      <playlists-list
        :playlists="playlists"
        :pending="playlistsPending"
        :options="options"
        :total-items="playlistsMeta.total_count"
        :updatable="hasUpdateAccess"
        :removable="hasDeleteAccess"
        :duplicable="hasCreateAccess"
        @edit="showEditPlaylistModal"
        @duplicate="showDuplicatePlaylistModal"
        @remove="showRemovePlaylistModal"
        @refresh="fetchList"
        @update:options="updateOptions"
      />
    </v-card>

    <c-fab-btn
      :has-access="hasCreateAccess"
      @create="showCreatePlaylistModal"
      @refresh="fetchList"
    />
  </div>
</template>

<script>
import { onMounted } from 'vue';

import { USER_PERMISSIONS, MODALS } from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { convertQueryToRequest } from '@/helpers/query';

import { useCRUDPermissions, useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { usePlaylist } from '@/hooks/store/modules/playlist';

import PlaylistsList from '@/components/other/playlists/playlists-list.vue';

export default {
  components: {
    PlaylistsList,
  },

  setup() {
    const {
      playlists,
      playlistsPending,
      playlistsMeta,
      fetchPlaylistsList,
      createPlaylist,
      updatePlaylist,
      removePlaylist,
    } = usePlaylist();

    const { hasCreateAccess, hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.playlist,
    );

    const modals = useModals();
    const { fetchCurrentUser } = useAuth();
    const { t } = useI18n();

    const {
      options,
      handler: fetchList,
      updateOptions,
    } = useLocalQueryWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        sortBy: [],
        sortDesc: [],
        search: '',
      },
      onUpdate: fetchQuery =>
        fetchPlaylistsList({
          params: convertQueryToRequest(fetchQuery),
        }),
    });

    const showRemovePlaylistModal = id =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('modals.confirmation.text'),
          action: async () => {
            await removePlaylist({ id });

            await Promise.all([fetchCurrentUser(), fetchList()]);
          },
        },
      });

    const refreshPlaylists = async () => {
      await Promise.all([fetchCurrentUser(), fetchList()]);
    };

    const showCreatePlaylistModal = () =>
      modals.show({
        name: MODALS.createPlaylist,
        config: {
          action: async (data) => {
            await createPlaylist({ data });
            await refreshPlaylists();
          },
        },
      });

    const showEditPlaylistModal = playlist =>
      modals.show({
        name: MODALS.createPlaylist,
        config: {
          playlist,
          action: async (data) => {
            await updatePlaylist({ id: playlist._id, data });
            await refreshPlaylists();
          },
        },
      });

    const showDuplicatePlaylistModal = playlist =>
      modals.show({
        name: MODALS.createPlaylist,
        config: {
          playlist,
          duplicating: true,
          action: async (data) => {
            await createPlaylist({ data });
            await refreshPlaylists();
          },
        },
      });

    onMounted(fetchList);

    return {
      playlists,
      playlistsPending,
      playlistsMeta,
      options,
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
      fetchList,
      updateOptions,
      showRemovePlaylistModal,
      showCreatePlaylistModal,
      showEditPlaylistModal,
      showDuplicatePlaylistModal,
    };
  },
};
</script>
