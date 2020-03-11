<template lang="pug">
  div
    v-textarea(
      v-field="value.message",
      :label="$t('modals.createAction.fields.message')"
    )
    duration-field(v-field="value.duration")
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import DurationField from '@/components/forms/fields/duration.vue';

export default {
  inject: ['$validator'],
  components: {
    DurationField,
  },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
