<template>
  <c-card-iterator-form
    :value="modelValue"
    :handle="`.${dragItemHandleClass}`"
    :name="name"
    item-key="key"
    @update:value="$emit('update:modelValue', $event)"
    @add="add"
  >
    <template #item="{ item, index }">
      <c-card-iterator-item
        :key="item.key"
        :drag-handle-class="dragItemHandleClass"
        small
        static-delete
        @remove="remove(index)"
      >
        <template #header>
          <c-text-field
            :model-value="modelValue[index].value"
            :label="t('common.regexp')"
            :max-length="FORM_MAX_STRING_LENGTH"
            :name="`${name}.${item.key}.value`"
            required
            @update:model-value="updateFormat(index, $event)"
          />
        </template>
      </c-card-iterator-item>
    </template>
  </c-card-iterator-form>
</template>

<script>
import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { ticketFormatToForm } from '@/helpers/entities/widget/ticket-format/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'formats',
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const dragItemHandleClass = 'ticket-format-drag-handle';
    const { addItemIntoArray, updateFieldInArrayItem, removeItemFromArray } = useArrayModelField(
      props,
      emit,
    );

    const add = () => addItemIntoArray(ticketFormatToForm());
    const updateFormat = (index, value) => updateFieldInArrayItem(index, 'value', value);

    return {
      t,
      FORM_MAX_STRING_LENGTH,
      dragItemHandleClass,
      add,
      updateFormat,
      remove: removeItemFromArray,
    };
  },
};
</script>
