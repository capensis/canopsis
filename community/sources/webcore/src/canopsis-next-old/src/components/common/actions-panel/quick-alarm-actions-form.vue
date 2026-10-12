<template>
  <c-card-iterator-form
    v-field="actions"
    :handle="`.${dragItemHandleClass}`"
    :iterator-class="{ empty: isActionsEmpty }"
    @add="add"
  >
    <template #item="{ item: action, index }">
      <quick-alarm-actions-form-item
        v-field="actions[index].value"
        :key="action.key"
        :drag-handle-class="dragItemHandleClass"
        :name="action.key"
        :selected-actions="actions"
        :massive="massive"
        @remove="remove(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { computed } from 'vue';

import { widgetQuickAlarmActionToForm } from '@/helpers/entities/widget/quick-action/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useAsyncBootingParent } from '@/hooks/render/async-booting';

import QuickAlarmActionsFormItem from './quick-alarm-actions-form-item.vue';

export default {
  inject: ['$validator'],
  components: { QuickAlarmActionsFormItem },
  model: {
    prop: 'actions',
    event: 'input',
  },
  props: {
    actions: {
      type: Array,
      default: () => [],
    },
    massive: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const dragItemHandleClass = 'action-drag-handle';

    const isActionsEmpty = computed(() => !props.actions.length);

    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    useAsyncBootingParent(2);

    const add = () => addItemIntoArray(widgetQuickAlarmActionToForm());
    const remove = index => removeItemFromArray(index);

    return {
      dragItemHandleClass,
      isActionsEmpty,

      add,
      remove,
    };
  },
};
</script>
