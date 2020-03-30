<template lang="pug">
  div
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.roles') }}
    div.white
      v-data-table(
        :headers="headers",
        :items="broadcastMessages",
        :loading="broadcastMessagesPending",
        :pagination.sync="pagination",
        :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
        :total-items="broadcastMessages.length",
        item-key="_id"
      )
        template(slot="items", slot-scope="props")
          tr(:data-test="`role-${props.item._id}`")
            td Expired
            td.broadcast-message-cell
              broadcast-message(:message="props.item.message", :color="props.item.color")
            td
              enabled-column(:value="props.item.enabled")
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
          span {{ $t('modals.createBroadcastMessage.title') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import viewQuery from '@/mixins/view/query';
import rightsTechnicalBroadcastMessageMixin from '@/mixins/rights/technical/broadcast-message';

import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';
import SearchField from '@/components/forms/fields/search-field.vue';
import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';
import EnabledColumn from '@/components/tables/enabled-column.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('broadcastMessage');

export default {
  components: {
    RefreshBtn,
    SearchField,
    BroadcastMessage,
    EnabledColumn,
  },
  mixins: [
    viewQuery,
    rightsTechnicalBroadcastMessageMixin,
  ],
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
          text: this.$t('common.enabled'),
          value: 'enabled',
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
  },
  methods: {
    ...mapActions({
      fetchBroadcastMessagesList: 'fetchList',
      createBroadcastMessage: 'create',
      updateBroadcastMessage: 'update',
      deleteBroadcastMessage: 'delete',
    }),

    showCreateBroadcastMessageModal() {
      this.$modals.show({
        name: MODALS.createBroadcastMessage,
        action: async () => {
          try {
            this.$popups.success({ text: this.$t('success.default') });
          } catch (err) {
            this.$popups.error({ text: this.$t('errors.default') });
          }
        },
      });
    },

    showEditBroadcastMessageModal(message) {
      this.$modals.show({
        name: MODALS.createBroadcastMessage,
        config: {
          message,
          action: async () => {
            try {
              this.$popups.success({ text: this.$t('success.default') });
            } catch (err) {
              this.$popups.error({ text: this.$t('errors.default') });
            }
          },
        },
      });
    },

    showRemoveBroadcastMessageModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              this.$popups.success({ text: this.$t('success.default') });
            } catch (err) {
              this.$popups.error({ text: this.$t('errors.default') });
            }
          },
        },
      });
    },

    fetchList() {
      this.fetchBroadcastMessagesList({ params: this.getQuery() });
    },
  },
};
</script>

<style lang="scss" scoped>
  .broadcast-message-cell {
    max-width: 350px;
  }
</style>
