<template>
  <div class="position-relative">
    <c-compiled-template :template="template" :sanitize-options="sanitizeOptions" />
  </div>
</template>

<script>
import { TEXT_WIDGET_SANITIZE_OPTIONS } from '@/constants/sanitize';

export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    sanitizeOptions() {
      return TEXT_WIDGET_SANITIZE_OPTIONS;
    },

    template() {
      return this.widget.parameters?.template;
    },
  },
};
</script>
