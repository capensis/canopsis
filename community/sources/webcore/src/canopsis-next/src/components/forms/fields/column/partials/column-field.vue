<template lang="pug">
  v-card
    v-card-text
      v-layout(row, align-center)
        span.handler.mr-1
          v-icon.draggable(:class="dragHandleClass") drag_indicator
        c-expand-btn.mr-1(
          v-model="expanded",
          :color="hasChildrenError ? 'error' : ''"
        )
        v-select(
          v-validate="'required'",
          :value="column.column",
          :items="availableColumns",
          :label="$tc('common.column', 1)",
          :error-messages="errors.collect(`${name}.column`)",
          :name="`${name}.column`",
          @change="changeColumn"
        )
        v-tooltip(left)
          template(#activator="{ on }")
            v-btn.mr-0(
              v-on="on",
              small,
              flat,
              icon,
              @click="$emit('remove')"
            )
              v-icon(color="error", small) close
          span {{ $t('common.delete') }}
      v-expand-transition(mode="out-in")
        column-field-expand-panel.pl-1(
          v-show="expanded",
          v-field="column",
          :name="name",
          :with-html="withHtml",
          :with-template="withTemplate",
          :with-color-indicator="withColorIndicator",
          :with-instructions="withInstructions",
          :with-simple-template="withSimpleTemplate",
          :optional-infos-attributes="optionalInfosAttributes",
          :without-infos-attributes="withoutInfosAttributes",
          :variables="variables"
        )
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
    withAlarmDisplayName: {
      type: Boolean,
      default: false,
    },
    withAlarmCreationDate: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      expanded: !this.column?.column,
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
      let columns = Object.values(CONTEXT_WIDGET_COLUMNS).map(value => ({
        value,
        text: this.$tc(ENTITY_FIELDS_TO_LABELS_KEYS[value], 2),
      }));

      if (!this.withAlarmDisplayName) {
        columns = columns.filter(column => column.value !== CONTEXT_WIDGET_COLUMNS.alarmDisplayName);
      }

      if (!this.withAlarmCreationDate) {
        columns = columns.filter(column => column.value !== CONTEXT_WIDGET_COLUMNS.alarmCreationDate);
      }

      return columns;
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
  },
};
</script>
