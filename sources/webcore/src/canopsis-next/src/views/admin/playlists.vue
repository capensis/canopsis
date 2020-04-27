<template lang="pug">
  div
    playlists-list(
      :playlists="playlists",
      :pending="playlistsPending",
      @edit="showEditPlaylistModal",
      @delete="showRemovePlaylistModal",
      @duplicate="showDuplicatePlaylistModal"
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

import rightsTechnicalPlaylistMixin from '@/mixins/rights/technical/playlist';
import entitiesPlaylistMixin from '@/mixins/entities/playlist';

import PlaylistsList from '@/components/other/playlists/playlists-list.vue';
import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';
import { omit } from 'lodash';

export default {
  components: {
    PlaylistsList,
    RefreshBtn,
  },
  mixins: [rightsTechnicalPlaylistMixin, entitiesPlaylistMixin],
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
          action: this.callActionWithFetching(newPlaylist => this.createPlaylist({ data: newPlaylist })),
        },
      });
    },

    showEditPlaylistModal(playlist) {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          title: this.$t('modals.createPlaylist.edit.title'),
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
          action: this.callActionWithFetching(() => this.removePlaylist({ id })),
        },
      });
    },

    showDuplicatePlaylistModal(playlist) {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          title: this.$t('modals.createPlaylist.duplicate.title'),
          playlist: omit(playlist, ['_id']),
          action: this.callActionWithFetching(newPlaylist => this.createPlaylist({ data: newPlaylist })),
        },
      });
    },

    fetchList() {
      this.fetchPlaylistsList();
    },
  },
};
</script>
