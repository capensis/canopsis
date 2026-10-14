<template>
  <c-form-block-array-field
    :form="externalData"
    :item-to-form="externalDataItemToForm"
    :label="$t('eventFilter.externalData')"
    :disabled="disabled"
    :add-button-label="$t('externalData.add')"
    @update:form="updateModel"
  >
    <template #item="{ item, index, remove }">
      <external-data-item-form
        :form="externalData[index]"
        :name="`${name}.${item.key}`"
        :server-error-name="`${name}.${index}`"
        :disabled="disabled"
        :types="types"
        :variables="variables"
        :optionally="optionally"
        @update:form="updateItemInArray(index, $event)"
        @remove="remove"
      />
    </template>
  </c-form-block-array-field>
</template>

<script>
import { externalDataItemToForm } from '@/helpers/entities/shared/external-data/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import ExternalDataItemForm from './external-data-item-form.vue';

export default {
  emits: ['update:external-data'],
  components: { ExternalDataItemForm },
  props: {
    externalData: {
      type: Array,
      required: true,
    },
    types: {
      type: Array,
      default: () => [],
    },
    variables: {
      type: Array,
      default: () => [],
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
  setup(props, { emit }) {
    const { updateModel, updateItemInArray } = useArrayModelField(props, emit);

    return {
      externalDataItemToForm,
      updateModel,
      updateItemInArray,
    };
  },
};
</script>
