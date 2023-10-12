<template lang="pug">
  v-layout(justify-center, column)
    template(v-if="!withoutInfosAttributes")
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
      c-number-field(
        v-field="column.inlineLinksCount",
        :label="$t('settings.columns.inlineLinksCount')"
      )
      v-switch.pa-0.my-2(
        v-field="column.onlyIcon",
        :label="$t('settings.columns.onlyIcon')",
        color="primary",
        hide-details
      )
      c-number-field(
        v-if="column.onlyIcon",
        v-field="column.linksInRowCount",
        :label="$t('settings.columns.linksInRowCount')"
      )
        template(#append="")
          c-help-icon(:text="$t('settings.columns.linksInRowCountTooltip')", left)
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
        @change="switchChangeTemplate($event)"
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
  COLOR_INDICATOR_TYPES,
  DEFAULT_COLUMN_TEMPLATE_VALUE,
  ALARM_INFOS_FIELDS,
  ENTITY_INFOS_FIELDS,
  ALARM_LIST_WIDGET_COLUMNS,
  CONTEXT_WIDGET_COLUMNS,
  ALARM_FIELDS,
} from '@/constants';

import { isLinksWidgetColumn } from '@/helpers/entities/widget/column/form';

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
    column: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: '',
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
  },
  methods: {
    updateCustomLabel(checked) {
      if (checked) {
        return;
      }

      this.updateField('label', '');
    },

    switchChangeTemplate(checked) {
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