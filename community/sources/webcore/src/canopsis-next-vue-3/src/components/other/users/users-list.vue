<template>
  <c-advanced-data-table
    :headers="headers"
    :items="users"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :select-all="removable"
    advanced-pagination
    search
    @update:options="$emit('update:options', $event)"
  >
    <template #mass-actions="{ selected, clearSelected }">
      <c-table-mass-actions-panel
        :items="selected"
        :enablable="updatable"
        :disablable="updatable"
        :removable="removable"
        user
        @clear:items="clearSelected"
        @refresh="$emit('refresh')"
      />
    </template>

    <template #enabled="{ item }">
      <c-enabled :value="item.enabled" />
    </template>

    <template #active="{ item }">
      <c-icon :color="item.active_connects > 0 ? 'success' : 'error'">
        {{ item.active_connects > 0 ? 'check_circle' : 'cancel' }}
      </c-icon>
    </template>

    <template #source="{ item }">
      {{ item.source || AUTH_METHODS.local }}
    </template>

    <template #roles="{ item }">
      <c-layout gap="1" wrap>
        <v-chip v-for="role in item.roles" :key="role._id">
          {{ role.name }}
        </v-chip>
      </c-layout>
    </template>

    <template #actions="{ item }">
      <c-layout gap="1">
        <c-action-btn v-if="updatable" type="edit" @click.stop="$emit('edit', item)" />
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

import { AUTH_METHODS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    users: {
      type: Array,
      default: () => [],
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    pending: {
      type: Boolean,
      default: false,
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
  },

  emits: ['update:options', 'edit', 'remove', 'refresh'],

  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.username'),
        value: 'name',
      },
      {
        text: t('user.displayName'),
        value: 'display_name',
      },
      {
        text: t('user.firstName'),
        value: 'firstname',
      },
      {
        text: t('user.lastName'),
        value: 'lastname',
      },
      {
        text: t('common.role', 2),
        value: 'roles',
        sortable: false,
      },
      {
        text: t('user.active'),
        value: 'active',
        sortable: false,
      },
      {
        text: t('common.enabled'),
        value: 'enabled',
      },
      {
        text: t('user.auth'),
        value: 'source',
      },
      {
        text: t('user.activeConnects'),
        value: 'active_connects',
        sortable: false,
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    return {
      AUTH_METHODS,
      headers,
    };
  },
};
</script>
