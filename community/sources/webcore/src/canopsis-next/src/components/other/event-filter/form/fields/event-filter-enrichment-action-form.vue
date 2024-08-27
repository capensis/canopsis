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
              :value="form.type"
              :items="eventFilterActionTypes"
              :label="$t('common.type')"
              @change="changeActionType"
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
            key="description"
            :label="$t('common.description')"
          />
          <v-layout v-if="isStringDictionaryValueType">
            <v-text-field
              v-field="form.value"
              v-validate="'required'"
              key="value"
              :label="$t('common.value')"
              :name="valueFieldName"
              :error-messages="errors.collect(valueFieldName)"
              :prefix="eventExtraPrefix"
              clearable
            />
          </v-layout>
          <v-layout v-else>
            <v-flex xs5>
              <c-name-field
                v-field="form.name"
                key="name"
                :name="nameFieldName"
                class="mr-2"
                required
              />
            </v-flex>
            <v-flex xs7>
              <c-payload-text-field
                v-if="isStringTemplateValueType"
                v-field="form.value"
                key="from"
                :label="$t('common.value')"
                :variables="variables"
                :name="valueFieldName"
                class="ml-2"
                required
                clearable
              />
              <v-combobox
                v-else-if="isStringCopyValueType"
                v-field="form.value"
                v-validate="'required'"
                key="value"
                :label="$t('common.value')"
                :error-messages="errors.collect(valueFieldName)"
                :items="copyValueVariables"
                :name="valueFieldName"
                class="ml-2"
              />
              <event-filter-enrichment-action-form-select-rags-value
                v-else-if="isSelectValueType"
                v-field="form.value"
                key="value"
                :items="setTagsItems"
                :name="valueFieldName"
              />
              <c-mixed-field
                v-else
                v-field="form.value"
                key="value"
                :label="$t('common.value')"
                :name="valueFieldName"
                class="ml-2"
              />
            </v-flex>
          </v-layout>
        </v-layout>
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import {
  ACTION_COPY_PAYLOAD_VARIABLES,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EVENT_FILTER_EVENT_EXTRA_PREFIX,
} from '@/constants';

import {
  eventFilterDictionaryActionValueToForm,
  formToEventFilterDictionaryActionValue,
} from '@/helpers/entities/event-filter/rule/form';

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
    eventExtraPrefix() {
      return EVENT_FILTER_EVENT_EXTRA_PREFIX;
    },
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

    isStringDictionaryValueType() {
      return EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary === this.form.type;
    },

    isSelectValueType() {
      return EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTags === this.form.type;
    },
  },
  methods: {
    changeActionType(type) {
      const newForm = {
        ...this.form,

        type,
      };

      if (this.form.type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary) {
        newForm.value = formToEventFilterDictionaryActionValue(this.form.value);
      } else if (type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary) {
        newForm.value = eventFilterDictionaryActionValueToForm(this.form.value);
      }

      this.updateModel(newForm);
      this.errors.clear();
    },

    remove() {
      this.$emit('remove', this.form);
    },
  },
};
</script>
