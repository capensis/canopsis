<template>
  <v-card class="column-field">
    <v-tooltip left>
      <template #activator="{ on }">
        <v-btn
          class="column-field__remove-btn"
          small
          text
          icon
          v-on="on"
          @click="$emit('remove')"
        >
          <v-icon
            color="error"
            small
          >
            close
          </v-icon>
        </v-btn>
      </template>
      <span>{{ $t('common.delete') }}</span>
    </v-tooltip>
    <v-card-text>
      <v-layout align-center>
        <span class="handler mr-1">
          <v-icon
            :class="dragHandleClass"
            class="draggable"
          >
            drag_indicator
          </v-icon>
        </span>
        <c-expand-btn
          v-model="expanded"
          :color="hasChildrenError ? 'error' : ''"
          class="mr-1"
        />
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
          @change="changeColumn"
        />
        <v-tooltip left>
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
      <v-expand-transition mode="out-in">
        <column-field-expand-panel
          v-if="bootedExpandPanel"
          v-show="expanded"
          v-field="column"
          :name="name"
          :with-label="!isCustom"
          :with-field="isCustom"
          :with-html="withHtml"
          :with-template="withTemplate"
          :with-color-indicator="withColorIndicator"
          :with-instructions="withInstructions"
          :with-simple-template="withSimpleTemplate"
          :optional-infos-attributes="optionalInfosAttributes"
          :without-infos-attributes="withoutInfosAttributes"
          :variables="variables"
          class="pl-1"
        />
      </v-expand-transition>
    </v-card-text>
  </v-card>
</template>

<script>
import { omit } from 'lodash';
import { computed, ref } from 'vue';

import {
  ENTITIES_TYPES,
  ALARM_LIST_WIDGET_COLUMNS,
  CONTEXT_WIDGET_COLUMNS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ALARM_OUTPUT_FIELDS,
} from '@/constants';

import { formToWidgetColumn, widgetColumnValueToForm } from '@/helpers/entities/widget/column/form';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';
import { useValidationChildren } from '@/hooks/validator/validation-children';
import { useModelField } from '@/hooks/form/model-field';
import { useAsyncBootingChild } from '@/hooks/render/async-booting';

import ColumnFieldExpandPanel from './column-field-expand-panel.vue';

export default {
  inject: ['$validator', '$asyncBooting'],
  components: { ColumnFieldExpandPanel },
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
    const expanded = ref(!props.column?.column);
    const isCustom = ref(false);

    const { tc } = useI18n();
    const validator = useValidator();
    const { updateModel } = useModelField(props, emit);
    const { booted: bootedExpandPanel } = useAsyncBootingChild(expanded.value);
    const { hasChildrenError } = useValidationChildren();

    /**
     * COMPUTED
     */
    const isAlarmType = computed(() => props.type === ENTITIES_TYPES.alarm);

    const alarmListAvailableColumns = computed(() => {
      const columns = props.withInstructions
        ? ALARM_LIST_WIDGET_COLUMNS
        : omit(ALARM_LIST_WIDGET_COLUMNS, ['assignedInstructions']);

      return Object.values(columns).map(value => ({
        value,
        text: tc(ALARM_FIELDS_TO_LABELS_KEYS[value], 2),
      }));
    });

    const contextAvailableColumns = computed(() => Object.values(CONTEXT_WIDGET_COLUMNS).map(value => ({
      value,
      text: tc(ENTITY_FIELDS_TO_LABELS_KEYS[value], 2),
    })));

    const availableColumns = computed(() => {
      const columns = isAlarmType.value
        ? alarmListAvailableColumns.value
        : contextAvailableColumns.value;

      return columns.filter(({ value }) => !props.excludedColumns.includes(value));
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

      this.updateModel(newValue);
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
      expanded,
      isCustom,
      bootedExpandPanel,
      availableColumns,
      columnLabelFieldName,
      columnLabelErrorMessages,
      hasChildrenError,

      changeColumn,
      convertToCustom,
    };
  },
};
</script>

<style lang="scss">
.column-field {
  position: relative;

  &__remove-btn.v-btn {
    position: absolute;
    right: 0;
    top: 0;
  }
}
</style>
