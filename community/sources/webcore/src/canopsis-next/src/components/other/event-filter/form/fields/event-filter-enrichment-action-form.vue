<template lang="pug">
  v-card
    v-card-text
      v-layout(row, align-start)
        v-icon.draggable.ml-0.mr-3.mt-3.action-drag-handler drag_indicator
        v-layout(column)
          v-layout(row)
            v-select(
              v-field="form.type",
              :items="eventFilterActionTypes",
              :label="$t('common.type')"
            )
            v-btn.mr-0(icon, @click="remove")
              v-icon(color="error") delete
          v-expand-transition
            event-filter-enrichment-action-form-type-info(v-if="form.type", :type="form.type")
          v-text-field(
            v-field="form.description",
            :label="$t('common.description')",
            key="description"
          )
          v-layout(v-if="isStringDictionaryValueType")
            c-payload-text-field(
              v-field="form.value",
              :label="$t('common.value')",
              :variables="variables",
              :name="valueFieldName",
              key="value",
              required,
              clearable
            )
          v-layout(v-else)
            v-flex(xs5)
              c-name-field(v-field="form.name", key="name", required)
            v-flex(xs7)
              c-payload-text-field.ml-2(
                v-if="isStringTemplateValueType",
                v-field="form.value",
                :label="$t('common.value')",
                :variables="variables",
                :name="valueFieldName",
                key="from",
                required,
                clearable
              )
              v-combobox.ml-2(
                v-else-if="isStringCopyValueType",
                v-field="form.value",
                v-validate="'required'",
                :label="$t('common.value')",
                :error-messages="errors.collect('value')",
                :items="copyValueVariables",
                :name="valueFieldName",
                key="from"
              )
              c-mixed-field.ml-2(
                v-else,
                v-field="form.value",
                :label="$t('common.value')",
                :name="valueFieldName",
                key="value"
              )
</template>

<script>
import { ACTION_COPY_PAYLOAD_VARIABLES, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES } from '@/constants';

import EventFilterEnrichmentActionFormTypeInfo from './event-filter-enrichment-action-form-type-info.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterEnrichmentActionFormTypeInfo },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    variables: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'action',
    },
  },
  computed: {
    valueFieldName() {
      return `${this.name}.value`;
    },

    eventFilterActionTypes() {
      return Object.values(EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES).map(value => ({
        value,

        text: this.$t(`eventFilter.actionsTypes.${value}.text`),
      }));
    },

    copyValueVariables() {
      return Object.values(ACTION_COPY_PAYLOAD_VARIABLES);
    },

    isStringCopyValueType() {
      return [
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo,
      ].includes(this.form.type);
    },

    isStringTemplateValueType() {
      return [
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,
      ].includes(this.form.type);
    },

    isStringDictionaryValueType() {
      return EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary === this.form.type;
    },
  },
  watch: {
    'form.type': function typeWatcher() {
      this.errors.clear();
    },
  },
  methods: {
    remove() {
      this.$emit('remove', this.form);
    },
  },
};
</script>
