<template>
  <c-advanced-data-table
    :headers="headers"
    :items="apiTokens"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    select-all
    advanced-pagination
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        :tooltip="t('apiToken.deleteSelected')"
        type="delete"
        small
        @click="$emit('remove-selected', selected)"
      />
    </template>

    <template #roles="{ item }">
      <div class="d-flex flex-wrap ga-1">
        <v-chip v-for="role in item.roles" :key="role._id">
          {{ role.name }}
        </v-chip>
      </div>
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
      <c-action-btn
        :tooltip="t('apiToken.delete')"
        type="delete"
        @click.stop="$emit('remove', item)"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    apiTokens: {
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
  },

  emits: ['update:options', 'remove', 'remove-selected'],

  setup() {
    const { t } = useI18n();
    const headers = computed(() => {
      const items = [
        {
          text: t('common.name'),
          value: 'name',
        },
        {
          text: t('common.value'),
          value: 'value',
          sortable: false,
        },
        {
          text: t('common.user'),
          value: 'user.display_name',
          sortable: false,
        },
        {
          text: t('common.role', 2),
          value: 'roles',
          sortable: false,
        },
        {
          text: t('common.created'),
          value: 'created',
        },
        {
          text: t('apiToken.expirationDate'),
          value: 'expired',
        },
        {
          text: t('apiToken.lastUsedDate'),
          value: 'accessed',
        },
      ];

      items.push({
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      });

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
