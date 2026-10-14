<template>
  <card-iterator-item
    :drag-handle-class="dragHandleClass"
    :default-expanded="isExpanded"
    class="column-field"
    small
    @remove="$emit('remove')"
  >
    <template #header>
      <v-layout align-center>
        <c-name-field
          v-if="isCustom"
          v-field="column.label"
          :name="columnLabelFieldName"
          :label="$t('common.label')"
          :error-messages="columnLabelErrorMessages"
          required
        />
        <v-select
          v-else
          v-validate="'required'"
          :value="column.column"
          :items="availableColumns"
          :label="$tc('common.column', 1)"
          :error-messages="errors.collect(`${name}.column`)"
          :name="`${name}.column`"
          :menu-props="selectMenuProps"
          @change="changeColumn"
        />
        <v-tooltip v-if="!withoutCustomLabel" left>
          <template #activator="{ on }">
            <v-btn
              :class="`mr-0 ${isCustom ? 'text--primary' : 'text--disabled'}`"
              small
              text
              icon
              v-on="on"
              @click="convertToCustom"
            >
              <v-icon small>
                tune
              </v-icon>
            </v-btn>
          </template>
          <span>{{ $t('common.convertToCustomColumn') }}</span>
        </v-tooltip>
      </v-layout>
    </template>

    <v-expand-transition v-if="!withoutCustomLabel" mode="out-in">
      <column-field-expand-panel
        v-if="bootedExpandPanel"
        v-field="column"
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
      />
    </v-expand-transition>
  </card-iterator-item>
</template>

<script>
import { computed, ref, toRef } from 'vue';

import { ENTITIES_TYPES, ALARM_OUTPUT_FIELDS } from '@/constants';

import { formToWidgetColumn, widgetColumnValueToForm } from '@/helpers/entities/widget/column/form';

import { useValidator } from '@/hooks/validator/validator';
import { useValidationChildren } from '@/hooks/validator/validation-children';
import { useModelField } from '@/hooks/form/model-field';
import { useAsyncBootingChild } from '@/hooks/render/async-booting';
import { useAvailableColumns } from '@/hooks/form/available-columns';

import CardIteratorItem from '@/components/forms/fields/card-iterator/c-card-iterator-item.vue';

import ColumnFieldExpandPanel from './column-field-expand-panel.vue';

export default {
  inject: ['$validator', '$asyncBooting'],
  components: { ColumnFieldExpandPanel, CardIteratorItem },
  model: {
    prop: 'column',
    event: 'input',
  },
  props: {
    type: {
      type: String,
      default: ENTITIES_TYPES.alarm,
    },
    column: {
      type: Object,
      default: () => ({}),
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
  setup(props, { emit }) {
    const isCustom = ref(false);
    const isExpanded = ref(!props.column.column);

    const validator = useValidator();
    const { updateModel } = useModelField(props, emit);
    const { booted: bootedExpandPanel } = useAsyncBootingChild(!props.column?.column);
    const { hasChildrenError } = useValidationChildren();

    const selectMenuProps = {
      contentClass: 'column-field-menu',
    };

    const { availableColumns } = useAvailableColumns({
      type: toRef(props, 'type'),
      items: toRef(props, 'items'),
      excludedColumns: toRef(props, 'excludedColumns'),
      withInstructions: toRef(props, 'withInstructions'),
    });

    const columnLabelFieldName = computed(() => `${props.name}.label`);
    const columnLabelErrorMessages = computed(() => validator.errors.collect(columnLabelFieldName.value));

    /**
     * METHODS
     */
    const changeColumn = (column) => {
      const newValue = {
        ...props.column,

        column,
      };

      if (props.withHtml) {
        newValue.isHtml = ALARM_OUTPUT_FIELDS.includes(column);
      }

      updateModel(newValue);
    };

    const convertToCustom = () => {
      isCustom.value = !isCustom.value;

      const newColumn = {
        ...props.column,
      };

      if (isCustom.value) {
        const { value } = formToWidgetColumn(props.column);

        const selectedColumn = availableColumns.value.find(column => column.value === props.column.column);

        const label = props.column.label || selectedColumn?.text || '';

        newColumn.column = value;
        newColumn.label = label;
        newColumn.field = '';
        newColumn.rule = '';
        newColumn.dictionary = '';
      } else {
        const { column: value, field, rule, dictionary } = widgetColumnValueToForm(props.column.column);

        const selectedColumn = availableColumns.value.find(column => column.value === value);

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
      bootedExpandPanel,
      availableColumns,
      columnLabelFieldName,
      columnLabelErrorMessages,
      hasChildrenError,
      selectMenuProps,

      changeColumn,
      convertToCustom,
    };
  },
};
</script>

<style lang="scss">
.column-field .v-input {
  min-width: 236px;
}
</style>
