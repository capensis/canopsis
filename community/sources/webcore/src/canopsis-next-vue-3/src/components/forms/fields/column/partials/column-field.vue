<template>
  <c-card-iterator-item
    :drag-handle-class="dragHandleClass"
    :default-expanded="isExpanded"
    :expandable="!withoutCustomLabel"
    class="column-field"
    small
    @remove="$emit('remove')"
  >
    <template #header>
      <c-layout class="column-field__header" align="center" fill>
        <c-name-field
          v-if="isCustom"
          :model-value="column.label"
          :name="columnLabelFieldName"
          :label="$t('common.label')"
          class="flex-grow-1"
          required
          @update:model-value="updateField('label', $event)"
        />
        <c-select
          v-else
          :model-value="column.column"
          :items="availableColumns"
          :label="$t('common.column', 1)"
          :name="`${name}.column`"
          class="flex-grow-1"
          item-title="text"
          item-value="value"
          required
          @update:model-value="changeColumn"
        />
        <c-tooltip v-if="!withoutCustomLabel" location="left">
          <template #activator="{ props }">
            <v-btn
              :class="isCustom ? 'text-primary' : 'text-disabled'"
              class="mr-0"
              size="small"
              variant="text"
              icon
              v-bind="props"
              @click="convertToCustom"
            >
              <v-icon size="small">
                tune
              </v-icon>
            </v-btn>
          </template>
          <span>{{ $t('common.convertToCustomColumn') }}</span>
        </c-tooltip>
      </c-layout>
    </template>

    <v-expand-transition v-if="!withoutCustomLabel" mode="out-in">
      <column-field-expand-panel
        :value="column"
        :name="name"
        :with-label="!isCustom"
        :with-field="isCustom"
        :with-html="withHtml"
        :with-template="withTemplate"
        :with-color-indicator="withColorIndicator"
        :with-instructions="withInstructions"
        :with-simple-template="withSimpleTemplate"
        :with-filter-on-click="withFilterOnClick"
        :optional-infos-attributes="optionalInfosAttributes"
        :without-infos-attributes="withoutInfosAttributes"
        :variables="variables"
        class="pl-1"
        @update:value="updateModel"
      />
    </v-expand-transition>
  </c-card-iterator-item>
</template>

<script>
import { computed, ref, toRef } from 'vue';

import { ALARM_OUTPUT_FIELDS, ENTITIES_TYPES } from '@/constants';

import { formToWidgetColumn, widgetColumnValueToForm } from '@/helpers/entities/widget/column/form';

import { useAvailableColumns } from '@/hooks/form/available-columns';
import { useModelField } from '@/hooks/form/model-field';

import ColumnFieldExpandPanel from './column-field-expand-panel.vue';

export default {
  components: { ColumnFieldExpandPanel },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    type: {
      type: String,
      default: ENTITIES_TYPES.alarm,
    },
    items: {
      type: Array,
      required: false,
    },
    name: {
      type: String,
      default: '',
    },
    dragHandleClass: {
      type: String,
      default: 'drag-handle',
    },
    withTemplate: {
      type: Boolean,
      default: false,
    },
    withHtml: {
      type: Boolean,
      default: false,
    },
    withColorIndicator: {
      type: Boolean,
      default: false,
    },
    withInstructions: {
      type: Boolean,
      default: false,
    },
    withoutInfosAttributes: {
      type: Boolean,
      default: false,
    },
    withoutCustomLabel: {
      type: Boolean,
      default: false,
    },
    withFilterOnClick: {
      type: Boolean,
      default: false,
    },
    optionalInfosAttributes: {
      type: Boolean,
      default: false,
    },
    withSimpleTemplate: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      required: false,
    },
    excludedColumns: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue', 'remove'],

  setup(props, { emit }) {
    const isCustom = ref(false);
    const isExpanded = ref(!props.modelValue.column);
    const { updateModel, updateField } = useModelField(props, emit);
    const { availableColumns } = useAvailableColumns({
      type: toRef(props, 'type'),
      items: toRef(props, 'items'),
      excludedColumns: toRef(props, 'excludedColumns'),
      withInstructions: toRef(props, 'withInstructions'),
    });

    const column = toRef(props, 'modelValue');
    const columnLabelFieldName = computed(() => `${props.name}.label`);

    const changeColumn = (columnValue) => {
      const newValue = {
        ...props.modelValue,
        column: columnValue,
      };

      if (props.withHtml) {
        newValue.isHtml = ALARM_OUTPUT_FIELDS.includes(columnValue);
      }

      updateModel(newValue);
    };

    const convertToCustom = () => {
      isCustom.value = !isCustom.value;

      const newColumn = { ...props.modelValue };

      if (isCustom.value) {
        const { value } = formToWidgetColumn(props.modelValue);
        const selectedColumn = availableColumns.value.find(
          columnItem => columnItem.value === props.modelValue.column,
        );
        const label = props.modelValue.label || selectedColumn?.text || '';

        newColumn.column = value;
        newColumn.label = label;
        newColumn.field = '';
        newColumn.rule = '';
        newColumn.dictionary = '';
      } else {
        const {
          column: value,
          field,
          rule,
          dictionary,
        } = widgetColumnValueToForm(props.modelValue.column);
        const selectedColumn = availableColumns.value.find(
          columnItem => columnItem.value === value,
        );

        newColumn.column = value === selectedColumn?.value ? value : '';
        newColumn.label = newColumn.label === selectedColumn?.text ? '' : newColumn.label;
        newColumn.field = field;
        newColumn.rule = rule;
        newColumn.dictionary = dictionary;
      }

      updateModel(newColumn);
    };

    return {
      isExpanded,
      isCustom,
      column,
      availableColumns,
      columnLabelFieldName,
      updateModel,
      updateField,
      changeColumn,
      convertToCustom,
    };
  },
};
</script>

<style lang="scss">
.column-field {
  &__header {
    min-width: 0;
  }

  .v-input {
    min-width: 0;
    width: 100%;
  }
}
</style>
