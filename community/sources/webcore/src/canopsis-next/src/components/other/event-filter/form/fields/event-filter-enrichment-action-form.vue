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
                v-validate="'required'"
                key="value"
                :value="form.value"
                :label="$t('common.value')"
                :error-messages="errors.collect(valueFieldName)"
                :items="copyVariables"
                :name="valueFieldName"
                class="ml-2"
                return-object
                @input="updateCopyValue"
              />
              <event-filter-enrichment-action-form-select-tags-value
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
import { computed } from 'vue';

import { EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES, EVENT_FILTER_EVENT_EXTRA_PREFIX } from '@/constants';

import {
  eventFilterDictionaryActionValueToForm,
  formToEventFilterDictionaryActionValue,
} from '@/helpers/entities/event-filter/rule/form';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';
import { useValidator } from '@/hooks/validator/validator';

import EventFilterEnrichmentActionFormTypeInfo from './event-filter-enrichment-action-form-type-info.vue';
import EventFilterEnrichmentActionFormSelectTagsValue from './event-filter-enrichment-action-form-select-tags-value.vue';

export default {
  inject: ['$validator'],
  components: {
    EventFilterEnrichmentActionFormTypeInfo,
    EventFilterEnrichmentActionFormSelectTagsValue,
  },
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
    copyVariables: {
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const validator = useValidator();
    const { updateField, updateModel } = useModelField(props, emit);

    const eventExtraPrefix = EVENT_FILTER_EVENT_EXTRA_PREFIX;

    const nameFieldName = computed(() => `${props.name}.name`);
    const valueFieldName = computed(() => `${props.name}.value`);

    const eventFilterActionTypes = computed(() => Object.values(EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES).map(value => ({
      value,
      text: t(`eventFilter.actionsTypes.${value}.text`),
    })));

    const isStringCopyValueType = computed(() => [
      EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy,
      EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo,
    ].includes(props.form.type));

    const isStringTemplateValueType = computed(() => [
      EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate,
      EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,
      EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTagsFromTemplate,
    ].includes(props.form.type));

    const isStringDictionaryValueType = computed(() => (
      EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary === props.form.type
    ));

    const isSelectValueType = computed(() => EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTags === props.form.type);

    /**
     * Changes the action type and handles value transformation for dictionary actions
     *
     * @param {string} type - The new action type from EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES
     */
    const changeActionType = (type) => {
      const newForm = {
        ...props.form,
        type,
      };

      if (props.form.type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary) {
        newForm.value = formToEventFilterDictionaryActionValue(props.form.value);
      } else if (type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary) {
        newForm.value = eventFilterDictionaryActionValueToForm(props.form.value);
      }

      updateModel(newForm);

      validator.errors.clear();
    };

    /**
     * Updates the copy value field by calling updateField with 'value' key
     *
     * @param {Object} [params={}] - Parameters object
     * @param {string} [params.value=''] - The value to set for the copy field
     */
    const updateCopyValue = ({ value = '' } = {}) => updateField('value', value);

    /**
     * Emits remove event to parent component to remove this action form
     */
    const remove = () => emit('remove', props.form);

    return {
      eventExtraPrefix,

      nameFieldName,
      valueFieldName,
      eventFilterActionTypes,
      isStringCopyValueType,
      isStringTemplateValueType,
      isStringDictionaryValueType,
      isSelectValueType,
      changeActionType,
      updateCopyValue,
      remove,
    };
  },
};
</script>
