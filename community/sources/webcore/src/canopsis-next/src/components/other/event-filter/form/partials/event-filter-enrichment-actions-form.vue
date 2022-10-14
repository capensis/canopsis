<template lang="pug">
  v-layout(column)
    c-draggable-list-field(v-field="actions", handle=".action-drag-handler")
      event-filter-enrichment-action-form.mb-3(
        v-for="(action, index) in actions",
        v-field="actions[index]",
        :key="action.key",
        @remove="removeItemFromArray(index)"
      )
      v-flex
        v-btn.ml-0.my-0(color="primary", outline, @click="addAction") {{ $t('eventFilter.addAction') }}
</template>

<script>
import { eventFilterActionToForm } from '@/helpers/forms/event-filter';

import { formArrayMixin } from '@/mixins/form';

import EventFilterEnrichmentActionForm from './event-filter-enrichment-action-form.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterEnrichmentActionForm },
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
  methods: {
    addAction() {
      this.addItemIntoArray(eventFilterActionToForm());
    },
  },
};
</script>
