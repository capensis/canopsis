<template lang="pug">
  div
    playlists-list(
      :playlists="playlists",
      :pending="playlistsPending",
      @edit="showEditPlaylistModal",
      @delete="showRemovePlaylistModal"
    )
    .fab(v-if="hasCreateAnyPlaylistAccess")
      v-layout(column)
        refresh-btn(@click="fetchList")
        v-tooltip(left)
          v-btn(
            slot="activator",
            color="primary",
            data-test="addButton",
            fab,
            @click.stop="showCreatePlaylistModal"
          )
            v-icon add
          span {{ $t('modals.createPlaylist.create.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import authMixin from '@/mixins/auth';
import rightsTechnicalPlaylistMixin from '@/mixins/rights/technical/playlist';
import entitiesPlaylistMixin from '@/mixins/entities/playlist';
import entitiesPlaylistRightMixin from '@/mixins/entities/playlist/right';

import PlaylistsList from '@/components/other/playlists/admin/playlists-list.vue';
import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';

export default {
  components: {
    PlaylistsList,
    RefreshBtn,
  },
  mixins: [
    authMixin,
    rightsTechnicalPlaylistMixin,
    entitiesPlaylistMixin,
    entitiesPlaylistRightMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    /**
     * Return function for calling of the action with popups and fetching
     *
     * @param {Function} action
     * @returns {Promise<void>}
     */
    callActionWithFetching(action) {
      return async (...args) => {
        try {
          await action(...args);

          this.fetchList();

          this.$popups.success({ text: this.$t('success.default') });
        } catch (err) {
          this.$popups.error({ text: this.$t('errors.default') });
        }
      };
    },

    showCreatePlaylistModal() {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          action: this.callActionWithFetching(async (newPlaylist) => {
            const { _id: playlistId } = await this.createPlaylist({ data: newPlaylist });

            return this.createRightByPlaylistId(playlistId);
          }),
        },
      });
    },

    showEditPlaylistModal(playlist) {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          playlist,

          action:
            this.callActionWithFetching(newPlaylist => this.updatePlaylist({ id: playlist._id, data: newPlaylist })),
        },
      });
    },

    showRemovePlaylistModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: this.callActionWithFetching(async () => {
            await this.removePlaylist({ id });

            return this.removeRightByPlaylistId(id);
          }),
        },
      });
    },

    fetchList() {
      this.fetchPlaylistsList();
    },
  },
};
</script>
