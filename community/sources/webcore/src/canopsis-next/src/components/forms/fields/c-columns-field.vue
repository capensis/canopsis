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
          :label="$t('common.column')",
          :error-messages="errors.collect(`${column.key}.column`)",
          :name="`${column.key}.column`",
          :return-object="false"
        )
        c-infos-attribute-field(
          v-if="hasDictionary(column.column)",
          v-field="columns[index]",
          :name="`${column.key}.column`",
          combobox,
          column
        )
        v-text-field(
          v-if="hasField(column.column)",
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
  COLOR_INDICATOR_TYPES,
  DEFAULT_COLUMN_TEMPLATE_VALUE,
  WIDGET_TYPES,
  ALARM_LIST_WIDGET_COLUMNS,
  CONTEXT_WIDGET_COLUMNS,
  ALARM_LIST_WIDGET_COLUMNS_TO_LABELS_KEYS,
  CONTEXT_WIDGET_COLUMNS_TO_LABELS_KEYS,
} from '@/constants';

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
      default: WIDGET_TYPES.alarmList,
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
    infosDictionary: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    alarmListAvailableColumns() {
      return Object.values(ALARM_LIST_WIDGET_COLUMNS).map(value => ({
        value,
        text: this.$t(ALARM_LIST_WIDGET_COLUMNS_TO_LABELS_KEYS[value]),
      }));
    },

    contextAvailableColumns() {
      return Object.values(CONTEXT_WIDGET_COLUMNS).map(value => ({
        value,
        text: this.$t(CONTEXT_WIDGET_COLUMNS_TO_LABELS_KEYS[value]),
      }));
    },

    availableColumns() {
      return this.type === WIDGET_TYPES.alarmList
        ? this.alarmListAvailableColumns
        : this.contextAvailableColumns;
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
    hasField(column) {
      return [
        ALARM_LIST_WIDGET_COLUMNS.links,
        CONTEXT_WIDGET_COLUMNS.links,
      ].includes(column);
    },

    hasDictionary(column) {
      return [
        ALARM_LIST_WIDGET_COLUMNS.infos,
        ALARM_LIST_WIDGET_COLUMNS.entityInfos,
        ALARM_LIST_WIDGET_COLUMNS.entityComponentInfos,
        CONTEXT_WIDGET_COLUMNS.infos,
        CONTEXT_WIDGET_COLUMNS.componentInfos,
      ].includes(column);
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
      const column = { label: '', value: '' };

      if (this.withHtml) {
        column.isHtml = false;
      }

      this.addItemIntoArray(column);
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
