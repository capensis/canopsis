<template>
  <div>
    <c-advanced-data-table
      :loading="pending"
      :options="options"
      :items="tokens"
      :headers="headers"
      :total-items="totalItems"
      expand
      search
      advanced-pagination
      @update:options="$emit('update:options', $event)"
    >
      <template #expiration_duration="{ item }">
        {{ item.expiration_duration | duration }}
      </template>
      <template #last_used="{ item }">
        {{ item.last_used | date }}
      </template>
      <template #updated="{ item }">
        {{ item.updated | date }}
      </template>
      <template #actions="{ item }">
        <v-layout>
          <c-action-btn
            type="edit"
            @click="$emit('edit', item)"
          />
          <c-action-btn
            type="delete"
            @click="$emit('remove', item)"
          />
        </v-layout>
      </template>
      <template #expand="{ item }">
        <external-auth-tokens-list-expand-panel :external-auth-token="item" />
      </template>
    </c-advanced-data-table>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import ExternalAuthTokensListExpandPanel from './partials/external-auth-tokens-list-expand-panel.vue';

export default {
  components: {
    ExternalAuthTokensListExpandPanel,
  },
  props: {
    tokens: {
      type: Array,
      default: () => [],
    },
    options: {
      type: Object,
      default: () => ({}),
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
  },
  setup() {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('externalAuthToken.tokenName'),
        value: 'name',
        sortable: false,
      },
      {
        text: t('common.description'),
        value: 'description',
        sortable: false,
      },
      {
        text: t('externalAuthToken.tokenExpirationTime'),
        value: 'expiration_duration',
        sortable: false,
      },
      {
        text: t('externalAuthToken.lastUsedDate'),
        value: 'last_used',
        sortable: false,
      },
      {
        text: t('externalAuthToken.lastUpdateDate'),
        value: 'updated',
        sortable: false,
      },
      {
        value: 'actions',
        text: t('common.actionsLabel'),
        sortable: false,
      },
    ]);

    return {
      headers,
    };
  },
};
</script>
