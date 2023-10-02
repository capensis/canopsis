<template lang="pug">
  v-list-group
    template(#activator="")
      v-list-tile
        span(:class="validationHeaderClass") {{ $tc('common.map') }}
    v-container
      c-map-field(v-field="value", name="parameters.map")
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formValidationHeaderMixin],
  props: {
    value: {
      type: String,
      required: false,
    },
  },
};
</script>
