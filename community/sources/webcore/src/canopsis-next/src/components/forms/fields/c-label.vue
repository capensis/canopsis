<template>
  <!-- eslint-disable-next-line vue/multiline-html-element-content-newline -->
  <v-label :aria-required="required" :class="labelClasses">{{ label ?? '' }}<slot />
    <span v-if="helpText" class="ml-1">
      <c-help-icon
        :text="helpText"
        icon="help"
        left
      />
    </span>
  </v-label>
</template>

<script>
import { computed } from 'vue';

export default {
  props: {
    required: {
      type: Boolean,
      default: false,
    },
    error: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      required: false,
    },
    helpText: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const labelClasses = computed(() => ({
      'error--text': props.error,
      'text-caption display-block': true,
    }));

    return {
      labelClasses,
    };
  },
};
</script>
