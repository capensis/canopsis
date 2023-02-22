<template lang="pug">
  v-tabs(slider-color="primary", color="transparent", fixed-tabs, centered)
    v-tab(:class="{ 'error--text': hasGeneralError }") {{ $t('common.general') }}
    v-tab-item
      link-rule-general-form.mt-2(ref="general", v-field="form")
    v-tab(:class="{ 'error--text': hasSimpleError }") {{ $t('linkRule.simpleMode') }}
    v-tab-item
      link-rule-simple-form.mt-2(ref="simple", v-field="form")
    v-tab(:class="{ 'error--text': hasAdvancedError }") {{ $t('linkRule.advancedMode') }}
    v-tab-item
      link-rule-advanced-form.mt-2(ref="advanced", v-field="form.code")
</template>

<script>
import LinkRuleGeneralForm from './link-rule-general-form.vue';
import LinkRuleSimpleForm from './link-rule-simple-form.vue';
import LinkRuleAdvancedForm from './link-rule-advanced-form.vue';

export default {
  components: {
    LinkRuleGeneralForm,
    LinkRuleSimpleForm,
    LinkRuleAdvancedForm,
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
      hasSimpleError: false,
      hasAdvancedError: false,
    };
  },
  mounted() {
    this.$watch(() => this.$refs.general.hasAnyError, (value) => {
      this.hasGeneralError = value;
    });

    this.$watch(() => this.$refs.simple.hasAnyError, (value) => {
      this.hasSimpleError = value;
    });

    this.$watch(() => this.$refs.advanced.hasAnyError, (value) => {
      this.hasAdvancedError = value;
    });
  },
};
</script>
