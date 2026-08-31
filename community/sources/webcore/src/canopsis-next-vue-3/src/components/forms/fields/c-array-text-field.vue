<template>
  <c-layout column>
    <v-label v-if="label" :class="{ 'text-error': errorMessages.length > 0 }" class="mb-2">
      {{ label }}
    </v-label>
    <c-layout v-for="(item, index) in values" :key="item.key" align="center">
      <c-text-field
        :model-value="values[index].value"
        :validation-rules="rules"
        :disabled="disabled"
        :label="$t('common.value')"
        :name="`${name}.${item.key}.value`"
        class="not-required"
        @update:model-value="updateItemValue(index, $event)"
      />
      <c-action-btn v-if="!disabled" type="delete" @click="removeItemFromArray(index)" />
    </c-layout>
    <v-messages
      v-if="errorMessages.length > 0"
      :messages="errorMessages"
      color="error"
      class="mb-2"
    />
    <div>
      <v-btn
        :disabled="disabled"
        :color="errorMessages.length ? 'error' : 'primary'"
        variant="outlined"
        @click="addItem"
      >
        {{ $t('common.add') }}
      </v-btn>
    </div>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { uid } from '@/helpers/uid';

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
    const rules = computed(() => (props.maxLength ? { max: props.maxLength } : null));

    const updateValues = (values) => {
      emit('update:values', values);
      emit('change', values);
    };

    const updateItemValue = (index, value) => updateValues(
        props.values.map((item, itemIndex) => (itemIndex === index ? { ...item, value } : item)),
      );

    const removeItemFromArray = index => updateValues(props.values.filter((item, itemIndex) => itemIndex !== index));

    /**
     * Adds an empty object item with key and value to the array
     */
    const addItem = () => updateValues([...props.values, { key: uid(), value: '' }]);

    return {
      rules,

      addItem,
      updateItemValue,
      removeItemFromArray,
    };
  },
};
</script>
