<template lang="pug">
  div
    c-column-field.my-2(
      v-for="(column, index) in columns",
      v-field="columns[index]",
      :key="column.key",
      :name="column.key",
      :type="type",
      :alarm-infos="alarmInfos",
      :entity-infos="entityInfos",
      :infos-pending="infosPending",
      :with-html="withHtml",
      :with-template="withTemplate",
      :with-color-indicator="withColorIndicator",
      :disabled-up="index === 0",
      :disabled-down="index === columns.length - 1",
      @up="up(index)",
      @down="down(index)",
      @remove="removeItemFromArray(index)"
    )
    v-btn.ml-0(color="primary", @click.prevent="add") {{ $t('common.add') }}
</template>

<script>
import {
  MODALS,
  ENTITIES_TYPES,
  COLOR_INDICATOR_TYPES,
  DEFAULT_COLUMN_TEMPLATE_VALUE,
  ALARM_INFOS_FIELDS,
  ENTITY_INFOS_FIELDS,
  ALARM_LIST_WIDGET_COLUMNS,
  CONTEXT_WIDGET_COLUMNS,
  ALARM_FIELDS_TO_LABELS_KEYS,
  ENTITY_FIELDS_TO_LABELS_KEYS,
} from '@/constants';

import { widgetColumnToForm } from '@/helpers/forms/shared/widget-column';

import { formArrayMixin, formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [
    formArrayMixin,
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    type: {
      type: String,
      default: ENTITIES_TYPES.alarm,
    },
    columns: {
      type: [Array, Object],
      default: () => [],
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
    alarmInfos: {
      type: Array,
      default: () => [],
    },
    entityInfos: {
      type: Array,
      default: () => [],
    },
    infosPending: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isAlarmType() {
      return this.type === ENTITIES_TYPES.alarm;
    },

    alarmListAvailableColumns() {
      return Object.values(ALARM_LIST_WIDGET_COLUMNS).map(value => ({
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

    infosFields() {
      return [
        ...ALARM_INFOS_FIELDS,
        ...ENTITY_INFOS_FIELDS,
      ];
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
  },
  methods: {
    isLinks(column) {
      return [
        ALARM_LIST_WIDGET_COLUMNS.links,
        CONTEXT_WIDGET_COLUMNS.links,
      ].includes(column);
    },

    isInfos(column) {
      return this.infosFields.includes(column);
    },

    getInfosByColumn(column) {
      return [
        ALARM_LIST_WIDGET_COLUMNS.entityInfos,
        ALARM_LIST_WIDGET_COLUMNS.entityComponentInfos,
        CONTEXT_WIDGET_COLUMNS.infos,
        CONTEXT_WIDGET_COLUMNS.componentInfos,
      ].includes(column) ? this.entityInfos : this.alarmInfos;
    },

    enableTemplate(index, checked) {
      const value = checked
        ? DEFAULT_COLUMN_TEMPLATE_VALUE
        : null;

      return this.updateFieldInArrayItem(index, 'template', value);
    },

    showEditTemplateModal(index) {
      const column = this.columns[index];

      this.$modals.show({
        name: MODALS.textEditor,
        config: {
          text: column.template,
          title: this.$t('settings.columns.withTemplate'),
          label: this.$t('common.template'),
          rules: {
            required: true,
          },
          action: value => this.updateFieldInArrayItem(index, 'template', value),
        },
      });
    },

    switchChangeColorIndicator(index, value) {
      return this.updateFieldInArrayItem(index, 'colorIndicator', value ? COLOR_INDICATOR_TYPES.state : null);
    },

    add() {
      this.addItemIntoArray(widgetColumnToForm());
    },

    up(index) {
      if (index > 0) {
        const columns = [...this.columns];
        const temp = columns[index];

        columns[index] = columns[index - 1];
        columns[index - 1] = temp;

        this.updateModel(columns);
      }
    },

    down(index) {
      if (index < this.columns.length - 1) {
        const columns = [...this.columns];
        const temp = columns[index];

        columns[index] = columns[index + 1];
        columns[index + 1] = temp;

        this.updateModel(columns);
      }
    },
  },
};
</script>
