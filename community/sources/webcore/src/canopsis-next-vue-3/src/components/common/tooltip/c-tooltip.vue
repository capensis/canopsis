<template>
  <v-tooltip
    v-model="opened"
    v-bind="$attrs"
    :disabled="tooltipDisabled"
    :eager="false"
    :location="tooltipLocation"
  >
    <template v-if="$slots.activator" #activator="activatorSlotProps">
      <slot
        v-bind="activatorSlotProps"
        :on="activatorSlotProps.props"
        name="activator"
      />
    </template>

    <slot v-if="$slots.default" />
    <template v-else-if="html">
      <span v-if="opened" v-html="sanitizedContent" />
    </template>
    <span v-else>{{ content }}</span>
  </v-tooltip>
</template>

<script>
import { computed, ref } from 'vue';

import { sanitizeHtml } from '@/helpers/html';

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
      type: Boolean,
      default: undefined,
    },
    content: {
      type: String,
      default: undefined,
    },
    location: {
      type: String,
      default: undefined,
    },
    top: {
      type: Boolean,
      default: false,
    },
    right: {
      type: Boolean,
      default: false,
    },
    bottom: {
      type: Boolean,
      default: false,
    },
    left: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    html: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const internallyOpened = ref(false);

    const opened = computed({
      get: () => props.modelValue ?? internallyOpened.value,
      set: (value) => {
        internallyOpened.value = value;
        emit('update:modelValue', value);
      },
    });

    const tooltipDisabled = computed(() => (
      props.disabled || (props.content !== undefined && !props.content)
    ));

    const tooltipLocation = computed(() => {
      if (props.location) {
        return props.location;
      }

      if (props.top) {
        return 'top';
      }

      if (props.right) {
        return 'right';
      }

      if (props.left) {
        return 'left';
      }

      if (props.bottom || props.content !== undefined) {
        return 'bottom';
      }

      return undefined;
    });

    const sanitizedContent = computed(() => (
      opened.value && props.html ? sanitizeHtml(props.content) : ''
    ));

    return {
      opened,
      tooltipDisabled,
      tooltipLocation,
      sanitizedContent,
    };
  },
};
</script>
