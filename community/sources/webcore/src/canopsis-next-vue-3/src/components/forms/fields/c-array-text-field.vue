<template>
  <c-layout column>
    <v-label v-if="label" :class="{ 'text-error': preparedErrorMessages.length > 0 }" class="mb-2">
      <span v-if="required" class="text-error"> * </span>
      {{ label }}
    </v-label>
    <c-layout v-for="(item, index) in values" :key="item.key" align="center">
      <c-text-field
        :model-value="values[index].value"
        :required="required"
        :max-length="maxLength"
        :disabled="disabled"
        :label="$t('common.value')"
        :name="`${name}.${item.key}.value`"
        class="not-required"
        @update:model-value="updateItemValue(index, $event)"
      />
      <c-action-btn v-if="!disabled" type="delete" @click="removeItemFromArray(index)" />
    </c-layout>
    <div>
      <c-btn-with-error
        :disabled="disabled"
        :error="errorMessage"
        variant="outlined"
        @click="addItem"
      >
        {{ $t('common.add') }}
      </c-btn-with-error>
    </div>
  </c-layout>
</template>

<script>
import { computed, toRef, watch } from 'vue';

import { uid } from '@/helpers/uid';

import { useField } from '@/hooks/form';

export default {
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

  emits: ['change', 'update:values'],

  setup(props, { emit }) {
    const fieldName = computed(() => props.name || props.label || 'values');
    const rules = computed(() => (props.required ? { required: true } : null));
    const { errorMessage, setValue, validate } = useField(fieldName, rules, {
      initialValue: props.values,
      label: toRef(() => props.label || props.name),
      validateOnValueUpdate: false,
    });

    const preparedErrorMessages = computed(() => {
      if (props.errorMessages.length) {
        return props.errorMessages;
      }

      return errorMessage.value ? [errorMessage.value] : [];
    });
    const errorMessageText = computed(() => preparedErrorMessages.value.join(', '));

    const updateValues = (values) => {
      emit('update:values', values);
      emit('change', values);
      setValue(values, false);

      if (errorMessage.value) {
        validate();
      }
    };

    const updateItemValue = (index, value) => updateValues(
        props.values.map((item, itemIndex) => (itemIndex === index ? { ...item, value } : item)),
      );

    const removeItemFromArray = index => updateValues(props.values.filter((item, itemIndex) => itemIndex !== index));

    /**
     * Adds an empty object item with key and value to the array
     */
    const addItem = () => updateValues([...props.values, { key: uid(), value: '' }]);

    watch(
      () => props.values,
      value => setValue(value, false),
    );

    return {
      preparedErrorMessages,
      errorMessage: errorMessageText,

      addItem,
      updateItemValue,
      removeItemFromArray,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-label.text-error {
  opacity: 1;
}
</style>
