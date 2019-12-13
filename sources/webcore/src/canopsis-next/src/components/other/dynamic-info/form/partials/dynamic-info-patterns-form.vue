<template lang="pug">
  div
    v-alert(:value="hasAnyError", type="error") {{ $t('modals.createDynamicInfo.steps.patterns.validationError') }}
    v-tabs(ref="tabs", centered, slider-color="primary")
      v-tab {{ $t('modals.createDynamicInfo.steps.patterns.alarmPatterns') }}
      v-tab {{ $t('modals.createDynamicInfo.steps.patterns.entityPatterns') }}
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
import { isEmpty } from 'lodash';

import vuetifyTabsMixin from '@/mixins/vuetify/tabs';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  components: {
    PatternsList,
  },
  mixins: [vuetifyTabsMixin, formValidationHeaderMixin],
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
  created() {
    this.$validator.attach({
      name: 'patterns',
      rules: 'required:true',
      getter: () => {
        const isAlarmPatternsEmpty = isEmpty(this.form.alarm_patterns);
        const isEntityPatternsEmpty = isEmpty(this.form.entity_patterns);

        return !isAlarmPatternsEmpty || !isEntityPatternsEmpty;
      },
      context: () => this,
      vm: this,
    });
  },
};
</script>
