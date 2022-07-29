<template lang="pug">
  v-layout(column)
    c-enabled-field(v-field="form.enabled")
    v-tabs(slider-color="primary", color="transparent", fixed-tabs, centered)
      v-tab(:class="{ 'error--text': hasGeneralError }") {{ $t('common.general') }}
      v-tab-item
        idle-rule-general-form(
          ref="general",
          v-field="form",
          :is-entity-type="isEntityType"
        )
      v-tab(:class="{ 'error--text': hasPatternsError }") {{ $tc('common.pattern') }}
      v-tab-item
        idle-rule-patterns-form.mt-2(
          ref="patterns",
          v-field="form.patterns",
          :is-entity-type="isEntityType"
        )
</template>

<script>
import { isIdleRuleEntityType } from '@/helpers/forms/idle-rule';

import IdleRuleGeneralForm from './partials/idle-rule-general-form.vue';
import IdleRulePatternsForm from './partials/idle-rule-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    IdleRuleGeneralForm,
    IdleRulePatternsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      hasGeneralError: false,
      hasPatternsError: false,
    };
  },
  computed: {
    isEntityType() {
      return isIdleRuleEntityType(this.form.type);
    },
  },
  mounted() {
    this.$watch(() => this.$refs.general.hasAnyError, (value) => {
      this.hasGeneralError = value;
    });

    this.$watch(() => this.$refs.patterns.hasAnyError, (value) => {
      this.hasPatternsError = value;
    });
  },
};
</script>
