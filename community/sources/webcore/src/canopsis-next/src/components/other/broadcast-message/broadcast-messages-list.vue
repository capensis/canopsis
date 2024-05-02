<template>
  <c-advanced-data-table
    :headers="headers"
    :items="preparedBroadcastMessages"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #status="{ item }">
      {{ $t(`broadcastMessage.statuses.${item.status}`) }}
    </template>
    <template #message="{ item }">
      <broadcast-message
        :message="item.message"
        :color="item.color"
      />
    </template>
    <template #start="{ item }">
      {{ item.start | date }}
    </template>
    <template #end="{ item }">
      {{ item.end | date }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="hasUpdateAnyBroadcastMessageAccess"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="hasDeleteAnyBroadcastMessageAccess"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

import { getNowTimestamp } from '@/helpers/date/date';

import { permissionsTechnicalBroadcastMessageMixin } from '@/mixins/permissions/technical/broadcast-message';

import BroadcastMessage from '@/components/other/broadcast-message/partials/broadcast-message.vue';

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
    options: {
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
          value: 'message',
          width: 300,
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
          value: 'actions',
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
