<template>
  <div>
    <c-advanced-data-table
      :loading="pending"
      :options="options"
      :items="preparedTokens"
      :headers="headers"
      :total-items="totalItems"
      :is-expandable-item="isExpandableItem"
      expand
      search
      advanced-pagination
      @update:options="$emit('update:options', $event)"
    >
      <template #expiration_duration="{ item }">
        {{ item.expiration_duration | duration }}
      </template>
      <template #last_used="{ item }">
        <span v-if="item.failed" class="error--text">
          {{ item.fail_reason }}
          <v-icon class="ml-2" color="error">warning</v-icon>
        </span>
        <span v-else>
          {{ item.last_used | date }}
        </span>
      </template>
      <template #updated="{ item }">
        {{ item.updated | date }}
      </template>
      <template #actions="{ item }">
        <v-layout>
          <c-action-btn
            v-if="editable"
            type="edit"
            @click="$emit('edit', item)"
          />
          <c-action-btn
            v-if="deletable"
            :disabled-button="!!item.deleteTooltip"
            :tooltip="item.deleteTooltip"
            type="delete"
            @click="$emit('remove', item)"
          />
        </v-layout>
      </template>
      <template #expand="{ item }">
        <external-auth-tokens-list-expand-panel :token="item" />
      </template>
    </c-advanced-data-table>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useLinkedRulesTooltips } from '@/hooks/table/linked-rules-tooltips';

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
    editable: {
      type: Boolean,
      default: false,
    },
    deletable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { getLinkedRulesMessage } = useLinkedRulesTooltips();

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

    const preparedTokens = computed(() => props.tokens.map((token) => {
      const linkedRulesTooltip = getLinkedRulesMessage(token.linked_rules);

      const deleteTooltip = t('externalAuthToken.tokenCanNotBeDeleted', { rules: linkedRulesTooltip });

      return {
        ...token,

        linkedRulesTooltip,
        deleteTooltip: deleteTooltip ? `<span class="pre-wrap">${deleteTooltip}</span>` : '',
      };
    }));

    const isExpandableItem = item => item?.failed;

    return {
      headers,
      preparedTokens,
      isExpandableItem,
    };
  },
};
</script>
