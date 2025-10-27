<template>
  <c-advanced-data-table
    :headers="headers"
    :items="preparedExternalDataTables"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    expand
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #name="{ item }">
      <span>{{ item.name }}</span>
      <v-tooltip v-if="item.removed_from_config" bottom>
        <template #activator="{ on }">
          <v-icon class="ml-2" color="error" v-on="on">
            warning
          </v-icon>
        </template>
        <span
          v-html="$t('externalData.tableRemovedFromConfig', { rules: item.linkedRulesTooltip })"
          class="pre-wrap"
        />
      </v-tooltip>
    </template>
    <template #type="{ item }">
      {{ $t(`externalData.tableTypes.${item.type}`) }}
    </template>
    <template #actions="{ item }">
      <v-layout align-center>
        <c-action-btn
          v-if="updatable"
          type="edit"
          @click="$emit('edit', item)"
        />
        <c-action-btn
          v-if="removable"
          :disabled-button="!!item.deleteTooltip"
          :tooltip="item.deleteTooltip"
          type="delete"
          @click="remove(item)"
        />
      </v-layout>
    </template>
    <template #expand="{ item }">
      <external-data-tables-list-expand-panel :external-data-table="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { MAX_EXTERNAL_DATA_TABLE_TOOLTIP_LINKED_RULES_COUNT } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import ExternalDataTablesListExpandPanel from './partials/external-data-tables-list-expand-panel.vue';

export default {
  components: {
    ExternalDataTablesListExpandPanel,
  },
  props: {
    externalDataTables: {
      type: Array,
      default: () => [],
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
    updatable: {
      type: Boolean,
      default: false,
    },
    removable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const headers = computed(() => [
      {
        text: t('common.name'),
        value: 'name',
      },
      {
        text: t('common.description'),
        value: 'description',
      },
      {
        text: t('common.database'),
        value: 'type',
      },
      {
        text: t('common.actionsLabel'),
        value: 'actions',
        sortable: false,
      },
    ]);

    /**
     * Formats a list of linked rules of a specific type for display in a tooltip
     *
     * @param {Array} [typeLinkedRules=[]] - Array of linked rules of a specific type
     * @param {string} typeLinkedRules[].name - The name of the linked rule
     * @returns {string} HTML string containing list items with rule names, limited to the maximum count
     */
    const getLinkedRulesMessageForType = (typeLinkedRules = []) => {
      const hasMore = typeLinkedRules.length > MAX_EXTERNAL_DATA_TABLE_TOOLTIP_LINKED_RULES_COUNT;

      const result = typeLinkedRules.slice(0, MAX_EXTERNAL_DATA_TABLE_TOOLTIP_LINKED_RULES_COUNT);

      if (hasMore) {
        result.push({ name: t('externalData.andMore') });
      }

      return result.map(item => `<li>${item.name}</li>`).join('');
    };

    /**
     * Generates a complete message about linked rules for an external data table
     *
     * Combines messages for different types of linked rules (widgets, event filters, and link rules)
     * into a single HTML string for display in a tooltip.
     *
     * @param {Object} [linkedRules={}] - Object containing arrays of different types of linked rules
     * @param {Array} [linkedRules.widget=[]] - Array of widget rules linked to the data table
     * @param {Array} [linkedRules.eventfilter=[]] - Array of event filter rules linked to the data table
     * @param {Array} [linkedRules.linkrule=[]] - Array of link rules linked to the data table
     * @returns {string} HTML string containing formatted messages for all types of linked rules
     */
    const getLinkedRulesMessage = (linkedRules = {}) => {
      const widgetsMessages = getLinkedRulesMessageForType(linkedRules.widget);
      const eventFiltersMessages = getLinkedRulesMessageForType(linkedRules.eventfilter);
      const linkRulesMessages = getLinkedRulesMessageForType(linkedRules.linkrule);

      let message = '';

      if (widgetsMessages) {
        message += t('externalData.linkedRules.widgets', { rules: widgetsMessages });
      }

      if (eventFiltersMessages) {
        message += t('externalData.linkedRules.eventFilters', { rules: eventFiltersMessages });
      }

      if (linkRulesMessages) {
        message += t('externalData.linkedRules.links', { rules: linkRulesMessages });
      }

      return message;
    };

    const preparedExternalDataTables = computed(() => props.externalDataTables.map((table) => {
      const linkedRulesTooltip = getLinkedRulesMessage(table.linked_rules);

      let deleteTooltip = '';

      if (table.from_config) {
        deleteTooltip = t('externalData.tableCanBeDeletedInConfig');
      } else if (linkedRulesTooltip) {
        deleteTooltip = t('externalData.tableCanBeDeletedAfter', { rules: linkedRulesTooltip });
      }

      return {
        ...table,

        linkedRulesTooltip,
        deleteTooltip: deleteTooltip ? `<span class="pre-wrap">${deleteTooltip}</span>` : '',
      };
    }));

    /**
     * Emits a remove event for the specified external data table
     *
     * @param {Object} externalDataTable - The external data table object to be removed
     */
    const remove = externalDataTable => emit('remove', externalDataTable);

    /**
     * Updates the table options and emits the changes to the parent component
     *
     * @param {Object} options - The updated options object
     */
    const updateOptions = options => emit('update:options', options);

    return {
      headers,
      preparedExternalDataTables,

      remove,
      updateOptions,
    };
  },
};
</script>
