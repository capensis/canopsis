<template lang="pug">
  v-layout(column)
    v-flex(xs12)
      v-alert.mb-3(:value="!form.length", type="info") {{ $t('eventFilter.noExternalData') }}
    event-filter-enrichment-external-data-item-form.mb-3(
      v-for="(item, index) in form",
      v-field="form[index]",
      :name="item.key",
      :key="item.key",
      :disabled="disabled",
      @remove="removeItem"
    )
    v-flex(v-if="!disabled")
      v-btn.ml-0.my-0(
        color="primary",
        outline,
        @click="addItem"
      ) {{ $t('eventFilter.addExternalData') }}
</template>

<script>
import { eventFilterExternalDataItemToForm } from '@/helpers/forms/event-filter';

import { formArrayMixin } from '@/mixins/form';

import EventFilterEnrichmentExternalDataItemForm from './event-filter-enrichment-external-data-item-form.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterEnrichmentExternalDataItemForm },
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    addItem() {
      this.addItemIntoArray(eventFilterExternalDataItemToForm());
    },

    removeItem(key) {
      this.removeItemFromArrayWith(item => key !== item.key);
    },
  },
};
</script>
