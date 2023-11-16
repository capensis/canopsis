<template>
  <v-layout column>
    <c-enabled-field v-field="form.enabled" />
    <v-tabs
      slider-color="primary"
      centered
    >
      <v-tab :class="{ 'error--text': hasGeneralError }">
        {{ $t('common.general') }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasPatternsError }">
        {{ $tc('common.pattern') }}
      </v-tab>

      <v-tab-item eager>
        <idle-rule-general-form
          ref="general"
          v-field="form"
          :is-entity-type="isEntityType"
        />
      </v-tab-item>
      <v-tab-item eager>
        <idle-rule-patterns-form
          class="mt-2"
          ref="patterns"
          v-field="form.patterns"
          :is-entity-type="isEntityType"
        />
      </v-tab-item>
    </v-tabs>
  </v-layout>
</template>

<script>
import { isIdleRuleEntityType } from '@/helpers/entities/idle-rule/form';

import IdleRuleGeneralForm from './idle-rule-general-form.vue';
import IdleRulePatternsForm from './idle-rule-patterns-form.vue';

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
