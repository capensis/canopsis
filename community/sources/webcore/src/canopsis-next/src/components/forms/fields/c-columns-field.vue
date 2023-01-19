<template lang="pug">
  div
    v-card.my-2(
      v-for="(column, index) in columns",
      :key="column.key"
    )
      v-layout.pt-2(justify-space-between)
        v-flex(xs3)
          v-layout.text-xs-center.pl-2(justify-space-between)
            v-btn(
              :disabled="index === 0",
              icon,
              @click.prevent="up(index)"
            )
              v-icon arrow_upward
            v-btn(
              :disabled="index === columns.length - 1",
              icon,
              @click.prevent="down(index)"
            )
              v-icon arrow_downward
        v-flex.text-xs-right.pr-2(xs3)
          v-btn(icon, @click.prevent="removeItemFromArray(index)")
            v-icon(color="red") close
      v-layout.px-3.pb-3(justify-center, column)
        v-select(
          v-field="columns[index].column",
          v-validate="'required'",
          :items="availableColumns",
          :label="$tc('common.column', 1)",
          :error-messages="errors.collect(`${column.key}.column`)",
          :name="`${column.key}.column`"
        )
        c-infos-attribute-field(
          v-if="isInfos(column.column)",
          v-field="columns[index]",
          :items="getInfosByColumn(column.column)",
          :pending="infosPending",
          :name="`${column.key}.column`",
          combobox,
          column
        )
        v-text-field(
          v-if="isLinks(column.column)",
          v-field="columns[index].field",
          :label="$t('common.field')"
        )
        v-layout(v-if="withTemplate", row)
          v-switch(
            :label="$t('settings.columns.withTemplate')",
            :input-value="!!column.template",
            color="primary",
            @change="enableTemplate(index, $event)"
          )
          v-btn.primary(v-if="column.template", small, @click="showEditTemplateModal(index)")
            span {{ $t('common.edit') }}
        v-switch(
          v-if="withHtml",
          v-field="columns[index].isHtml",
          :label="$t('settings.columns.isHtml')",
          :disabled="!!column.template",
          color="primary"
        )
        v-switch(
          v-if="withColorIndicator",
          :label="$t('settings.colorIndicator.title')",
          :input-value="!!column.colorIndicator",
          :disabled="!!column.template",
          color="primary",
          @change="switchChangeColorIndicator(index, $event)"
        )
        v-layout(v-if="column.colorIndicator", row)
          c-color-indicator-field(
            v-field="columns[index].colorIndicator",
            :disabled="!!column.template"
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
