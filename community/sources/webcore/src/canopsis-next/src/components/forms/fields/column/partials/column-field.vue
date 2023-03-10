<template lang="pug">
  v-card
    v-layout.pt-2(justify-space-between)
      v-flex(xs5)
        span.pl-2
          v-btn(
            :disabled="disabledUp",
            icon,
            @click.prevent="$emit('up')"
          )
            v-icon arrow_upward
          v-btn(
            :disabled="disabledDown",
            icon,
            @click.prevent="$emit('down')"
          )
            v-icon arrow_downward
      v-flex.text-xs-right.pr-2(xs3)
        v-btn(icon, @click.prevent="$emit('remove')")
          v-icon(color="error") close
    v-layout.px-3.pb-3(justify-center, column)
      v-select(
        v-field="column.column",
        v-validate="'required'",
        :items="availableColumns",
        :label="$tc('common.column', 1)",
        :error-messages="errors.collect(`${name}.column`)",
        :name="`${name}.column`"
      )
      c-alarm-infos-attribute-field(
        v-if="isAlarmInfos",
        v-field="column",
        :rules="alarmInfosRules",
        :pending="infosPending",
        :name="`${name}.column`"
      )
      c-infos-attribute-field(
        v-else-if="isInfos",
        v-field="column",
        :items="infosItems",
        :pending="infosPending",
        :name="`${name}.column`",
        combobox,
        column
      )
      template(v-if="isLinks")
        column-links-category-field(v-field="column.field")
        v-switch.pa-0.my-2(
          v-field="column.onlyIcon",
          :label="$t('settings.columns.onlyIcon')",
          color="primary",
          hide-details
        )
      v-switch.pa-0.my-2(
        v-model="customLabel",
        :label="$t('settings.columns.customLabel')",
        color="primary",
        hide-details,
        @change="updateCustomLabel"
      )
      v-text-field(
        v-if="customLabel",
        v-field="column.label",
        v-validate="'required'",
        :label="$t('common.label')",
        :error-messages="errors.collect(`${name}.label`)",
        :name="`${name}.label`"
      )
      v-layout(v-if="withTemplate", row, align-center)
        v-switch.pa-0.my-2(
          :label="$t('settings.columns.withTemplate')",
          :input-value="!!column.template",
          color="primary",
          hide-details,
          @change="enableTemplate($event)"
        )
        v-btn.primary(
          v-if="column.template",
          small,
          @click="showEditTemplateModal"
        )
          span {{ $t('common.edit') }}
      v-switch.pa-0.my-2(
        v-if="withHtml",
        v-field="column.isHtml",
        :label="$t('settings.columns.isHtml')",
        :disabled="!!column.template",
        color="primary",
        hide-details
      )
      v-switch.pa-0.my-2(
        v-if="withColorIndicator",
        :label="$t('settings.colorIndicator.title')",
        :input-value="!!column.colorIndicator",
        :disabled="!!column.template",
        color="primary",
        hide-details,
        @change="switchChangeColorIndicator($event)"
      )
      v-layout(v-if="column.colorIndicator", row)
        c-color-indicator-field(
          v-field="column.colorIndicator",
          :disabled="!!column.template"
        )
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
  ALARM_FIELDS,
} from '@/constants';

import { isLinksWidgetColumn } from '@/helpers/forms/shared/widget-column';

import { formMixin } from '@/mixins/form';
import { entitiesInfosMixin } from '@/mixins/entities/infos';

import ColumnLinksCategoryField from './column-links-category-field.vue';

export default {
  inject: ['$validator'],
  components: { ColumnLinksCategoryField },
  mixins: [
    formMixin,
    entitiesInfosMixin,
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
    name: {
      type: String,
      default: '',
    },
    disabledUp: {
      type: Boolean,
      default: false,
    },
    disabledDown: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      customLabel: !!this.column.label,
    };
  },
  computed: {
    infosFields() {
      return [
        ...ALARM_INFOS_FIELDS,
        ...ENTITY_INFOS_FIELDS,
      ];
    },

    isAlarmType() {
      return this.type === ENTITIES_TYPES.alarm;
    },

    isLinks() {
      return isLinksWidgetColumn(this.column?.column);
    },

    isAlarmInfos() {
      return ALARM_FIELDS.infos === this.column?.column;
    },

    isInfos() {
      return this.infosFields.includes(this.column?.column);
    },

    infosItems() {
      return [
        ALARM_LIST_WIDGET_COLUMNS.entityInfos,
        ALARM_LIST_WIDGET_COLUMNS.entityComponentInfos,
        CONTEXT_WIDGET_COLUMNS.infos,
        CONTEXT_WIDGET_COLUMNS.componentInfos,
      ].includes(this.column?.column) ? this.entityInfos : this.alarmInfos;
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
    updateCustomLabel(checked) {
      if (checked) {
        return;
      }

      this.updateField('label', '');
    },

    enableTemplate(checked) {
      const template = checked
        ? DEFAULT_COLUMN_TEMPLATE_VALUE
        : null;

      return this.updateModel({
        ...this.column,

        template,
        isHtml: checked && this.column.isHtml ? false : this.column.isHtml,
        colorIndicator: checked && this.column.isHtml ? null : this.column.isHtml,
      });
    },

    switchChangeColorIndicator(colorIndicator) {
      const value = colorIndicator
        ? COLOR_INDICATOR_TYPES.state
        : null;

      return this.updateField('colorIndicator', value);
    },

    showEditTemplateModal() {
      this.$modals.show({
        name: MODALS.textEditor,
        config: {
          text: this.column?.template ?? '',
          title: this.$t('settings.columns.withTemplate'),
          label: this.$t('common.template'),
          rules: {
            required: true,
          },
          action: value => this.updateField('template', value),
        },
      });
    },
  },
};
</script>
