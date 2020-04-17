<template lang="pug">
  div
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.broadcastMessages') }}
    div.white
      v-data-table(
        :headers="headers",
        :items="preparedBroadcastMessages",
        :loading="broadcastMessagesPending",
        :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
        item-key="_id"
      )
        template(slot="items", slot-scope="props")
          tr(:data-test="`role-${props.item._id}`")
            td {{ $t(`tables.broadcastMessages.statuses.${props.item.status}`) }}
            td.broadcast-message-cell
              broadcast-message(:message="props.item.message", :color="props.item.color")
            td {{ props.item.start | date('long', true) }}
            td {{ props.item.end | date('long', true) }}
            td
              v-btn.ma-0(
                v-if="hasUpdateAnyBroadcastMessageAccess",
                data-test="editButton",
                icon,
                @click="showEditBroadcastMessageModal(props.item)"
              )
                v-icon edit
              v-btn.ma-0(
                v-if="hasDeleteAnyBroadcastMessageAccess",
                data-test="deleteButton",
                icon,
                @click="showRemoveBroadcastMessageModal(props.item._id)"
              )
                v-icon(color="error") delete
    .fab(v-if="hasCreateAnyBroadcastMessageAccess")
      v-layout(column)
        refresh-btn(@click="fetchList")
        v-tooltip(left)
          v-btn(
            slot="activator",
            color="primary",
            data-test="addButton",
            fab,
            @click.stop="showCreateBroadcastMessageModal"
          )
            v-icon add
          span {{ $t('modals.createBroadcastMessage.create.title') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import moment from 'moment';

import { MODALS, BROADCAST_MESSAGES_STATUSES } from '@/constants';

import rightsTechnicalBroadcastMessageMixin from '@/mixins/rights/technical/broadcast-message';

import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';
import SearchField from '@/components/forms/fields/search-field.vue';
import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('broadcastMessage');

export default {
  components: {
    RefreshBtn,
    SearchField,
    BroadcastMessage,
  },
  mixins: [rightsTechnicalBroadcastMessageMixin],
  computed: {
    ...mapGetters({
      broadcastMessages: 'items',
      broadcastMessagesPending: 'pending',
    }),

    headers() {
      return [
        {
          text: this.$t('common.status'),
          value: 'status',
        },
        {
          text: this.$t('common.preview'),
          sortable: false,
        },
        {
          text: this.$t('common.start'),
          value: 'start',
        },
        {
          text: this.$t('common.end'),
          value: 'end',
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
      fetchBroadcastMessagesList: 'fetchList',
      fetchActiveBroadcastMessagesList: 'fetchActiveList',
      createBroadcastMessage: 'create',
      updateBroadcastMessage: 'update',
      removeBroadcastMessage: 'remove',
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
        this.fetchActiveBroadcastMessagesList();

        this.$popups.success({ text: this.$t('success.default') });
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    showCreateBroadcastMessageModal() {
      this.$modals.show({
        name: MODALS.createBroadcastMessage,
        config: {
          action: newMessage =>
            this.callActionWithFetching(() => this.createBroadcastMessage({ data: newMessage })),
        },
      });
    },

    showEditBroadcastMessageModal(message) {
      this.$modals.show({
        name: MODALS.createBroadcastMessage,
        config: {
          message,

          action: newMessage =>
            this.callActionWithFetching(() => this.updateBroadcastMessage({ id: message._id, data: newMessage })),
        },
      });
    },

    showRemoveBroadcastMessageModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.callActionWithFetching(() => this.removeBroadcastMessage({ id })),
        },
      });
    },

    fetchList() {
      this.fetchBroadcastMessagesList();
    },
  },
};
</script>

<style lang="scss" scoped>
  .broadcast-message-cell {
    max-width: 300px;
  }
</style>
