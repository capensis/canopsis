<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $tc('common.pattern', 2) }}
    v-tab-item(lazy)
      v-layout.py-3
        v-flex(xs12, md8, offset-md2)
          v-card
            v-card-text
              tag-patterns-form(:form="patterns", readonly)
</template>

<script>
import { tagFilterPatternsToForm } from '@/helpers/entities/tag/form';

import TagPatternsForm from '@/components/other/tag/form/fields/tag-patterns-form.vue';

export default {
  components: { TagPatternsForm },
  props: {
    tag: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    patterns() {
      return tagFilterPatternsToForm(this.tag);
    },
  },
};
</script>
