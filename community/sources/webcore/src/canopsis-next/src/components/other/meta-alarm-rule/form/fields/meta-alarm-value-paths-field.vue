<template>
  <v-layout column>
    <v-layout
      v-for="(item, index) in items"
      :key="item[itemKey]"
      justify-space-between
      align-center
    >
      <c-name-field
        :value="item[itemValue]"
        :label="label"
        :name="getFieldName(item[itemKey])"
        :required="required"
        @input="updateFieldInArrayItem(index, itemValue, $event)"
      />
      <c-action-btn type="delete" @click="removeItemFromArray(index)" />
      <c-help-icon :text="$t('metaAlarmRule.valuePathHelpText')" icon="help" top />
    </v-layout>
    <c-btn-with-error :error="hasValuePathsErrors ? $t('metaAlarmRule.errors.noValuePaths'): ''" @click="addNewItem" />
  </v-layout>
</template>

<script>
import { computed, nextTick, onBeforeUnmount, watch } from 'vue';

import { defaultPrimitiveArrayItemCreator } from '@/helpers/entities/shared/form';

import { useArrayModelField } from '@/hooks/form/useArrayModelField';
import { useInjectValidator } from '@/hooks/validator/inject-validator';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';

export default {
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
  },
  setup(props, { emit }) {
    const validator = useInjectValidator();
    const { addItemIntoArray, updateFieldInArrayItem, removeItemFromArray } = useArrayModelField(props, emit);
    const {
      attachRequiredRule,
      detachRequiredRule,
      validateRequiredRule,
    } = useValidationAttachRequired(props.name);

    const hasValuePathsErrors = computed(() => validator.errors.has(props.name));

    const getFieldName = key => `${props.name}[${key}].${props.itemValue}`;

    const addNewItem = () => {
      addItemIntoArray(defaultPrimitiveArrayItemCreator());
    };

    const isValuePathsExist = () => props.items && props.items.length > 0;

    watch(
      () => props.required,
      () => {
        if (props.required) {
          attachRequiredRule(isValuePathsExist);
        } else {
          detachRequiredRule();
        }
      },
      { immediate: true },
    );

    watch(() => props.items, () => {
      nextTick(validateRequiredRule);
    });

    onBeforeUnmount(detachRequiredRule);

    return {
      hasValuePathsErrors,
      addNewItem,
      getFieldName,
      updateFieldInArrayItem,
      removeItemFromArray,
    };
  },
};
</script>
