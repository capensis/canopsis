<template>
  <widget-settings-item :title="t('settings.fastPbehavior')">
    <c-card-iterator-form
      :value="modelValue"
      :handle="`.${dragHandleClass}`"
      @update:value="$emit('update:modelValue', $event)"
      @add="add"
    >
      <template #header>
        <span class="text-body-2">
          {{ t('settings.fastPbehavior') }}
        </span>
        <span>{{ t('settings.fastPbehaviorDescription') }}</span>
      </template>
      <template #item="{ index }">
        <c-card-iterator-item
          :drag-handle-class="dragHandleClass"
          small
          @remove="remove(index)"
        >
          <template #header>
            <c-layout
              gap="2"
              column
              fill
            >
              <c-name-field
                :model-value="modelValue[index].name_prefix"
                :label="t('common.namePrefix')"
                :max-length="FORM_MAX_STRING_LENGTH"
                :name="getName(modelValue[index], 'name_prefix')"
                required
                @update:model-value="updateField(index, 'name_prefix', $event)"
              />
              <c-pbehavior-type-field
                :model-value="modelValue[index].type"
                :label="t('common.type')"
                :types="types"
                :name="getName(modelValue[index], 'type')"
                independent
                clearable
                required
                @update:model-value="updateField(index, 'type', $event)"
              />
              <c-pbehavior-reason-field
                :model-value="modelValue[index].reason"
                :label="t('common.reason')"
                :name="getName(modelValue[index], 'reason')"
                clearable
                required
                @update:model-value="updateField(index, 'reason', $event)"
              />
            </c-layout>
          </template>
        </c-card-iterator-item>
      </template>
    </c-card-iterator-form>
  </widget-settings-item>
</template>

<script>
import { FORM_MAX_STRING_LENGTH, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { alarmListFastPbehaviorParametersItemToForm } from '@/helpers/entities/widget/forms/alarm';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const types = [PBEHAVIOR_TYPE_TYPES.pause];
    const dragHandleClass = 'fast-pbehavior-drag-handle';

    const { addItemIntoArray, updateFieldInArrayItem, removeItemFromArray } = useArrayModelField(
      props,
      emit,
    );

    const add = () => addItemIntoArray(alarmListFastPbehaviorParametersItemToForm());
    const updateField = (index, field, value) => updateFieldInArrayItem(index, field, value);
    const getName = (item, suffix) => `${item.key}${suffix}`;

    return {
      t,
      FORM_MAX_STRING_LENGTH,
      types,
      dragHandleClass,
      add,
      updateField,
      remove: removeItemFromArray,
      getName,
    };
  },
};
</script>
