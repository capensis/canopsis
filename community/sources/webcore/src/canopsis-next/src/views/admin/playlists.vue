<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <playlists-list
        :playlists="playlists"
        :pending="playlistsPending"
        :options.sync="options"
        :total-items="playlistsMeta.total_count"
        :creatable="hasCreateAnyPlaylistAccess"
        :updatable="hasUpdateAnyPlaylistAccess"
        :removable="hasDeleteAnyPlaylistAccess"
        @edit="showEditPlaylistModal"
        @remove="showRemovePlaylistModal"
        @duplicate="showDuplicatePlaylistModal"
      />
    </v-card>
    <c-fab-btn
      v-if="hasCreateAnyPlaylistAccess"
      @refresh="fetchList"
      @create="showCreatePlaylistModal"
    >
      <span>{{ $t('modals.createPlaylist.create.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { permissionsTechnicalPlaylistMixin } from '@/mixins/permissions/technical/playlist';
import { entitiesPlaylistMixin } from '@/mixins/entities/playlist';
import { localQueryMixin } from '@/mixins/query-local/query';

import PlaylistsList from '@/components/other/playlists/playlists-list.vue';

export default {
  components: {
    PlaylistsList,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    permissionsTechnicalPlaylistMixin,
    entitiesPlaylistMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    showCreatePlaylistModal() {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          action: async (newPlaylist) => {
            await this.createPlaylist({ data: newPlaylist });

            await Promise.all([
              this.fetchCurrentUser(),
              this.fetchList(),
            ]);
          },
        },
      });
    },

    showEditPlaylistModal(playlist) {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          playlist,

          title: this.$t('modals.createPlaylist.edit.title'),
          action: async (newPlaylist) => {
            await this.updatePlaylist({ id: playlist._id, data: newPlaylist });

            await Promise.all([
              this.fetchCurrentUser(),
              this.fetchList(),
            ]);
          },
        },
      });
    },

    showRemovePlaylistModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removePlaylist({ id });

            await Promise.all([
              this.fetchCurrentUser(),
              this.fetchList(),
            ]);
          },
        },
      });
    },

    showDuplicatePlaylistModal(playlist) {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          title: this.$t('modals.createPlaylist.duplicate.title'),
          playlist: omit(playlist, ['_id']),
          action: async (newPlaylist) => {
            await this.createPlaylist({ data: newPlaylist });

            await Promise.all([
              this.fetchCurrentUser(),
              this.fetchList(),
            ]);
          },
        },
      });
    },

    fetchList() {
      return this.fetchPlaylistsList({ params: this.getQuery() });
    },
  },
};
</script>
