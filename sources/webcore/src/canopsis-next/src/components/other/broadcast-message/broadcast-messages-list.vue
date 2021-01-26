<template lang="pug">
  c-advanced-data-table.white(
    :headers="headers",
    :items="preparedBroadcastMessages",
    :loading="pending",
    :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES"
  )
    template(slot="items", slot-scope="props")
      tr
        td {{ $t(`tables.broadcastMessages.statuses.${props.item.status}`) }}
        td.broadcast-message-cell
          broadcast-message(:message="props.item.message", :color="props.item.color")
        td {{ props.item.start | date('long', true) }}
        td {{ props.item.end | date('long', true) }}
        td
          v-layout(row)
            c-action-btn(
              v-if="hasUpdateAnyBroadcastMessageAccess",
              type="edit",
              @click="$emit('edit', props.item)"
            )
            c-action-btn(
              v-if="hasDeleteAnyBroadcastMessageAccess",
              type="delete",
              @click="$emit('remove', props.item)"
            )
</template>

<script>
import moment from 'moment';

import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

import rightsTechnicalBroadcastMessageMixin from '@/mixins/rights/technical/broadcast-message';

import BroadcastMessage from '@/components/other/broadcast-message/broadcast-message.vue';

export default {
  components: {
    BroadcastMessage,
  },
  mixins: [rightsTechnicalBroadcastMessageMixin],
  props: {
    broadcastMessages: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
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
