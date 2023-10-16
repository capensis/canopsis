<template>
  <v-tabs
    slider-color="primary"
    color="transparent"
    fixed-tabs="fixed-tabs"
    centered="centered"
  >
    <v-tab :class="{ 'error--text': hasGeneralError }">
      {{ $t('common.general') }}
    </v-tab>
    <v-tab-item>
      <link-rule-general-form
        class="mt-2"
        ref="general"
        v-field="form"
      />
    </v-tab-item>
    <v-tab
      :class="{ 'error--text': hasSimpleError || errors.has('links') }"
      :disabled="sourceCodeWasChanged"
    >
      {{ $t('linkRule.simpleMode') }}
    </v-tab>
    <v-tab-item>
      <c-alert
        :value="errors.has('links')"
        transition="fade-transition"
        type="error"
      >
        {{ $t('linkRule.linksEmptyError') }}
      </c-alert>
      <link-rule-simple-form
        class="mt-2"
        ref="simple"
        v-field="form.links"
        :type="form.type"
        @input="resetLinksErrors"
      />
    </v-tab-item>
    <v-tab :class="{ 'error--text': hasAdvancedError || errors.has('links') }">
      {{ $t('linkRule.advancedMode') }}
    </v-tab>
    <v-tab-item>
      <c-alert
        :value="errors.has('links')"
        transition="fade-transition"
        type="error"
      >
        {{ $t('linkRule.linksEmptyError') }}
      </c-alert>
      <link-rule-advanced-form
        class="mt-2"
        ref="advanced"
        v-field="form.source_code"
        @input="resetLinksErrors"
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { isDefaultSourceCode } from '@/helpers/entities/link/form';

import LinkRuleGeneralForm from './link-rule-general-form.vue';
import LinkRuleSimpleForm from './link-rule-simple-form.vue';
import LinkRuleAdvancedForm from './link-rule-advanced-form.vue';

export default {
  inject: ['$validator'],
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
  computed: {
    sourceCodeWasChanged() {
      return !isDefaultSourceCode(this.form.source_code);
    },
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

    this.attachRequiredLinksRules();
  },
  beforeDestroy() {
    this.detachLinksRules();
  },
  methods: {
    resetLinksErrors() {
      this.$validator.reset({ name: 'links' });
    },

    attachRequiredLinksRules() {
      this.$validator.attach({
        name: 'links',
        rules: 'required:true',
        getter: () => !!this.form.links.length || !isDefaultSourceCode(this.form.source_code),
        vm: this,
      });
    },

    detachLinksRules() {
      this.$validator.detach('links');
    },
  },
};
</script>
