<template>
  <c-card-iterator-form
    :value="modelValue"
    :handle="`.${dragItemHandleClass}`"
    :empty-message="t('clock.validation.atLeastOneClock')"
    :required-error-message="t('clock.validation.atLeastOneClock')"
    name="parameters.clocks"
    item-key="key"
    add-button-justify="end"
    required
    @add="addClock"
    @update:value="$emit('update:modelValue', $event)"
  >
    <template #item="{ item: clock, index }">
      <clock-item-field
        :key="clock.key"
        :model-value="modelValue[index]"
        :name="`parameters.clocks.${clock.key}`"
        :drag-handle-class="dragItemHandleClass"
        :removable="modelValue.length > 1"
        @remove="removeItemFromArray(index)"
        @update:model-value="updateItemInArray(index, $event)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { сlockWidgetClockToForm } from '@/helpers/entities/widget/forms/clock';
import { uid } from '@/helpers/uid';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';

import ClockItemField from './clock-item-field.vue';

export default {
  components: { ClockItemField },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const dragItemHandleClass = 'clock-drag-handle';
    const { addItemIntoArray, removeItemFromArray, updateItemInArray } = useArrayModelField(
      props,
      emit,
    );
    const addClock = () => addItemIntoArray({
      key: uid(),

      ...сlockWidgetClockToForm(),
    });

    return {
      t,
      dragItemHandleClass,
      addClock,
      removeItemFromArray,
      updateItemInArray,
    };
  },
};
</script>
