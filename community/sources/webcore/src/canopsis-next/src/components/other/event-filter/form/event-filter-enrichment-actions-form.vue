<template lang="pug">
  v-layout(column)
    draggable(v-field="actions", :options="draggableOptions")
      event-filter-enrichment-action-form.mb-3(
        v-for="(action, index) in actions",
        v-field="actions[index]",
        :key="action.key",
        @remove="removeAction"
      )
      v-flex
        v-btn.ml-0.my-0(color="primary", outline, @click="addAction") {{ $t('eventFilter.addAction') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { eventFilterActionToForm } from '@/helpers/forms/event-filter';

import { formArrayMixin } from '@/mixins/form';

import EventFilterEnrichmentActionForm from './event-filter-enrichment-action-form.vue';

export default {
  inject: ['$validator'],
  components: { Draggable, EventFilterEnrichmentActionForm },
  mixins: [formArrayMixin],
  model: {
    prop: 'actions',
    event: 'input',
  },
  props: {
    actions: {
      type: Array,
      required: true,
    },
  },
  computed: {
    draggableOptions() {
      return {
        handle: '.action-drag-handler',
      };
    },
  },
  methods: {
    addAction() {
      this.addItemIntoArray(eventFilterActionToForm());
    },

    removeAction(action) {
      this.removeItemFromArrayWith(({ key }) => key !== action.key);
    },
  },
};
</script>
