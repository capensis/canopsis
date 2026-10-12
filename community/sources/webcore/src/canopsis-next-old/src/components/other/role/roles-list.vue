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
        small
        role
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>
    <template #auth_config.inactivity_interval="{ item }">
      {{ durationToString(item.auth_config.inactivity_interval) }}
    </template>
    <template #auth_config.expiration_interval="{ item }">
      {{ durationToString(item.auth_config.expiration_interval) }}
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="updatable"
          :disabled="!item.editable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="duplicable"
          type="duplicate"
          @click="$emit('duplicate', item)"
        />
        <c-action-btn
          v-if="removable"
          :disabled="!item.deletable"
          type="delete"
          @click="$emit('remove', item._id)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    roles: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    options: {
      type: Object,
      required: false,
    },
    totalItems: {
      type: Number,
      required: false,
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

    /**
     * Convert duration to string
     */
    const durationToString = duration => (
      duration ? `${duration.value}${duration.unit}` : t('common.notAvailable')
    );

    /**
     * Check if role is disabled
     */
    const isDisabledRole = ({ deletable = true }) => !deletable;

    return {
      headers,
      durationToString,
      isDisabledRole,
    };
  },
};
</script>
