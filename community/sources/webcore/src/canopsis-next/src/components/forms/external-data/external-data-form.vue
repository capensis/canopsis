<template lang="pug">
  v-layout(column)
    c-alert(
      :value="!form.length",
      type="info",
      transition="fade-transition"
    ) {{ $t('externalData.empty') }}
    external-data-item-form.mb-3(
      v-for="(item, index) in form",
      v-field="form[index]",
      :name="item.key",
      :key="item.key",
      :disabled="disabled",
      :types="types",
      :variables="variables",
      @remove="removeItemFromArray(index)"
    )
    v-flex(v-if="!disabled")
      v-btn.ml-0.my-0(
        color="primary",
        outline,
        @click="addItem"
      ) {{ $t('externalData.add') }}
</template>

<script>
import { externalDataItemToForm } from '@/helpers/forms/shared/external-data';

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
  },
  methods: {
    addItem() {
      this.addItemIntoArray(externalDataItemToForm());
    },
  },
};
</script>
