<template>
  <c-layout gap="3" column>
    <slot name="header" />
    <c-alert v-if="required && isEmpty && emptyMessage" type="info">
      {{ emptyMessage }}
    </c-alert>
    <c-card-iterator-field
      :data="value"
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
    <c-layout justify="end">
      <slot name="add-button">
        <c-btn-with-error
          v-if="required"
          :error="
            required && hasErrors ? requiredErrorMessage || $t('common.addAtLeastOneItem') : ''
          "
          :fab="!addButtonLabel"
          :small="!addButtonLabel"
          variant="outlined"
          @click="$emit('add')"
        >
          <span v-if="addButtonLabel">
            {{ addButtonLabel }}
          </span>
          <v-icon v-else>
            add
          </v-icon>
        </c-btn-with-error>
        <v-tooltip v-else location="left">
          <template #activator="{ props: tooltipProps }">
            <v-btn
              :error="
                required && hasErrors ? requiredErrorMessage || $t('common.addAtLeastOneItem') : ''
              "
              color="primary"
              size="small"
              variant="outlined"
              icon
              v-bind="tooltipProps"
              @click="$emit('add')"
            >
              <v-icon>add</v-icon>
            </v-btn>
          </template>
          <span>{{ $t('common.add') }}</span>
        </v-tooltip>
      </slot>
    </c-layout>
  </c-layout>
</template>

<script>
import { computed, watch, onMounted, onBeforeUnmount } from 'vue';

import { useValidationAttachMinValue } from '@/hooks/validator/validation-attach-min-value';
import { useComponentInstance } from '@/hooks/vue';

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
    addButtonLabel: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const instance = useComponentInstance();
    const { validator, attachMinValueRule, detachMinValueRule, validateMinValueRule } =
      useValidationAttachMinValue(props.name);

    const isEmpty = computed(() => !props.value?.length);
    const hasErrors = computed(() => validator?.errors?.has(props.name) ?? false);

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
          instance.$nextTick(validateMinValueRule);
        }
      },
    );

    onMounted(() => props.required && attachRequiredMinValueRule());

    onBeforeUnmount(() => detachMinValueRule());

    return {
      isEmpty,
      hasErrors,
    };
  },
};
</script>
