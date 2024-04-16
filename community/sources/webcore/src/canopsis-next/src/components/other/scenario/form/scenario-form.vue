<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <c-duration-field
      v-field="form.delay"
      :label="$t('common.delay')"
      :units-label="$t('common.unit')"
      name="delay"
      clearable
    />
    <c-enabled-field v-field="form.enabled" />
    <c-triggers-field
      :value="form.triggers"
      @input="updateField('triggers', $event)"
    />
    <c-disable-during-periods-field v-field="form.disable_during_periods" />
    <c-priority-field v-field="form.priority" />
    <v-tabs
      slider-color="primary"
      centered
    >
      <v-tab :class="{ 'error--text': hasParametersChildrenError }">
        {{ $tc('common.action', 2) }}
      </v-tab>
      <v-tab>{{ $t('common.testQuery') }}</v-tab>

      <v-tab-item eager>
        <scenario-actions-form
          v-field="form.actions"
          class="mt-2"
          name="actions"
        />
      </v-tab-item>
      <v-tab-item>
        <v-layout class="mt-2">
          <c-alert :value="!isWebhookActionExist" type="error">
            {{ $t('scenario.errors.testQueryRequireSteps') }}
          </c-alert>
          <scenario-test-query v-if="isWebhookActionExist" :form="form" />
        </v-layout>
      </v-tab-item>
    </v-tabs>
  </v-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { isWebhookActionType } from '@/helpers/entities/action';

import { formMixin } from '@/mixins/form';

import ScenarioTestQuery from '@/components/other/scenario/partials/scenario-test-query.vue';

import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import ScenarioActionsForm from './scenario-actions-form.vue';

export default {
  inject: ['$validator'],
  components: { ScenarioTestQuery, ScenarioActionsForm },
  mixins: [formMixin],
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
    const isWebhookActionExist = computed(() => props.form.actions.some(({ type }) => isWebhookActionType(type)));

    const parametersFormElement = ref();
    const { hasChildrenError: hasParametersChildrenError } = useValidationElementChildren(parametersFormElement);

    return {
      isWebhookActionExist,
      parametersFormElement,
      hasParametersChildrenError,
    };
  },
};
</script>
