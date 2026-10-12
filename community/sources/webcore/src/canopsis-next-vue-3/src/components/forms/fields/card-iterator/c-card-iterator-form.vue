<template>
  <c-layout gap="3" column>
    <slot name="header" />
    <c-alert v-if="required && isEmpty && emptyMessage" type="info">
      {{ emptyMessage }}
    </c-alert>
    <c-card-iterator-field
      :data="value"
      :disabled="disabled"
      v-bind="$attrs"
      :class="iteratorClass"
      @update:data="$emit('update:value', $event)"
    >
      <template #item="{ item, index, handle: itemHandle }">
        <slot
          :item="item"
          :index="index"
          :handle="itemHandle"
          name="item"
        />
      </template>
    </c-card-iterator-field>
    <c-layout :justify="addButtonJustify">
      <c-tooltip v-if="!addButtonLabel && !addButtonError" location="left">
        <template #activator="{ props: tooltipProps }">
          <c-btn-with-error
            :error="addButtonError"
            :disabled="disabled"
            :variant="addButtonVariant"
            size="small"
            icon
            v-bind="tooltipProps"
            @click="$emit('add')"
          >
            <v-icon>add</v-icon>
          </c-btn-with-error>
        </template>
        <span>{{ $t('common.add') }}</span>
      </c-tooltip>

      <c-btn-with-error
        v-else
        :error="addButtonError"
        :disabled="disabled"
        :icon="!addButtonLabel"
        :size="!addButtonLabel ? 'small' : undefined"
        :variant="addButtonVariant"
        @click="$emit('add')"
      >
        <span v-if="addButtonLabel">
          {{ addButtonLabel }}
        </span>
        <v-icon v-else>
          add
        </v-icon>
      </c-btn-with-error>
    </c-layout>
  </c-layout>
</template>

<script>
import { computed, nextTick, onBeforeUnmount, onMounted, toRef, watch } from 'vue';

import { useFieldError, useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useValidationAttachMinValue } from '@/hooks/validator/validation-attach-min-value';

export default {
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  emits: ['update:value', 'add'],
  props: {
    value: {
      type: Array,
      default: () => [],
    },
    iteratorClass: {
      type: [String, Object],
      default: null,
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'items',
    },
    emptyMessage: {
      type: String,
      default: '',
    },
    requiredErrorMessage: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    addButtonLabel: {
      type: String,
      required: false,
    },
    addButtonJustify: {
      type: String,
      default: '',
    },
    addButtonVariant: {
      type: String,
      default: 'outlined',
    },
  },
  setup(props) {
    const { t } = useI18n();
    const fieldName = toRef(props, 'name');
    useFormContext();
    const fieldError = useFieldError(fieldName);
    const { hasError, attachMinValueRule, detachMinValueRule, validateMinValueRule } =
      useValidationAttachMinValue(fieldName);

    const isEmpty = computed(() => !props.value?.length);
    const hasErrors = computed(() => hasError(fieldName.value));
    const requiredAddButtonError = computed(() => props.required && hasErrors.value
        ? props.requiredErrorMessage || fieldError.value || t('common.addAtLeastOneItem')
        : '',
    );
    const addButtonError = computed(() => requiredAddButtonError.value || fieldError.value);

    /**
     * Attaches min_value:1 validation rule to require at least one item when required prop is true
     */
    const attachRequiredMinValueRule = () => {
      if (props.required) {
        attachMinValueRule(() => props.value?.length ?? 0);
      }
    };

    watch(
      () => props.value,
      () => {
        if (props.required) {
          nextTick(validateMinValueRule);
        }
      },
    );

    onMounted(() => props.required && attachRequiredMinValueRule());

    onBeforeUnmount(() => detachMinValueRule());

    return {
      isEmpty,
      hasErrors,
      addButtonError,
    };
  },
};
</script>
