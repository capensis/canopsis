<template>
  <v-layout class="gap-3" column>
    <!-- eslint-disable vue/multiline-html-element-content-newline -->
    <c-label
      v-if="label || $slots.label"
      :required="required"
      :error="displayErrorMessages.length > 0"
    ><slot name="label">{{ label }}</slot></c-label>
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
        :remove-item-from-array="removeItemFromArray"
        :update-field-in-array-item="updateFieldInArrayItem"
        :update-item-in-array="updateItemInArray"
        name="item"
      />
    </div>

    <div v-if="!disabled">
      <slot
        :add="add"
        name="add-button"
      >
        <v-btn
          :color="displayErrorMessages.length ? 'error' : 'primary'"
          outlined
          @click="add"
        >
          {{ resolvedAddButtonLabel }}
        </v-btn>
      </slot>
    </div>
    <v-messages
      v-if="displayErrorMessages.length > 0"
      :value="displayErrorMessages"
      color="error"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useArrayModelField } from '@/hooks/form/array-model-field';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
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
    requiredErrorMessage: {
      type: String,
      default: '',
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

    const resolvedAddButtonLabel = computed(
      () => props.addButtonLabel || t('common.add'),
    );

    const displayErrorMessages = computed(() => {
      if (props.requiredErrorMessage && props.errorMessages.length > 0) {
        return [props.requiredErrorMessage];
      }

      return props.errorMessages;
    });

    const {
      addItemIntoArray,
      removeItemFromArray,
      updateItemInArray,
      updateFieldInArrayItem,
    } = useArrayModelField(props, emit);

    const add = () => addItemIntoArray(props.itemToForm());

    const getRemove = index => () => removeItemFromArray(index);

    /**
     * Stable key for list rows when items expose `itemKey`, otherwise fallback to index.
     *
     * @param {*} item
     * @param {number} index
     * @returns {string|number}
     */
    const getItemRowKey = (item, index) => {
      const keyField = props.itemKey;

      if (item && typeof item === 'object' && keyField in item && item[keyField] !== undefined) {
        return item[keyField];
      }

      return index;
    };

    return {
      add,
      resolvedAddButtonLabel,
      displayErrorMessages,
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
