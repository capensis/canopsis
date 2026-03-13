<template>
  <widget-settings-item :title="$t('settings.defaultSortColumns')">
    <v-layout column>
      <c-widget-template-field
        :value="template"
        :templates="templates"
        :pending="templatesPending"
        clearable
        @input="updateTemplate"
      />
      <c-alert
        v-if="missingColumns.length > 0"
        type="warning"
      >
        <div>{{ $t('settings.defaultSortColumnsHasAnotherColumns') }}</div>
        <div class="mt-2">
          <strong>{{ $t('settings.missingColumns') }}:</strong>
          <ul class="mb-0 mt-1">
            <li
              v-for="column in missingColumns"
              :key="column.sort_by"
            >
              {{ column.label }}
            </li>
          </ul>
        </div>
      </c-alert>
      <widget-template-sort-columns-form :value="sortColumns" :items="columns" @input="updateValue" />
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import { ALARM_FIELDS_TO_LABELS_KEYS } from '@/constants';

import { getWidgetColumnLabel } from '@/helpers/entities/widget/list';

import { useWidgetTemplateField } from '@/hooks/widget/widget-template';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import WidgetTemplateSortColumnsForm from '@/components/other/widget-template/form/widget-template-sort-columns-form.vue';

export default {
  components: { WidgetSettingsItem, WidgetTemplateSortColumnsForm },
  model: {
    prop: 'sortColumns',
    event: 'input',
  },
  props: {
    sortColumns: {
      type: Array,
      default: () => [],
    },
    columns: {
      type: Array,
      default: () => [],
    },
    template: {
      type: [String, Symbol],
      required: false,
    },
    templates: {
      type: Array,
      default: () => [],
    },
    templatesPending: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { updateTemplate, updateValue } = useWidgetTemplateField(props, 'sort_columns', emit);

    const selectedTemplate = computed(() => {
      if (!props.template || !props.templates.length) {
        return null;
      }

      return props.templates.find(t => t._id === props.template);
    });

    const currentColumnValues = computed(() => props.columns.map(column => column.value));

    const missingColumns = computed(() => {
      if (!selectedTemplate.value?.sort_columns) {
        return [];
      }

      return selectedTemplate.value.sort_columns
        .filter(sortColumn => sortColumn.sort_by && !currentColumnValues.value.includes(sortColumn.sort_by))
        .map(sortColumn => ({
          ...sortColumn,
          label: getWidgetColumnLabel({ value: sortColumn.sort_by }, ALARM_FIELDS_TO_LABELS_KEYS),
        }));
    });

    return {
      missingColumns,
      updateValue,
      updateTemplate,
    };
  },
};
</script>
