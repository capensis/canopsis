<template>
  <v-layout column>
    <v-label v-if="label" :class="{ 'error--text': errorMessages.length > 0 }" class="mb-2 label-required">
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
    <c-alert
      :value="errorMessages.length > 0"
      class="mb-2"
      type="error"
    >
      {{ errorMessages.join('\n') }}
    </c-alert>
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
import { nextTick, onBeforeUnmount, watch } from 'vue';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useComponentInstance } from '@/hooks/vue';
import { useValidator } from '@/hooks/validator/validator';

export default {
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
  },
  setup(props, { emit }) {
    const validator = useValidator();
    const instance = useComponentInstance();
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    let valuesWatcher = null;

    /**
     * Adds an empty string item to the array
     */
    const addItem = () => addItemIntoArray('');

    /**
     * Attaches the array_required validation rule to the validator
     * and sets up a watcher to re-validate when values change
     */
    const attachArrayRequiredRule = () => {
      validator?.attach?.({
        name: props.name,
        rules: 'array_required',
        getter: () => props.values.length,
        vm: instance,
      });

      valuesWatcher = watch(() => props.values, () => nextTick(() => validator?.validate?.(props.name)));
    };

    /**
     * Detaches the validation rule from the validator
     * and cleans up the values watcher
     */
    const detachArrayRequiredRule = () => {
      if (valuesWatcher) {
        valuesWatcher();
        valuesWatcher = null;
      }
      validator?.detach?.(props.name);
    };

    watch(
      () => props.required,
      value => (value ? attachArrayRequiredRule() : detachArrayRequiredRule()),
      { immediate: true },
    );

    onBeforeUnmount(detachArrayRequiredRule);

    return {
      addItem,
      removeItemFromArray,
    };
  },
};
</script>
