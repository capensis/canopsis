<template>
  <c-form-block-array-field
    v-field="form"
    :item-to-form="externalDataItemToForm"
    :label="$t('eventFilter.externalData')"
    :disabled="disabled"
    :add-button-label="$t('externalData.add')"
  >
    <template #item="{ item, index, remove }">
      <external-data-item-form
        v-field="form[index]"
        :name="`${name}.${item.key}`"
        :server-error-name="`${name}.${index}`"
        :disabled="disabled"
        :types="types"
        :variables="variables"
        :optionally="optionally"
        @remove="remove"
      />
    </template>
  </c-form-block-array-field>
</template>

<script>
import { externalDataItemToForm } from '@/helpers/entities/shared/external-data/form';

import ExternalDataItemForm from './external-data-item-form.vue';

export default {
  inject: ['$validator'],
  components: { ExternalDataItemForm },
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
    optionally: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    return {
      externalDataItemToForm,
    };
  },
};
</script>
