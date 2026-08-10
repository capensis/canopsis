<template>
  <c-advanced-data-table
    :headers="headers"
    :items="themes"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    :is-disabled-item="isDisabledTheme"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-show="removable"
        type="delete"
        small
        @click="$emit('remove-selected', selected)"
      />
    </template>

    <template #updated="{ item }">
      {{ dateToString(item.updated) }}
    </template>

    <template #actions="{ item }">
      <c-layout gap="1">
        <c-action-btn
          v-if="updatable"
          :tooltip="item.deletable ? t('common.edit') : t('theme.defaultTheme')"
          :disabled-button="!item.deletable"
          type="edit"
          @click.stop="$emit('edit', item)"
        />
        <c-action-btn v-if="duplicable" type="duplicate" @click.stop="$emit('duplicate', item)" />
        <c-action-btn
          v-if="removable"
          :tooltip="item.deletable ? t('common.delete') : t('theme.defaultTheme')"
          :disabled-button="!item.deletable"
          type="delete"
          @click.stop="$emit('remove', item._id)"
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
    themes: {
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
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:options', 'edit', 'remove', 'duplicate', 'remove-selected'],

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.title'),
        value: 'name',
      },
      {
        text: t('common.lastModifiedOn'),
        value: 'updated',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const dateToString = date => convertDateToString(date);
    const isDisabledTheme = ({ deletable }) => !deletable;

    return {
      t,
      headers,
      dateToString,
      isDisabledTheme,
    };
  },
};
</script>
