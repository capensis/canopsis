<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile.widget-settings-item-title
        v-layout(:class="validationHeaderClass", align-center)
          slot(name="title") {{ title }}
          span.font-italic.caption.ml-1(v-if="optional || subTitle")
            span(v-if="optional") ({{ $t('common.optional') }})
            span(v-if="subTitle") {{ subTitle }}
    v-container
      slot
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  props: {
    title: {
      type: String,
      required: false,
    },
    subTitle: {
      type: String,
      required: false,
    },
    optional: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
