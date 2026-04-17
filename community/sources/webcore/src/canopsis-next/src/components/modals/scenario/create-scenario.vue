<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper text-class="position-relative" close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
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
        <ai-chat-sidebar
          v-if="chatShown"
          v-bind="chatOptions.bind"
          v-on="chatOptions.on"
        />
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
          :disabled="isDisabled || chatOptions.bind.pending"
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
import {
  computed,
  ref,
  inject,
  set,
  toRef,
} from 'vue';

import { LLM_SOCKET_CONTEXTS, MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { formToScenario, scenarioToForm } from '@/helpers/entities/scenario/form';

import { useAiChatForm } from '@/hooks/ai/ai-chat-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import AiChatSidebar from '@/components/other/llm/chat/ai-chat-sidebar.vue';
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
    AiChatSidebar,
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

    const formRef = computed({
      get: () => form.value.actions,
      set: actions => set(form, 'actions', actions),
    });

    const {
      shown: chatShown,
      options: chatOptions,
    } = useAiChatForm({
      form: formRef,

      modal: toRef(props, 'modal'),
      ruleId: props.modal.config?.scenario?._id,
      context: LLM_SOCKET_CONTEXTS.scenario,
    });

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        const result = await config.value.action?.(formToScenario(form.value, system.timezone));

        await config.value.afterSubmit?.(result);

        close();

        return result;
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
      chatShown,
      chatOptions,
      submit,
      close,
    };
  },
};
</script>
