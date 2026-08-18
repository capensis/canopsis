<template>
  <c-layout column>
    <c-draggable-list-field
      :value="actions"
      handle=".action-drag-handler"
      @update:value="updateModel"
    >
      <event-filter-enrichment-action-field
        v-for="(action, index) in actions"
        :key="action.key"
        :action="action"
        :name="`${name}.${action.key}`"
        :variables="variables"
        :copy-variables="copyVariables"
        :set-tags-items="setTagsItems"
        class="mb-3"
        @update:action="updateItemInArray(index, $event)"
        @remove="removeItemFromArray(index)"
      />
      <v-btn :color="hasError ? 'error' : 'primary'" outlined @click="addAction">
        {{ $t('eventFilter.addAction') }}
      </v-btn>
    </c-draggable-list-field>
  </c-layout>
</template>

<script>
import { computed, toRef } from 'vue';

import { eventFilterActionToForm } from '@/helpers/entities/event-filter/rule/form';

import { useFieldError } from '@/hooks/form';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import EventFilterEnrichmentActionField from './event-filter-enrichment-action-field.vue';

export default {
  emits: ['update:actions'],
  components: { EventFilterEnrichmentActionField },
  props: {
    actions: {
      type: Array,
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
      default: 'config.actions',
    },
    setTagsItems: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { updateModel, addItemIntoArray, updateItemInArray, removeItemFromArray } =
      useArrayModelField(props, emit);

    const errorMessage = useFieldError(toRef(props, 'name'));
    const hasError = computed(() => !!errorMessage.value);
    const addAction = () => addItemIntoArray(eventFilterActionToForm());

    return {
      hasError,
      updateModel,
      updateItemInArray,
      removeItemFromArray,
      addAction,
    };
  },
};
</script>
