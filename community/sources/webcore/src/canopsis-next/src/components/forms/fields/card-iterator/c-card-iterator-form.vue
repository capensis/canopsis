<template>
  <v-layout class="gap-2" column>
    <slot name="header" />
    <c-alert
      v-if="required && isEmpty && emptyMessage"
      type="info"
    >
      {{ emptyMessage }}
    </c-alert>
    <c-card-iterator-field
      v-field="value"
      v-bind="$attrs"
      :class="iteratorClass"
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
    <v-layout justify-end>
      <slot name="add-button">
        <c-btn-with-error
          v-if="required"
          :error="required && hasErrors ? (requiredErrorMessage || $t('common.addAtLeastOneItem')) : ''"
          :fab="!addButtonLabel"
          class="mt-3"
          small
          @click="$emit('add')"
        >
          <span v-if="addButtonLabel">{{ addButtonLabel }}</span>
          <v-icon v-else>
            add
          </v-icon>
        </c-btn-with-error>
        <v-tooltip v-else left>
          <template #activator="{ on }">
            <v-btn
              :error="required && hasErrors ? (requiredErrorMessage || $t('common.addAtLeastOneItem')) : ''"
              class="mt-3"
              color="primary"
              fab
              small
              v-on="on"
              @click="$emit('add')"
            >
              <v-icon>add</v-icon>
            </v-btn>
          </template>
          <span>{{ $t('common.add') }}</span>
        </v-tooltip>
      </slot>
    </v-layout>
  </v-layout>
</template>

<script>
import { computed, watch, onMounted, onBeforeUnmount } from 'vue';

import { useValidator } from '@/hooks/validator/validator';
import { useComponentInstance } from '@/hooks/vue';

export default {
  inject: ['$validator'],
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
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
    const validator = useValidator();
    const instance = useComponentInstance();

    const isEmpty = computed(() => !props.value?.length);
    const hasErrors = computed(() => validator?.errors?.has(props.name) ?? false);

    /**
     * Attaches min_value:1 validation rule to require at least one item when required prop is true
     */
    const attachMinValueRule = () => {
      if (validator?.attach && props.required) {
        validator.attach({
          name: props.name,
          rules: 'min_value:1',
          getter: () => props.value?.length ?? 0,
          vm: instance,
        });
      }
    };

    /**
     * Detaches the min_value validation rule from the validator
     */
    const detachMinValueRule = () => validator?.detach?.(props.name);

    watch(() => props.value, () => {
      if (validator?.validate && props.required) {
        instance.$nextTick(() => validator.validate(props.name));
      }
    });

    onMounted(() => props.required && attachMinValueRule());

    onBeforeUnmount(() => detachMinValueRule());

    return {
      isEmpty,
      hasErrors,
    };
  },
};
</script>
