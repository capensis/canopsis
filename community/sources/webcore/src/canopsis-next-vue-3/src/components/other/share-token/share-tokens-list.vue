<template>
  <c-advanced-data-table
    :headers="headers"
    :items="shareTokens"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        :tooltip="t('shareToken.revokeSelectedTokens')"
        type="delete"
        small
        @click="$emit('remove-selected', selected)"
      />
    </template>

    <template #created="{ item }">
      {{ dateToString(item.created) }}
    </template>

    <template #expired="{ item }">
      {{ dateToString(item.expired, notAvailableText) }}
    </template>

    <template #accessed="{ item }">
      {{ dateToString(item.accessed, notAvailableText) }}
    </template>

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn
          v-if="removable"
          :tooltip="t('shareToken.revokeToken')"
          type="delete"
          @click.stop="$emit('remove', item)"
        />
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    shareTokens: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    options: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:options', 'remove', 'remove-selected'],

  setup(props) {
    const { t } = useI18n();

    const headers = computed(() => {
      const items = [
        {
          text: t('common.id'),
          value: '_id',
          sortable: false,
        },
        {
          text: t('common.description'),
          value: 'description',
        },
        {
          text: t('common.value'),
          value: 'value',
          sortable: false,
        },
        {
          text: t('common.user'),
          value: 'user.display_name',
        },
        {
          text: t('common.role'),
          value: 'role.name',
          sortable: false,
        },
        {
          text: t('common.created'),
          value: 'created',
        },
        {
          text: t('common.expired'),
          value: 'expired',
        },
        {
          text: t('common.accessed'),
          value: 'accessed',
        },
      ];

      if (props.removable) {
        items.push({
          text: t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        });
      }

      return items;
    });

    const notAvailableText = computed(() => t('common.notAvailable'));
    const dateToString = (date, defaultValue) => convertDateToString(date, undefined, defaultValue);

    return {
      t,
      headers,
      notAvailableText,
      dateToString,
    };
  },
};
</script>
