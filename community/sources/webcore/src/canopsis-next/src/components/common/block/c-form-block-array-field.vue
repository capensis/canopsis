<template>
  <c-layout gap="3" column>
    <!-- eslint-disable vue/multiline-html-element-content-newline -->
    <c-label
      v-if="label || $slots.label"
      :label="label"
      :required="required"
      :error="errorMessages.length > 0"
    ><slot name="label"/></c-label>
    <!-- eslint-enable vue/multiline-html-element-content-newline -->

    <div
      v-for="(item, index) in form"
      :key="getItemRowKey(item, index)"
      class="c-form-block-array-field__row"
    >
      <slot
        :add-item-into-array="addItemIntoArray"
        :form="form"
        :item="item"
        :index="index"
        :remove="getRemove(index)"
        name="item"
      />
    </div>

    <div v-if="!disabled">
      <slot :add="add" name="add-button">
        <v-btn :color="errorMessages.length ? 'error' : 'primary'" variant="outlined" @click="add">
          {{ resolvedAddButtonLabel }}
        </v-btn>
      </slot>
    </div>
    <v-messages v-if="errorMessages.length > 0" :messages="errorMessages" color="error" />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  emits: ['update:form'],
  props: {
    form: {
      type: Array,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    required: {
      type: Boolean,
      default: false,
    },
    itemToForm: {
      type: Function,
      required: true,
    },
    itemKey: {
      type: String,
      default: 'key',
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    addButtonLabel: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const resolvedAddButtonLabel = computed(() => props.addButtonLabel || t('common.add'));

    const { addItemIntoArray, removeItemFromArray, updateItemInArray, updateFieldInArrayItem } =
      useArrayModelField(props, emit);

    const add = () => addItemIntoArray(props.itemToForm());

    const getRemove = index => () => removeItemFromArray(index);

    const getItemRowKey = (item, index) => item?.[props.itemKey] ?? index;

    return {
      add,
      resolvedAddButtonLabel,
      addItemIntoArray,
      removeItemFromArray,
      updateItemInArray,
      updateFieldInArrayItem,
      getRemove,
      getItemRowKey,
    };
  },
};
</script>

<style scoped lang="scss">
.c-form-block-array-field__row {
  display: contents;
}
</style>
