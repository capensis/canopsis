<template>
  <div class="position-relative">
    <c-compiled-template :template="template" :sanitize-options="sanitizeOptions" />
  </div>
</template>

<script>
import { DEFAULT_SANITIZE_OPTIONS } from '@/config';

import { setSeveralFields } from '@/helpers/immutable';

export default {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    sanitizeOptions() {
      return setSeveralFields(DEFAULT_SANITIZE_OPTIONS, {
        allowedTags: tags => [...tags, 'iframe'],
        'allowedAttributes.iframe': [
          'allowtransparency',
          'frameborder',
          'hspace',
          'marginheight',
          'marginwidth',
          'sandbox',
          'scrolling',
          'seamless',
          'src',
          'srcdoc',
          'vspace',
        ],
      });
    },

    template() {
      return this.widget.parameters?.template;
    },
  },
};
</script>
