<template lang="pug">
  div
    c-the-page-header {{ $t('common.broadcastMessages') }}
    broadcast-messages-list(
      :broadcast-messages="broadcastMessages",
      :pending="broadcastMessagesPending",
      @edit="showEditBroadcastMessageModal",
      @remove="showRemoveBroadcastMessageModal"
    )
    c-fab-btn(
      v-if="hasCreateAnyBroadcastMessageAccess",
      @refresh="fetchList",
      @create="showCreateBroadcastMessageModal"
    )
      span {{ $t('modals.createBroadcastMessage.create.title') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import rightsTechnicalBroadcastMessageMixin from '@/mixins/rights/technical/broadcast-message';

import BroadcastMessagesList from '@/components/other/broadcast-message/broadcast-messages-list.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('broadcastMessage');

export default {
  components: {
    BroadcastMessagesList,
  },
  mixins: [rightsTechnicalBroadcastMessageMixin],
  computed: {
    ...mapGetters({
      broadcastMessages: 'items',
      broadcastMessagesPending: 'pending',
    }),
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
