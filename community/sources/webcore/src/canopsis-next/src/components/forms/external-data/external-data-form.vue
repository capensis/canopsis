<template>
  <v-layout column>
    <c-alert
      :value="!form.length"
      type="info"
      transition="fade-transition"
    >
      {{ $t('externalData.empty') }}
    </c-alert>
    <external-data-item-form
      v-for="(item, index) in form"
      v-field="form[index]"
      :key="item.key"
      :name="`${name}.${item.key}`"
      :disabled="disabled"
      :types="types"
      :variables="variables"
      class="mb-3"
      @remove="removeItemFromArray(index)"
    />
    <v-flex v-if="!disabled">
      <v-btn
        class="ml-0 my-0"
        color="primary"
        outlined
        @click="addItem"
      >
        {{ $t('externalData.add') }}
      </v-btn>
    </v-flex>
  </v-layout>
</template>

<script>
import { externalDataItemToForm } from '@/helpers/entities/shared/external-data/form';

import { formArrayMixin } from '@/mixins/form';

import ExternalDataItemForm from './external-data-item-form.vue';

export default {
  inject: ['$validator'],
  components: { ExternalDataItemForm },
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
    types: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => ([]),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'external_data',
    },
  },
  methods: {
    addItem() {
      this.addItemIntoArray(externalDataItemToForm());
    },
  },
};
</script>
