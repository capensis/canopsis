<template>
  <c-layout justify="center" column>
    <c-text-field
      v-if="withField"
      :model-value="formValue.column"
      :label="$t('common.field')"
      :name="columnValueFieldName"
      required
      @update:model-value="updateField('column', $event)"
    />

    <c-text-field
      v-if="isAlarmTicketData"
      :model-value="formValue.dictionary"
      :label="$t('common.dictionary')"
      :name="columnDictionaryName"
      required
      @update:model-value="updateField('dictionary', $event)"
    />

    <template v-if="!withoutInfosAttributes">
      <c-alarm-infos-attribute-field
        v-if="isAlarmInfos"
        :value="formValue"
        :rules="alarmInfosRules"
        :pending="infosPending"
        :name="`${name}.column`"
        :required="!optionalInfosAttributes"
        @update:value="updateModel"
      />
      <c-infos-attribute-field
        v-else-if="isInfos"
        :value="formValue"
        :items="infosItems"
        :pending="infosPending"
        :name="`${name}.column`"
        combobox
        column
        @update:value="updateModel"
      />
    </template>

    <template v-if="isLinks">
      <column-links-category-field
        :model-value="formValue.field"
        @update:model-value="updateField('field', $event)"
      />
      <c-number-field
        :value="formValue.inlineLinksCount"
        :label="$t('settings.columns.inlineLinksCount')"
        @update:value="updateField('inlineLinksCount', $event)"
      />
      <v-switch
        :model-value="formValue.onlyIcon"
        :label="$t('settings.columns.onlyIcon')"
        class="pa-0 my-0"
        color="primary"
        hide-details
        @update:model-value="updateField('onlyIcon', $event)"
      />
      <c-number-field
        v-if="formValue.onlyIcon"
        :value="formValue.linksInRowCount"
        :label="$t('settings.columns.linksInRowCount')"
        @update:value="updateField('linksInRowCount', $event)"
      >
        <template #append-inner>
          <c-help-icon :text="$t('settings.columns.linksInRowCountTooltip')" />
        </template>
      </c-number-field>
    </template>

    <template v-if="withLabel">
      <v-switch
        v-model="customLabel"
        :label="$t('settings.columns.customLabel')"
        class="pa-0 my-0"
        color="primary"
        hide-details
        @update:model-value="updateCustomLabel"
      />
      <c-text-field
        v-if="customLabel"
        :model-value="formValue.label"
        :label="$t('common.label')"
        :name="`${name}.label`"
        required
        @update:model-value="updateField('label', $event)"
      />
    </template>

    <c-layout
      v-if="withTemplate || withSimpleTemplate"
      align="center"
      justify="space-between"
    >
      <v-switch
        :model-value="!!formValue.template"
        class="pa-0 my-0"
        color="primary"
        hide-details
        @update:model-value="switchChangeTemplate"
      >
        <template #label>
          <span @click.prevent.stop="switchChangeTemplate(!formValue.template)">
            {{ $t('settings.columns.withTemplate') }}
          </span>
        </template>
      </v-switch>
      <v-btn
        v-if="formValue.template"
        color="primary"
        size="small"
        @click="showEditTemplateModal"
      >
        <span>{{ $t('common.edit') }}</span>
      </v-btn>
    </c-layout>

    <v-switch
      v-if="withHtml"
      :model-value="formValue.isHtml"
      :label="$t('settings.columns.isHtml')"
      :disabled="!!formValue.template"
      class="pa-0 my-0"
      color="primary"
      hide-details
      @update:model-value="switchChangeIsHtml"
    />

    <v-switch
      v-if="withColorIndicator"
      :label="$t('settings.colorIndicator.title')"
      :model-value="!!formValue.colorIndicator"
      :disabled="!!formValue.template"
      class="pa-0 my-0"
      color="primary"
      hide-details
      @update:model-value="switchChangeColorIndicator"
    />

    <c-color-indicator-field
      v-if="formValue.colorIndicator"
      :model-value="formValue.colorIndicator"
      :disabled="!!formValue.template"
      @update:model-value="updateField('colorIndicator', $event)"
    />

    <v-switch
      v-if="filterOnClickAvailable"
      :model-value="formValue.isFilter"
      :label="$t('settings.columns.filterOnClick')"
      :disabled="!!formValue.template || !!formValue.isHtml"
      class="pa-0 my-0"
      color="primary"
      hide-details
      @update:model-value="updateField('isFilter', $event)"
    />
  </c-layout>
</template>

<script>
import {
  ALARM_FIELDS,
  ALARM_FIELDS_WITHOUT_FILTER_ON_CLICK,
  ALARM_INFOS_FIELDS,
  ALARM_LIST_WIDGET_COLUMNS,
  COLOR_INDICATOR_TYPES,
  CONTEXT_WIDGET_COLUMNS,
  DEFAULT_COLUMN_TEMPLATE_VALUE,
  ENTITY_INFOS_FIELDS,
  MODALS,
} from '@/constants';

import { isLinksWidgetColumn } from '@/helpers/entities/widget/column/form';

import { useInfos } from '@/hooks/store/modules/infos';
import { formMixin } from '@/mixins/form';

import ColumnLinksCategoryField from './column-links-category-field.vue';

export default {
  components: { ColumnLinksCategoryField },
  mixins: [formMixin],

  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    column: {
      type: Object,
      default: undefined,
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
    withoutInfosAttributes: {
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

  emits: ['update:value'],

  setup() {
    const {
      alarmInfos,
      alarmInfosRules,
      entityInfos,
      infosPending,
    } = useInfos();

    return {
      alarmInfos,
      alarmInfosRules,
      entityInfos,
      infosPending,
    };
  },

  data() {
    return {
      customLabel: false,
    };
  },

  computed: {
    formValue() {
      return this.column ?? this.value;
    },

    infosFields() {
      return [
        ...ALARM_INFOS_FIELDS,
        ...ENTITY_INFOS_FIELDS,
      ];
    },

    isLinks() {
      return isLinksWidgetColumn(this.formValue?.column);
    },

    isAlarmInfos() {
      return ALARM_FIELDS.infos === this.formValue?.column;
    },

    isInfos() {
      return this.infosFields.includes(this.formValue?.column);
    },

    isAlarmTicketData() {
      return ALARM_FIELDS.ticketData === this.formValue?.column;
    },

    infosItems() {
      return [
        ALARM_LIST_WIDGET_COLUMNS.entityInfos,
        ALARM_LIST_WIDGET_COLUMNS.entityComponentInfos,
        CONTEXT_WIDGET_COLUMNS.infos,
        CONTEXT_WIDGET_COLUMNS.componentInfos,
      ].includes(this.formValue?.column) ? this.entityInfos : this.alarmInfos;
    },

    templateModalConfig() {
      return {
        text: this.formValue?.template ?? '',
        title: this.$t('settings.columns.withTemplate'),
        label: this.$t('common.template', 1),
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

    columnDictionaryName() {
      return `${this.name}.dictionary`;
    },

    filterOnClickAvailable() {
      return this.withFilterOnClick && !ALARM_FIELDS_WITHOUT_FILTER_ON_CLICK.includes(this.formValue?.column);
    },
  },

  watch: {
    withLabel: {
      immediate: true,
      handler() {
        this.customLabel = !!this.formValue.label;
      },
    },
  },

  methods: {
    updateModel(model) {
      this.$emit('update:value', model);

      return model;
    },

    updateCustomLabel(checked) {
      if (checked) {
        return;
      }

      this.updateField('label', '');
    },

    updateModelByTemplate(checked, template) {
      return this.updateModel({
        ...this.formValue,
        template,
        isHtml: checked && this.formValue.isHtml ? false : this.formValue.isHtml,
        colorIndicator: checked && this.formValue.colorIndicator ? null : this.formValue.colorIndicator,
        isFilter: checked && this.formValue.isFilter ? false : this.formValue.isFilter,
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
      return this.updateField('colorIndicator', colorIndicator ? COLOR_INDICATOR_TYPES.state : null);
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

    switchChangeIsHtml(isHtml) {
      this.updateModel({
        ...this.formValue,
        isHtml,
        isFilter: isHtml ? false : this.formValue.isFilter,
      });
    },
  },
};
</script>
