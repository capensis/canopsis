<template>
  <v-card>
    <v-card-text>
      <v-layout align-start>
        <v-icon class="draggable ml-0 mr-3 mt-3 action-drag-handler">
          drag_indicator
        </v-icon>
        <v-layout column>
          <v-layout>
            <v-select
              v-field="form.type"
              :items="eventFilterActionTypes"
              :label="$t('common.type')"
            />
            <v-btn
              class="mr-0"
              icon
              @click="remove"
            >
              <v-icon color="error">
                delete
              </v-icon>
            </v-btn>
          </v-layout>
          <v-expand-transition>
            <event-filter-enrichment-action-form-type-info
              v-if="form.type"
              :type="form.type"
            />
          </v-expand-transition>
          <v-text-field
            v-field="form.description"
            :label="$t('common.description')"
            key="description"
          />
          <v-layout>
            <v-flex xs5>
              <c-name-field
                class="mr-2"
                v-field="form.name"
                :name="nameFieldName"
                key="name"
                required
              />
            </v-flex>
            <v-flex xs7>
              <c-payload-text-field
                v-if="isStringTemplateValueType"
                v-field="form.value"
                :label="$t('common.value')"
                :variables="variables"
                :name="valueFieldName"
                key="value"
                required
                clearable
              />
              <v-combobox
                v-else-if="isStringCopyValueType"
                v-field="form.value"
                v-validate="'required'"
                :label="$t('common.value')"
                :error-messages="errors.collect(valueFieldName)"
                :items="copyValueVariables"
                :name="valueFieldName"
                key="value"
              />
              <event-filter-enrichment-action-form-select-rags-value
                v-else-if="isSelectValueType"
                v-field="form.value"
                :items="setTagsItems"
                :name="valueFieldName"
                key="value"
              />
              <c-mixed-field
                v-else
                v-field="form.value"
                :label="$t('common.value')"
                :name="valueFieldName"
                key="value"
              />
            </v-flex>
          </v-layout>
        </v-layout>
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { ACTION_COPY_PAYLOAD_VARIABLES, EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

import EventFilterEnrichmentActionFormTypeInfo from './event-filter-enrichment-action-form-type-info.vue';
import EventFilterEnrichmentActionFormSelectRagsValue from './event-filter-enrichment-action-form-select-tags-value.vue';

export default {
  inject: ['$validator'],
  components: {
    EventFilterEnrichmentActionFormTypeInfo,
    EventFilterEnrichmentActionFormSelectRagsValue,
  },
  mixins: [formMixin],
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
    setTagsItems: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    nameFieldName() {
      return `${this.name}.name`;
    },

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
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTagsFromTemplate,
      ].includes(this.form.type);
    },

    isSelectValueType() {
      return EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTags === this.form.type;
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
