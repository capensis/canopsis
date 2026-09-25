<template>
  <c-card-iterator-form
    :value="actions"
    :handle="`.${dragItemHandleClass}`"
    :name="name"
    :iterator-class="{ empty: isActionsEmpty }"
    @update:value="$emit('update:actions', $event)"
    @add="add"
  >
    <template #item="{ item: action, index }">
      <quick-alarm-actions-form-item
        :key="action.key"
        :model-value="actions[index].value"
        :drag-handle-class="dragItemHandleClass"
        :name="`${name}.${action.key}`"
        :selected-actions="actions"
        :massive="massive"
        @update:model-value="updateAction(index, $event)"
        @remove="remove(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { computed } from 'vue';

import { widgetQuickAlarmActionToForm } from '@/helpers/entities/widget/quick-action/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import QuickAlarmActionsFormItem from './quick-alarm-actions-form-item.vue';

export default {
  components: { QuickAlarmActionsFormItem },

  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    massive: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'actions',
    },
  },

  emits: ['update:actions'],

  setup(props, { emit }) {
    const dragItemHandleClass = 'action-drag-handle';

    const isActionsEmpty = computed(() => !props.actions.length);
    const { addItemIntoArray, updateFieldInArrayItem, removeItemFromArray } = useArrayModelField(
      props,
      emit,
    );

    const add = () => addItemIntoArray(widgetQuickAlarmActionToForm());
    const updateAction = (index, value) => updateFieldInArrayItem(index, 'value', value);
    const remove = index => removeItemFromArray(index);

    return {
      dragItemHandleClass,
      isActionsEmpty,
      add,
      updateAction,
      remove,
    };
  },
};
</script>
