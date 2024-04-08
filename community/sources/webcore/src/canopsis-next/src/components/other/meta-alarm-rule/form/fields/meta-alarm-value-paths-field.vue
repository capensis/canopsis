<template>
  <v-layout column>
    <v-layout
      v-for="(item, index) in items"
      :key="item[itemKey]"
      justify-space-between
      align-center
    >
      <v-text-field
        v-validate="validationRules"
        :value="item[itemValue]"
        :label="label"
        :name="getFieldName(item[itemKey])"
        :error-messages="errors.collect(getFieldName(item[itemKey]))"
        @input="updateFieldInArrayItem(index, itemValue, $event)"
      />
      <c-action-btn type="delete" @click="removeItemFromArray(index)" />
      <c-help-icon :text="$t('metaAlarmRule.valuePathHelpText')" icon="help" top />
    </v-layout>
    <c-btn-with-error :error="error ? $t('metaAlarmRule.errors.noValuePaths'): ''" @click="addNewItem" />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { defaultPrimitiveArrayItemCreator } from '@/helpers/entities/shared/form';

import { useArrayModelField } from '@/hooks/form/useArrayModelField';
import { useInjectValidator } from '@/hooks/form/useValidationChildren';

export default {
  inject: ['$validator'],
  model: {
    prop: 'items',
    event: 'input',
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    itemKey: {
      type: String,
      default: 'key',
    },
    name: {
      type: String,
      default: 'item',
    },
    required: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const validator = useInjectValidator();
    const { addItemIntoArray, updateFieldInArrayItem, removeItemFromArray } = useArrayModelField(props, emit);

    const validationRules = computed(() => ({
      required: props.required,
    }));

    const getNamePrefix = key => `${props.name}[${key}]`;
    const getFieldName = key => `${getNamePrefix(key)}.${props.itemValue}`;

    const addNewItem = () => {
      addItemIntoArray(defaultPrimitiveArrayItemCreator());
    };

    return {
      errors: validator.errors,
      validationRules,
      addNewItem,
      getFieldName,
      updateFieldInArrayItem,
      removeItemFromArray,
    };
  },
};
</script>
