<template>
  <v-layout column>
    <c-card-iterator-field
      :data="actions"
      :handle="`.${dragItemHandleClass}`"
      :class="{ empty: isActionsEmpty }"
      @input="$emit('input', $event)"
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
    </c-card-iterator-field>
    <v-layout justify-end>
      <v-tooltip left>
        <template #activator="{ on }">
          <v-btn
            class="mt-3"
            color="primary"
            fab
            small
            v-on="on"
            @click.prevent="add"
          >
            <v-icon>add</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('common.add') }}</span>
      </v-tooltip>
    </v-layout>
  </v-layout>
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
