<template>
  <v-layout column>
    <v-label v-if="label" :class="{ 'error--text': errorMessages.length > 0 }" class="mb-2">
      {{ label }}
    </v-label>
    <v-layout
      v-for="(item, index) in values"
      :key="item.key"
      align-center
    >
      <v-text-field
        v-field="values[index].value"
        v-validate="rules"
        :disabled="disabled"
        :label="$t('common.value')"
        :name="`${name}.${item.key}.value`"
        :error-messages="errors.collect(`${name}.${item.key}.value`)"
        class="not-required"
      />
      <c-action-btn
        v-if="!disabled"
        type="delete"
        @click="removeItemFromArray(index)"
      />
    </v-layout>
    <v-messages
      v-if="errorMessages.length > 0"
      :value="errorMessages"
      color="error"
      class="mb-2"
    />
    <v-flex>
      <v-btn
        :disabled="disabled"
        :color="errorMessages.length ? 'error' : 'primary'"
        outlined
        @click="addItem"
      >
        {{ $t('common.add') }}
      </v-btn>
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { uid } from '@/helpers/uid';

import { useArrayModelField } from '@/hooks/form/array-model-field';

export default {
  $_veeValidate: {
    name() {
      return this.name;
    },
    value() {
      return this.values;
    },
  },
  inject: ['$validator'],
  model: {
    prop: 'values',
    event: 'change',
  },
  props: {
    name: {
      type: String,
      default: '',
    },
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
    required: {
      type: Boolean,
      default: false,
    },
    maxLength: {
      type: Number,
      default: null,
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const rules = computed(() => (props.maxLength ? { max: props.maxLength } : null));
    /**
     * Adds an empty object item with key and value to the array
     */
    const addItem = () => addItemIntoArray({ key: uid(), value: '' });

    return {
      rules,

      addItem,
      removeItemFromArray,
    };
  },
};
</script>
