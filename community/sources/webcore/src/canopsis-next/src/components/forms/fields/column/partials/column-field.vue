<template>
  <v-card class="column-field">
    <v-tooltip left>
      <template #activator="{ on }">
        <v-btn
          class="column-field__remove-btn"
          v-on="on"
          small
          text
          icon
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
            class="draggable"
            :class="dragHandleClass"
          >
            drag_indicator
          </v-icon>
        </span>
        <c-expand-btn
          class="mr-1"
          v-model="expanded"
          :color="hasChildrenError ? 'error' : ''"
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
              v-on="on"
              small
              text
              icon
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
          class="pl-1"
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
        />
      </v-expand-transition>
    </v-card-text>
  </v-card>
</template>

<script>
import { omit } from 'lodash';

import {
  ENTITIES_TYPES,
  ALARM_LIST_WIDGET_COLUMNS,
  CONTEXT_WIDGET_COLUMNS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
  ALARM_OUTPUT_FIELDS,
} from '@/constants';

import { formToWidgetColumn, widgetColumnValueToForm } from '@/helpers/entities/widget/column/form';

import { formBaseMixin, validationChildrenMixin } from '@/mixins/form';

import ColumnFieldExpandPanel from './column-field-expand-panel.vue';

export default {
  inject: ['$validator'],
  components: { ColumnFieldExpandPanel },
  mixins: [
    formBaseMixin,
    validationChildrenMixin,
  ],
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
  },
  data() {
    return {
      expanded: !this.column?.column,
      isCustom: false,
    };
  },
  computed: {
    isAlarmType() {
      return this.type === ENTITIES_TYPES.alarm;
    },

    alarmListAvailableColumns() {
      const columns = this.withInstructions
        ? ALARM_LIST_WIDGET_COLUMNS
        : omit(ALARM_LIST_WIDGET_COLUMNS, ['assignedInstructions']);

      return Object.values(columns).map(value => ({
        value,
        text: this.$tc(ALARM_FIELDS_TO_LABELS_KEYS[value], 2),
      }));
    },

    contextAvailableColumns() {
      return Object.values(CONTEXT_WIDGET_COLUMNS).map(value => ({
        value,
        text: this.$tc(ENTITY_FIELDS_TO_LABELS_KEYS[value], 2),
      }));
    },

    availableColumns() {
      return this.isAlarmType
        ? this.alarmListAvailableColumns
        : this.contextAvailableColumns;
    },

    columnLabelFieldName() {
      return `${this.name}.label`;
    },

    columnLabelErrorMessages() {
      return this.errors.collect(this.columnLabelFieldName);
    },
  },
  watch: {
    type() {
      const columns = this.columns.map(({ key }) => ({
        key,
        column: '',
      }));

      this.updateModel(columns);
    },

    availableColumns: {
      immediate: true,
      handler(columns) {
        this.isCustom = this.column.column
          ? columns.every(column => column.value !== this.column.column)
          : false;
      },
    },
  },
  methods: {
    changeColumn(column) {
      const newValue = {
        ...this.column,

        column,
      };

      if (this.withHtml) {
        newValue.isHtml = ALARM_OUTPUT_FIELDS.includes(column);
      }

      this.updateModel(newValue);
    },

    convertToCustom() {
      this.isCustom = !this.isCustom;

      const newColumn = {
        ...this.column,
      };

      if (this.isCustom) {
        const { value } = formToWidgetColumn(this.column);

        const selectedColumn = this.availableColumns.find(column => column.value === this.column.column);

        const label = this.column.label || selectedColumn?.text || '';

        newColumn.column = value;
        newColumn.label = label;
        newColumn.field = '';
        newColumn.rule = '';
        newColumn.dictionary = '';
      } else {
        const { column: value, field, rule, dictionary } = widgetColumnValueToForm(this.column.column);

        const selectedColumn = this.availableColumns.find(column => column.value === value);

        newColumn.column = value === selectedColumn?.value ? value : '';
        newColumn.label = newColumn.label === selectedColumn?.text ? '' : newColumn.label;
        newColumn.field = field;
        newColumn.rule = rule;
        newColumn.dictionary = dictionary;
      }

      this.updateModel(newColumn);
    },
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
