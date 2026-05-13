<template>
  <v-layout column>
    <c-enabled-field v-field="form.enabled" with-background />
    <c-form-general-patterns-tabs>
      <template #general="{ setRef }">
        <idle-rule-general-form
          v-field="form"
          :ref="setRef"
          :is-entity-type="isEntityType"
        />
      </template>
      <template #patterns="{ setRef }">
        <idle-rule-patterns-form
          v-field="form.patterns"
          :ref="setRef"
          :is-entity-type="isEntityType"
        />
      </template>
    </c-form-general-patterns-tabs>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { isIdleRuleEntityType } from '@/helpers/entities/idle-rule/form';

import IdleRuleGeneralForm from './idle-rule-general-form.vue';
import IdleRulePatternsForm from './idle-rule-patterns-form.vue';

export default {
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
  setup(props) {
    const isEntityType = computed(() => isIdleRuleEntityType(props.form.type));

    return {
      isEntityType,
    };
  },
};
</script>
