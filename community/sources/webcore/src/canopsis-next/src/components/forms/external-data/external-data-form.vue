<template>
  <v-layout class="gap-3" column>
    <c-label>{{ $t('externalData.title') }}</c-label>
    <external-data-item-form
      v-for="(item, index) in form"
      v-field="form[index]"
      :key="item.key"
      :name="`${name}.${item.key}`"
      :server-error-name="`${name}.${index}`"
      :disabled="disabled"
      :types="types"
      :variables="variables"
      :optionally="optionally"
      @remove="removeItemFromArray(index)"
    />
    <v-flex v-if="!disabled">
      <v-btn
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

import { useArrayModelField } from '@/hooks/form/array-model-field';

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
  setup(props, { emit }) {
    const { addItemIntoArray } = useArrayModelField(props, emit);

    const addItem = () => addItemIntoArray(externalDataItemToForm());

    return {
      addItem,
    };
  },
};
</script>
