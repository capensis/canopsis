<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="preparedBroadcastMessages",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    advanced-pagination,
    search,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#items="{ item }")
      tr
        td {{ $t(`broadcastMessage.statuses.${item.status}`) }}
        td.broadcast-message-cell
          broadcast-message(:message="item.message", :color="item.color")
        td {{ item.start | date }}
        td {{ item.end | date }}
        td
          v-layout(row)
            c-action-btn(
              v-if="hasUpdateAnyBroadcastMessageAccess",
              type="edit",
              @click="$emit('edit', item)"
            )
            c-action-btn(
              v-if="hasDeleteAnyBroadcastMessageAccess",
              type="delete",
              @click="$emit('remove', item._id)"
            )
</template>

<script>
import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

import { getNowTimestamp } from '@/helpers/date/date';

import { permissionsTechnicalBroadcastMessageMixin } from '@/mixins/permissions/technical/broadcast-message';

import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';

export default {
  components: {
    BroadcastMessage,
  },
  mixins: [permissionsTechnicalBroadcastMessageMixin],
  props: {
    broadcastMessages: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pagination: {
      type: Object,
      required: true,
    },
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.status'),
          value: 'status',
          sortable: false,
        },
        {
          text: this.$t('common.preview'),
          sortable: false,
        },
        {
          text: this.$t('common.start'),
          value: 'start',
          sortable: false,
        },
        {
          text: this.$t('common.end'),
          value: 'end',
          sortable: false,
        },
        {
          text: this.$t('common.actionsLabel'),
          sortable: false,
        },
      ];
    },

    preparedBroadcastMessages() {
      return this.broadcastMessages.map((message) => {
        const now = getNowTimestamp();
        let status = BROADCAST_MESSAGES_STATUSES.pending;

        if (now >= message.start) {
          status = now <= message.end
            ? BROADCAST_MESSAGES_STATUSES.active
            : BROADCAST_MESSAGES_STATUSES.expired;
        }

        return {
          ...message,

          status,
        };
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .broadcast-message-cell {
    max-width: 300px;
  }
</style>
