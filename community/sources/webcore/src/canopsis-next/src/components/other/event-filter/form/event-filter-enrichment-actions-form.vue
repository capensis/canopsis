<template lang="pug">
  v-layout(column)
    event-filter-enrichment-action-form.mb-3(
      v-for="(action, index) in form",
      v-field="form[index]",
      :key="action.key",
      @remove="removeAction"
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
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Array,
      required: true,
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
