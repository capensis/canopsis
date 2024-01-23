<template>
  <v-layout column>
    <c-draggable-list-field
      v-field="actions"
      handle=".action-drag-handler"
    >
      <event-filter-enrichment-action-form
        v-field="actions[index]"
        v-for="(action, index) in actions"
        :key="action.key"
        :name="`${name}.${action.key}`"
        :variables="variables"
        class="mb-3"
        @remove="removeItemFromArray(index)"
      />
      <v-flex>
        <v-btn
          class="ml-0 my-0"
          color="primary"
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
import { eventFilterActionToForm } from '@/helpers/entities/event-filter/rule/form';

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
    variables: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'actions',
    },
  },
  methods: {
    addAction() {
      this.addItemIntoArray(eventFilterActionToForm());
    },
  },
};
</script>
