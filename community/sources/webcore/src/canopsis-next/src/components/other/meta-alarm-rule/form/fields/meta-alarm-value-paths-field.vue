<template>
  <c-form-block-row :label="t('metaAlarmRule.valuePath', 2)" indented>
    <c-layout gap="3" column>
      <c-label
        :label="t('metaAlarmRule.valuePath', 2)"
        :help-text="t('metaAlarmRule.valuePathHelpText')"
        :required="required"
      />

      <c-layout
        v-for="(item, index) in items"
        :key="item[itemKey]"
        gap="2"
        align="center"
      >
        <c-name-field
          :model-value="item[itemValue]"
          :label="label"
          :name="getFieldName(item[itemKey])"
          :required="required"
          @update:model-value="updateFieldInArrayItem(index, itemValue, $event)"
        />
        <c-action-btn type="delete" @click="removeItemFromArray(index)" />
      </c-layout>

      <c-btn-with-error
        :error="addButtonError"
        variant="outlined"
        @click="addNewItem"
      >
        {{ t('metaAlarmRule.addValuePath') }}
      </c-btn-with-error>
    </c-layout>
  </c-form-block-row>
</template>

<script>
import { computed, watch } from 'vue';

import { defaultPrimitiveArrayItemCreator } from '@/helpers/entities/shared/form';

import { useArrayModelField, useField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
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
      default: 'value_paths',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:items'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { addItemIntoArray, updateFieldInArrayItem, removeItemFromArray } =
      useArrayModelField(props, emit);

    const fieldRules = () => {
      if (!props.required) {
        return true;
      }

      return props.items.length > 0 || t('metaAlarmRule.errors.noValuePaths');
    };

    const { errorMessage, validate } = useField(props.name, fieldRules, {
      initialValue: props.items.length,
      validateOnValueUpdate: false,
    });

    const addButtonError = computed(() => errorMessage.value);
    const getFieldName = key => `${props.name}.${key}.${props.itemValue}`;
    const addNewItem = () => addItemIntoArray(defaultPrimitiveArrayItemCreator());

    watch(
      () => props.items.length,
      () => validate(),
      { flush: 'post' },
    );

    return {
      t,
      addButtonError,
      addNewItem,
      getFieldName,
      updateFieldInArrayItem,
      removeItemFromArray,
    };
  },
};
</script>
