<template>
  <v-layout column>
    <c-draggable-list-field
      v-field="actions"
      handle=".action-drag-handler"
    >
      <event-filter-enrichment-action-form
        v-for="(action, index) in actions"
        v-field="actions[index]"
        :key="action.key"
        :name="`${name}.${action.key}`"
        :variables="variables"
        :copy-variables="copyVariables"
        :set-tags-items="setTagsItems"
        class="mb-3"
        @remove="removeItemFromArray(index)"
      />
      <v-flex>
        <v-btn
          :color="errors.has(name) ? 'error' : 'primary'"
          outlined
          @click="addAction"
        >
          {{ $t('eventFilter.addAction') }}
        </v-btn>
      </v-flex>
    </c-draggable-list-field>
  </v-layout>
</template>

<script>
import { watch, onBeforeUnmount } from 'vue';

import { eventFilterActionToForm } from '@/helpers/entities/event-filter/rule/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';

import EventFilterEnrichmentActionForm from './event-filter-enrichment-action-form.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterEnrichmentActionForm },
  model: {
    prop: 'actions',
    event: 'input',
  },
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
      default: 'actions',
    },
    setTagsItems: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const addAction = () => addItemIntoArray(eventFilterActionToForm());

    const {
      attachRequiredRule,
      detachRequiredRule,
      validateRequiredRule,
    } = useValidationAttachRequired(props.name);

    watch(() => props.actions, validateRequiredRule);

    attachRequiredRule(() => props.actions);

    onBeforeUnmount(detachRequiredRule);

    return {
      addAction,
      removeItemFromArray,
    };
  },
};
</script>
