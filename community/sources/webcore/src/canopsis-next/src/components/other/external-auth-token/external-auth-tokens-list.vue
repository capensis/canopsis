<template>
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
      {{ formatDuration(item.expiration_duration) }}
    </template>

    <template #last_used="{ item }">
      <span v-if="item.failed" class="text-error">
        {{ item.fail_reason }}
        <v-icon class="ml-2" color="error"> warning </v-icon>
      </span>
      <span v-else>
        {{ formatDate(item.last_used) }}
      </span>
    </template>

    <template #updated="{ item }">
      {{ formatDate(item.updated) }}
    </template>

    <template #actions="{ item }">
      <c-layout>
        <c-action-btn v-if="editable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn
          v-if="deletable"
          :disabled-button="!!item.deleteTooltip"
          :tooltip="item.deleteTooltip"
          type="delete"
          @click="$emit('remove', item)"
        />
      </c-layout>
    </template>

    <template #expand="{ item }">
      <external-auth-tokens-list-expand-panel :token="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';
import { useLinkedRulesTooltips } from '@/hooks/table/linked-rules-tooltips';

import ExternalAuthTokensListExpandPanel from './partials/external-auth-tokens-list-expand-panel.vue';

const htmlToTooltipText = (html = '') => html
    .replace(/<br\s*\/?>/gi, '\n')
    .replace(/<\/(li|ul)>/gi, '\n')
    .replace(/<[^>]+>/g, '')
    .replace(/\n\s*\n/g, '\n')
    .trim();

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

  emits: ['edit', 'remove', 'update:options'],

  setup(props) {
    const { t } = useI18n();
    const { getLinkedRulesMessage } = useLinkedRulesTooltips();

    const headers = computed(() => [
      {
        text: t('externalAuthToken.tokenName'),
        value: 'name',
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
      },
      {
        text: t('externalAuthToken.lastUpdateDate'),
        value: 'updated',
      },
      {
        value: 'actions',
        text: t('common.actionsLabel'),
        sortable: false,
      },
    ]);

    const preparedTokens = computed(() => props.tokens.map((token) => {
        const linkedRulesTooltip = getLinkedRulesMessage(token.linked_rules);

        if (!linkedRulesTooltip.trim()) {
          return token;
        }

        const deleteTooltip = t('externalAuthToken.tokenCanNotBeDeleted', {
          rules: htmlToTooltipText(linkedRulesTooltip),
        });

        return {
          ...token,

          linkedRulesTooltip,
          deleteTooltip,
        };
      }),
    );

    const isExpandableItem = item => item?.failed;
    const formatDate = date => convertDateToString(date);
    const formatDuration = duration => convertDurationToString(duration);

    return {
      headers,
      preparedTokens,
      isExpandableItem,
      formatDate,
      formatDuration,
    };
  },
};
</script>
