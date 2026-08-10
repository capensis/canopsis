<template>
  <c-advanced-data-table
    :headers="headers"
    :items="roles"
    :loading="pending"
    :options="options"
    :total-items="totalItems"
    :select-all="removable"
    :is-disabled-item="isDisabledRole"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :removable="removable"
        role
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>

    <template #auth_config.inactivity_interval="{ item }">
      {{ durationToString(item.auth_config?.inactivity_interval) }}
    </template>

    <template #auth_config.expiration_interval="{ item }">
      {{ durationToString(item.auth_config?.expiration_interval) }}
    </template>

    <template #actions="{ item }">
      <c-layout gap="1">
        <c-action-btn
          v-if="updatable"
          :disabled-button="!item.editable"
          type="edit"
          @click.stop="$emit('edit', item)"
        />
        <c-action-btn v-if="duplicable" type="duplicate" @click.stop="$emit('duplicate', item)" />
        <c-action-btn
          v-if="removable"
          :disabled-button="!item.deletable"
          type="delete"
          @click.stop="$emit('remove', item)"
        />
      </c-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDurationToString } from '@/helpers/date/duration';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    roles: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    options: {
      type: Object,
      required: true,
    },
    totalItems: {
      type: Number,
      default: 0,
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

  emits: ['update:options', 'edit', 'duplicate', 'remove', 'refresh'],

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'name',
      },
      {
        text: t('role.inactivityInterval'),
        value: 'auth_config.inactivity_interval',
        sortable: false,
      },
      {
        text: t('role.expirationInterval'),
        value: 'auth_config.expiration_interval',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    const durationToString = duration =>
      duration ? convertDurationToString(duration) : t('common.notAvailable');
    const isDisabledRole = ({ deletable = true }) => !deletable;

    return {
      headers,
      durationToString,
      isDisabledRole,
    };
  },
};
</script>
