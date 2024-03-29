<template>
  <v-layout
    justify-center
    column
  >
    <v-text-field
      v-if="withField"
      v-validate="'required'"
      v-field="column.column"
      :label="$t('common.field')"
      :error-messages="columnValueErrorMessages"
      :name="columnValueFieldName"
    />
    <template v-if="!withoutInfosAttributes">
      <c-alarm-infos-attribute-field
        v-if="isAlarmInfos"
        v-field="column"
        :rules="alarmInfosRules"
        :pending="infosPending"
        :name="`${name}.column`"
        :required="!optionalInfosAttributes"
      />
      <c-infos-attribute-field
        v-else-if="isInfos"
        v-field="column"
        :items="infosItems"
        :pending="infosPending"
        :name="`${name}.column`"
        combobox
        column
      />
    </template>
    <template v-if="isLinks">
      <column-links-category-field v-field="column.field" />
      <c-number-field
        v-field="column.inlineLinksCount"
        :label="$t('settings.columns.inlineLinksCount')"
      />
      <v-switch
        v-field="column.onlyIcon"
        :label="$t('settings.columns.onlyIcon')"
        class="pa-0 my-2"
        color="primary"
        hide-details
      />
      <c-number-field
        v-if="column.onlyIcon"
        v-field="column.linksInRowCount"
        :label="$t('settings.columns.linksInRowCount')"
      >
        <template #append="">
          <c-help-icon
            :text="$t('settings.columns.linksInRowCountTooltip')"
            left
          />
        </template>
      </c-number-field>
    </template>
    <template v-if="withLabel">
      <v-switch
        v-model="customLabel"
        :label="$t('settings.columns.customLabel')"
        class="pa-0 my-2"
        color="primary"
        hide-details
        @change="updateCustomLabel"
      />
      <v-text-field
        v-if="customLabel"
        v-field="column.label"
        v-validate="'required'"
        :label="$t('common.label')"
        :error-messages="errors.collect(`${name}.label`)"
        :name="`${name}.label`"
      />
    </template>
    <v-layout
      v-if="withTemplate || withSimpleTemplate"
      align-center
      justify-space-between
    >
      <v-switch
        :label="$t('settings.columns.withTemplate')"
        :input-value="!!column.template"
        :true-value="true"
        :false-value="false"
        :value-comparator="isCustomTemplate"
        class="pa-0 my-2"
        color="primary"
        hide-details
        @change="switchChangeTemplate($event)"
      />
      <v-btn
        v-if="column.template"
        class="primary"
        small
        @click="showEditTemplateModal"
      >
        <span>{{ $t('common.edit') }}</span>
      </v-btn>
    </v-layout>
    <v-switch
      v-if="withHtml"
      v-field="column.isHtml"
      :label="$t('settings.columns.isHtml')"
      :disabled="!!column.template"
      class="pa-0 my-2"
      color="primary"
      hide-details
    />
    <v-switch
      v-if="withColorIndicator"
      :label="$t('settings.colorIndicator.title')"
      :input-value="!!column.colorIndicator"
      :disabled="!!column.template"
      class="pa-0 my-2"
      color="primary"
      hide-details
      @change="switchChangeColorIndicator($event)"
    />
    <v-layout v-if="column.colorIndicator">
      <c-color-indicator-field
        v-field="column.colorIndicator"
        :disabled="!!column.template"
      />
    </v-layout>
  </v-layout>
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
    withSimpleTemplate: {
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
    variables: {
      type: Array,
      required: false,
    },
    withLabel: {
      type: Boolean,
      default: false,
    },
    withField: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      customLabel: false,
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

    templateModalConfig() {
      return {
        text: this.column?.template ?? '',
        title: this.$t('settings.columns.withTemplate'),
        label: this.$t('common.template'),
        variables: this.variables,
        rules: {
          required: true,
        },
      };
    },

    templateModalName() {
      return this.withSimpleTemplate ? MODALS.payloadTextareaEditor : MODALS.textEditor;
    },

    columnValueFieldName() {
      return `${this.name}.column`;
    },

    columnValueErrorMessages() {
      return this.errors.collect(this.columnValueFieldName);
    },
  },
  watch: {
    withLabel: {
      immediate: true,
      handler() {
        this.customLabel = !!this.column.label;
      },
    },
  },
  methods: {
    updateCustomLabel(checked) {
      if (checked) {
        return;
      }

      this.updateField('label', '');
    },

    updateModelByTemplate(checked, template) {
      return this.updateModel({
        ...this.column,

        template,
        isHtml: checked && this.column.isHtml ? false : this.column.isHtml,
        colorIndicator: checked && this.column.colorIndicator ? null : this.column.colorIndicator,
      });
    },

    switchChangeTemplate(checked) {
      if (!checked) {
        this.updateModelByTemplate(checked, null);

        return;
      }

      this.$modals.show({
        name: this.templateModalName,
        config: {
          ...this.templateModalConfig,
          text: this.withSimpleTemplate ? '' : DEFAULT_COLUMN_TEMPLATE_VALUE,
          action: value => this.updateModelByTemplate(checked, value),
        },
      });
    },

    switchChangeColorIndicator(colorIndicator) {
      const value = colorIndicator
        ? COLOR_INDICATOR_TYPES.state
        : null;

      return this.updateField('colorIndicator', value);
    },

    isCustomTemplate() {
      return !!this.column.template;
    },

    showEditTemplateModal() {
      this.$modals.show({
        name: this.templateModalName,
        config: {
          ...this.templateModalConfig,
          action: value => this.updateField('template', value),
        },
      });
    },
  },
};
</script>
