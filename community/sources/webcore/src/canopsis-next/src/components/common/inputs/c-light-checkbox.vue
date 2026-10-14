<template>
  <label :class="classes" class="c-light-checkbox">
    <input
      ref="inputRef"
      :checked="modelValue"
      :disabled="disabled"
      class="c-light-checkbox__input"
      type="checkbox"
      @change="updateValue"
    />

    <span class="c-light-checkbox__control">
      <c-icon v-if="modelValue" class="c-light-checkbox__icon" size="small"> check </c-icon>
      <span v-else-if="indeterminate" class="c-light-checkbox__indeterminate" />
    </span>

    <span v-if="label" class="c-light-checkbox__label">
      {{ label }}
    </span>
  </label>
</template>

<script>
import { computed, ref, watch } from 'vue';

export default {
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    indeterminate: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const inputRef = ref();

    const classes = computed(() => ({
      'c-light-checkbox--checked': props.modelValue,
      'c-light-checkbox--indeterminate': props.indeterminate,
      'c-light-checkbox--disabled': props.disabled,
    }));

    const syncIndeterminate = () => {
      if (inputRef.value) {
        inputRef.value.indeterminate = props.indeterminate;
      }
    };

    const updateValue = event => emit('update:modelValue', event.target.checked);

    watch(() => props.indeterminate, syncIndeterminate, { immediate: true, flush: 'post' });

    return {
      inputRef,
      classes,
      updateValue,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-light-checkbox {
  align-items: center;
  cursor: pointer;
  display: inline-flex;
  gap: 4px;
  min-height: 24px;
  user-select: none;
  vertical-align: middle;

  &--disabled {
    cursor: default;
    opacity: 0.38;
  }

  &__input {
    height: 0;
    opacity: 0;
    pointer-events: none;
    position: absolute;
    width: 0;
  }

  &__control {
    align-items: center;
    border: 2px solid rgba(var(--v-theme-on-surface), 0.54);
    border-radius: 2px;
    box-sizing: border-box;
    display: inline-flex;
    height: 18px;
    justify-content: center;
    transition:
      background-color 0.15s ease,
      border-color 0.15s ease;
    width: 18px;
  }

  &--checked &__control,
  &--indeterminate &__control {
    background-color: rgb(var(--v-theme-primary));
    border-color: rgb(var(--v-theme-primary));
    color: rgb(var(--v-theme-on-primary));
  }

  &__icon {
    line-height: 1;
  }

  &__indeterminate {
    background-color: currentColor;
    display: block;
    height: 2px;
    width: 10px;
  }

  &__label {
    font-size: 0.875rem;
    line-height: 1;
  }
}
</style>
