<template lang="pug">
  div
    v-tabs(centered, slider-color="primary")
      v-tab Alarm patterns
      v-tab Entity patterns
      v-tab-item
        patterns-list(
          v-field="form.alarm_patterns",
          :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS"
        )
      v-tab-item
        patterns-list(
          v-field="form.entity_patterns",
          :operators="$constants.WEBHOOK_EVENT_FILTER_RULE_OPERATORS"
        )
</template>

<script>
import formValidationHeaderMixin from '@/mixins/form/validation-header';
import formMixin from '@/mixins/form';

import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  components: {
    PatternsList,
  },
  mixins: [formMixin, formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>
