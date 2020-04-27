<template lang="pug">
  div
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.broadcastMessages') }}
    div.white
      v-data-table(
        :headers="headers",
        :items="playlists",
        :loading="playlistsPending",
        :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
        item-key="_id"
      )
        template(slot="items", slot-scope="props")
          tr(:data-test="`role-${props.item._id}`")
            td {{ props.item.name }}
            td
              enabled-column(:value="props.item.fullscreen")
            td
              enabled-column(:value="props.item.enabled")
            td {{ props.item.interval | interval }}
            td
              v-btn.ma-0(
                v-if="hasUpdateAnyPlaylistAccess",
                data-test="editButton",
                icon,
                @click="showEditPlaylistModal(props.item)"
              )
                v-icon edit
              v-btn.ma-0(
                v-if="hasDeleteAnyPlaylistAccess",
                data-test="deleteButton",
                icon,
                @click="showRemovePlaylistModal(props.item._id)"
              )
                v-icon(color="error") delete
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
          span {{ $t('modals.createBroadcastMessage.create.title') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import moment from 'moment';

import { MODALS, BROADCAST_MESSAGES_STATUSES } from '@/constants';

import rightsTechnicalPlaylistMixin from '@/mixins/rights/technical/playlist';

import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';
import SearchField from '@/components/forms/fields/search-field.vue';
import EnabledColumn from '@/components/tables/enabled-column.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('playlist');

export default {
  filters: {
    interval(interval) {
      return interval && interval.interval && `${interval.interval}${interval.unit}`;
    },
  },
  components: {
    RefreshBtn,
    SearchField,
    EnabledColumn,
  },
  mixins: [rightsTechnicalPlaylistMixin],
  computed: {
    ...mapGetters({
      playlists: 'items',
      playlistsPending: 'pending',
    }),

    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.fullscreen'),
          value: 'fullscreen',
        },
        {
          text: this.$t('common.enabled'),
          value: 'enabled',
        },
        {
          text: this.$t('common.interval'),
          value: 'interval',
        },
        {
          text: this.$t('common.actionsLabel'),
          sortable: false,
        },
      ];
    },

    preparedBroadcastMessages() {
      return this.broadcastMessages.map((message) => {
        const now = moment().unix();
        let status = BROADCAST_MESSAGES_STATUSES.pending;

        if (now >= message.start) {
          if (now <= message.end) {
            status = BROADCAST_MESSAGES_STATUSES.active;
          } else {
            status = BROADCAST_MESSAGES_STATUSES.expired;
          }
        }

        return {
          ...message,

          status,
        };
      });
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchPlaylistsList: 'fetchList',
      createPlaylist: 'create',
      updatePlaylist: 'update',
      removePlaylist: 'remove',
    }),

    /**
     * Function for calling of the action with popups and fetching
     *
     * @param {Function} action
     * @returns {Promise<void>}
     */
    async callActionWithFetching(action) {
      try {
        await action();

        this.fetchList();

        this.$popups.success({ text: this.$t('success.default') });
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    showCreatePlaylistModal() {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          action: newPlaylist =>
            this.callActionWithFetching(() => this.createPlaylist({ data: newPlaylist })),
        },
      });
    },

    showEditPlaylistModal(playlist) {
      this.$modals.show({
        name: MODALS.createPlaylist,
        config: {
          playlist,

          action: newPlaylist =>
            this.callActionWithFetching(() => this.updatePlaylist({ id: playlist._id, data: newPlaylist })),
        },
      });
    },

    showRemovePlaylistModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.callActionWithFetching(() => this.removePlaylist({ id })),
        },
      });
    },

    fetchList() {
      this.fetchPlaylistsList();
    },
  },
};
</script>
