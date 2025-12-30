<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <v-layout class="gap-2" column>
          <c-enabled-field v-model="form.enabled" />
          <template-testing-test-variables-wrapper
            v-model="form"
            :rule-id="scenarioId"
            :type="type"
          >
            <template #default="{ templateVars }">
              <scenario-form
                v-model="form"
                :template-vars="templateVars"
              />
            </template>
          </template-testing-test-variables-wrapper>
        </v-layout>
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="close"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref, inject } from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { formToScenario, scenarioToForm } from '@/helpers/entities/scenario/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import ScenarioForm from '@/components/other/scenario/form/scenario-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/test-variables/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createScenario,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ScenarioForm,
    TemplateTestingTestVariablesWrapper,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const type = TEMPLATE_TESTING_TEST_TYPES.scenario;

    const system = inject('$system');

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();

    const form = ref(scenarioToForm(config.value.scenario, system.timezone));

    const scenarioId = computed(() => config.value.scenario?._id);
    const title = computed(() => config.value.title ?? t('modals.createScenario.create.title'));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const scenario = await config.value.action?.(formToScenario(form.value, system.timezone));

        close();

        return scenario;
      },
    });

    useEntityInfoPropertyFetching();
    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      scenarioId,
      type,
      title,
      isDisabled,
      submitting,
      submit,
      close,
    };
  },
};
</script>
