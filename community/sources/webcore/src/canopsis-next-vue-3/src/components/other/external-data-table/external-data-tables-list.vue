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
      <c-tooltip v-if="item.removed_from_config" location="bottom">
        <template #activator="{ props: tooltipProps }">
          <v-icon v-bind="tooltipProps" class="ml-2" color="error">
            warning
          </v-icon>
        </template>
        <span
          v-html="t('externalData.tableRemovedFromConfig', { rules: item.linkedRulesTooltip })"
          class="pre-wrap"
        />
      </c-tooltip>
    </template>

    <template #type="{ item }">
      {{ t(`externalData.tableTypes.${item.type}`) }}
    </template>

    <template #actions="{ item }">
      <c-layout align="center">
        <c-action-btn v-if="updatable" type="edit" @click="$emit('edit', item)" />
        <c-action-btn
          v-if="removable"
          :disabled-button="!!item.deleteTooltip"
          :tooltip="item.deleteTooltip"
          type="delete"
          @click="remove(item)"
        />
      </c-layout>
    </template>

    <template #expand="{ item }">
      <external-data-tables-list-expand-panel :external-data-table="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useLinkedRulesTooltips } from '@/hooks/table/linked-rules-tooltips';

import ExternalDataTablesListExpandPanel from './partials/external-data-tables-list-expand-panel.vue';

const htmlToTooltipText = (html = '') => html
    .replace(/<br\s*\/?>/gi, '\n')
    .replace(/<\/(li|ul)>/gi, '\n')
    .replace(/<[^>]+>/g, '')
    .replace(/\n\s*\n/g, '\n')
    .trim();

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

  emits: ['edit', 'remove', 'update:options'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { getLinkedRulesMessage } = useLinkedRulesTooltips();

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

    const preparedExternalDataTables = computed(() => props.externalDataTables.map((table) => {
        const linkedRulesTooltip = getLinkedRulesMessage(table.linked_rules);

        let deleteTooltip = '';

        if (table.from_config) {
          deleteTooltip = t('externalData.tableCanBeDeletedInConfig');
        } else if (linkedRulesTooltip) {
          deleteTooltip = t('externalData.tableCanBeDeletedAfter', {
            rules: htmlToTooltipText(linkedRulesTooltip),
          });
        }

        return {
          ...table,

          linkedRulesTooltip,
          deleteTooltip,
        };
      }),
    );

    /**
     * Emits a remove event for the specified external data table
     *
     * @param {Object} externalDataTable - The external data table object to be removed
     */
    const remove = externalDataTable => emit('remove', externalDataTable);

    /**
     * Updates the table options and emits the changes to the parent component
     *
     * @param {Object} newOptions - The updated options object
     */
    const updateOptions = newOptions => emit('update:options', newOptions);

    return {
      t,
      headers,
      preparedExternalDataTables,

      remove,
      updateOptions,
    };
  },
};
</script>
