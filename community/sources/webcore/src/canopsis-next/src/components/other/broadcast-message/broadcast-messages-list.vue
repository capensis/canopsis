<template>
  <c-advanced-data-table
    :headers="headers"
    :items="broadcastMessages"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    search
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #message="{ item }">
      <broadcast-message :message="item.message" :color="item.color" />
    </template>

    <template #status="{ item }">
      <span
        :class="{
          'text-success': item.status === activeStatus,
          'text-grey': item.status !== activeStatus,
        }"
      >
        {{ t(`broadcastMessage.statuses.${item.status}`) }}
      </span>
    </template>

    <template #start="{ item }">
      {{ convertDateToString(item.start) }}
    </template>

    <template #end="{ item }">
      {{ convertDateToString(item.end) }}
    </template>

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn v-if="editable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn v-if="deletable" type="delete" @click="$emit('remove', item._id)" />
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { BROADCAST_MESSAGES_STATUSES } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import BroadcastMessage from './partials/broadcast-message.vue';

export default {
  components: {
    BroadcastMessage,
  },

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
    editable: {
      type: Boolean,
      default: false,
    },
    deletable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:options', 'edit', 'remove'],

  setup() {
    const { t } = useI18n();
    const activeStatus = BROADCAST_MESSAGES_STATUSES.active;

    const headers = computed(() => [
      {
        text: t('common.preview'),
        value: 'message',
        width: 400,
      },
      {
        text: t('common.status'),
        value: 'status',
      },
      {
        text: t('common.start'),
        value: 'start',
      },
      {
        text: t('common.end'),
        value: 'end',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    return {
      t,
      activeStatus,
      headers,
      convertDateToString,
    };
  },
};
</script>
