<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <broadcast-messages-list
        :broadcast-messages="broadcastMessages"
        :pending="broadcastMessagesPending"
        :options.sync="options"
        :total-items="broadcastMessagesMeta.total_count"
        @edit="showEditBroadcastMessageModal"
        @remove="showRemoveBroadcastMessageModal"
      />
    </v-card>
    <c-fab-btn
      :has-access="hasCreateAnyBroadcastMessageAccess"
      @refresh="fetchList"
      @create="showCreateBroadcastMessageModal"
    >
      <span>{{ $t('modals.createBroadcastMessage.create.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { permissionsTechnicalBroadcastMessageMixin } from '@/mixins/permissions/technical/broadcast-message';
import { entitiesBroadcastMessageMixin } from '@/mixins/entities/broadcast-message';
import { localQueryMixin } from '@/mixins/query-local/query';

import BroadcastMessagesList from '@/components/other/broadcast-message/broadcast-messages-list.vue';

export default {
  components: {
    BroadcastMessagesList,
  },
  mixins: [
    localQueryMixin,
    permissionsTechnicalBroadcastMessageMixin,
    entitiesBroadcastMessageMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
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
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },

    showCreateBroadcastMessageModal() {
      this.$modals.show({
        name: MODALS.createBroadcastMessage,
        config: {
          action: newMessage => this.callActionWithFetching(
            () => this.createBroadcastMessage({ data: newMessage }),
          ),
        },
      });
    },

    showEditBroadcastMessageModal(message) {
      this.$modals.show({
        name: MODALS.createBroadcastMessage,
        config: {
          message,
          title: this.$t('modals.createBroadcastMessage.edit.title'),

          action: newMessage => this.callActionWithFetching(
            () => this.updateBroadcastMessage({ id: message._id, data: newMessage }),
          ),
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
      this.fetchBroadcastMessagesList({ params: this.getQuery() });
    },
  },
};
</script>
