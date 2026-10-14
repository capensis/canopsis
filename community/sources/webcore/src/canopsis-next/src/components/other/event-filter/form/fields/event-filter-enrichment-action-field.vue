<template>
  <v-card>
    <v-card-text>
      <c-layout align="start">
        <v-icon class="draggable ml-0 mr-3 mt-3 action-drag-handler">
          drag_indicator
        </v-icon>
        <c-layout class="w-100" column>
          <c-layout align="center" gap="2">
            <c-select
              :model-value="action.type"
              :items="eventFilterActionTypes"
              :label="$t('common.type')"
              item-title="text"
              item-value="value"
              @update:model-value="changeActionType"
            >
              <template #selection="{ item }">
                {{ item.groupTitle }} {{ item.text }}
              </template>
            </c-select>
            <v-btn
              class="mr-0"
              density="compact"
              size="small"
              variant="text"
              icon
              @click="remove"
            >
              <v-icon color="error">
                delete
              </v-icon>
            </v-btn>
          </c-layout>
          <v-expand-transition>
            <event-filter-enrichment-action-form-type-info v-if="action.type" :type="action.type" />
          </v-expand-transition>
          <c-text-field
            key="description"
            :model-value="action.description"
            :label="$t('common.description')"
            @update:model-value="updateField('description', $event)"
          />
          <c-layout v-if="isStringDictionaryValueType">
            <c-text-field
              key="value"
              :model-value="action.value"
              :label="$t('common.value')"
              :name="`${name}.value`"
              :prefix="eventExtraPrefix"
              class="event-filter-enrichment-action-field__dictionary-value"
              active
              clearable
              required
              @update:model-value="updateField('value', $event)"
            />
          </c-layout>
          <v-row v-else>
            <v-col cols="5">
              <c-combobox
                v-if="isSetFieldValueType"
                key="name-set-field"
                :model-value="action.name"
                :items="EVENT_FILTER_SET_FIELD_ITEMS"
                :label="$t('common.name')"
                :name="`${name}.name`"
                class="mr-2"
                required
                @update:model-value="updateField('name', $event)"
              />
              <c-name-field
                v-else
                key="name"
                :model-value="action.name"
                :name="`${name}.name`"
                class="mr-2"
                required
                @update:model-value="updateField('name', $event)"
              />
            </v-col>
            <v-col cols="7">
              <c-payload-text-field
                v-if="isStringTemplateValueType"
                key="from"
                :value="action.value"
                :label="$t('common.value')"
                :variables="variables"
                :name="`${name}.value`"
                class="ml-2"
                required
                clearable
                @update:value="updateField('value', $event)"
              />
              <c-combobox
                v-else-if="isStringCopyValueType"
                key="value"
                :model-value="action.value"
                :label="$t('common.value')"
                :items="copyVariables"
                :name="`${name}.value`"
                :menu-props="comboboxMenuProps"
                :return-object="false"
                children-key="variables"
                class="ml-2"
                required
                @update:model-value="updateCopyValue"
              />
              <event-filter-enrichment-action-form-select-tags-value
                v-else-if="isSelectValueType"
                key="value-tags"
                :value="action.value"
                :items="setTagsItems"
                :name="`${name}.value`"
                @update:value="updateField('value', $event)"
              />
              <c-mixed-field
                v-else
                key="value-mixed"
                :value="action.value"
                :label="$t('common.value')"
                :name="`${name}.value`"
                :types="mixedFieldTypes"
                class="ml-2"
                @update:value="updateField('value', $event)"
              />
            </v-col>
          </v-row>
        </c-layout>
      </c-layout>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed } from 'vue';

import {
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES,
  EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_GROUPS,
  EVENT_FILTER_EVENT_EXTRA_PREFIX,
  EVENT_FILTER_SET_FIELD_ITEMS,
  PATTERN_FIELD_TYPES,
} from '@/constants';

import {
  eventFilterDictionaryActionValueToForm,
  formToEventFilterDictionaryActionValue,
} from '@/helpers/entities/event-filter/rule/form';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

import EventFilterEnrichmentActionFormSelectTagsValue from './event-filter-enrichment-action-form-select-tags-value.vue';
import EventFilterEnrichmentActionFormTypeInfo from './event-filter-enrichment-action-form-type-info.vue';

export default {
  emits: ['update:action', 'remove'],
  components: {
    EventFilterEnrichmentActionFormTypeInfo,
    EventFilterEnrichmentActionFormSelectTagsValue,
  },
  props: {
    action: {
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
    const { updateField, updateModel } = useModelField(props, emit);

    const eventExtraPrefix = EVENT_FILTER_EVENT_EXTRA_PREFIX;
    const mixedFieldTypes = [
      { value: PATTERN_FIELD_TYPES.string },
      { value: PATTERN_FIELD_TYPES.number },
      { value: PATTERN_FIELD_TYPES.boolean },
      { value: PATTERN_FIELD_TYPES.stringArray },
    ];

    const comboboxMenuProps = { minWidth: 200 };

    const eventFilterActionTypes = computed(() =>
      Object.entries(EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES_GROUPS).reduce((acc, [group, groupItems]) => {
        const groupTitle = t(`eventFilter.actionsTypesGroups.${group}`);

        acc.push(
          { header: groupTitle },
          ...groupItems.map(value => ({
            value,
            text: t(`eventFilter.actionsTypes.${value}.text`),
            groupTitle,
          })),
        );

        return acc;
      }, []),
    );

    const isSetFieldValueType = computed(
      () => EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setField === props.action.type,
    );

    const isStringCopyValueType = computed(() => [
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copy,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.copyToEntityInfo,
      ].includes(props.action.type),
    );

    const isStringTemplateValueType = computed(() => [
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setFieldFromTemplate,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromTemplate,
        EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTagsFromTemplate,
      ].includes(props.action.type),
    );

    const isStringDictionaryValueType = computed(
      () => EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary === props.action.type,
    );

    const isSelectValueType = computed(
      () => EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTags === props.action.type,
    );

    /**
     * Changes the action type and handles value transformation for dictionary and select actions
     *
     * @param {string} type - The new action type from EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES
     */
    const changeActionType = (type) => {
      const newAction = {
        ...props.action,
        type,
      };

      if (type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setTags && props.action.type !== type) {
        newAction.value = '';
      } else if (props.action.type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary) {
        newAction.value = formToEventFilterDictionaryActionValue(props.action.value);
      } else if (type === EVENT_FILTER_ENRICHMENT_ACTIONS_TYPES.setEntityInfoFromDictionary) {
        newAction.value = eventFilterDictionaryActionValueToForm(props.action.value);
      }

      updateModel(newAction);
    };

    /**
     * Updates the copy value field by calling updateField with 'value' key
     *
     * @param {Object|string} [params={}] - Parameters object or string
     * @param {string} [params.value=''] - The value to set for the copy field
     */
    const updateCopyValue = newValue => updateField('value', newValue?.value || newValue);

    /**
     * Emits remove event to parent component to remove this action field
     */
    const remove = () => emit('remove', props.action);

    return {
      EVENT_FILTER_SET_FIELD_ITEMS,

      eventExtraPrefix,
      mixedFieldTypes,

      comboboxMenuProps,
      eventFilterActionTypes,
      isSetFieldValueType,
      isStringCopyValueType,
      isStringTemplateValueType,
      isStringDictionaryValueType,
      isSelectValueType,
      changeActionType,
      updateCopyValue,
      updateField,
      remove,
    };
  },
};
</script>

<style lang="scss" scoped>
.event-filter-enrichment-action-field__dictionary-value {
  :deep(.v-text-field__prefix) {
    padding-top: 20px;
  }
}
</style>
