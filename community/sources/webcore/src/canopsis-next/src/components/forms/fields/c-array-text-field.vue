<template>
  <v-layout column>
    <v-label v-if="label" class="mb-2">
      {{ label }}
    </v-label>
    <v-layout
      v-for="(value, index) in values"
      :key="index"
      align-center
    >
      <v-text-field
        v-field="values[index]"
        :disabled="disabled"
        :label="$t('common.value')"
      />
      <c-action-btn
        v-if="!disabled"
        type="delete"
        @click="removeItemFromArray(index)"
      />
    </v-layout>
    <v-messages
      :value="errorMessages"
      color="error"
    />
    <v-flex>
      <v-btn
        :disabled="disabled"
        color="primary"
        outlined
        @click="addItem"
      >
        {{ $t('common.add') }}
      </v-btn>
    </v-flex>
  </v-layout>
</template>

<script>
import { useArrayModelField } from '@/hooks/form/array-model-field';

export default {
  inject: ['$validator'],
  model: {
    prop: 'values',
    event: 'change',
  },
  props: {
    values: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const addItem = () => addItemIntoArray('');

    return {
      addItem,
      removeItemFromArray,
    };
  },
};
</script>
