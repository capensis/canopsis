<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <template-testing-test-variables-wrapper
        v-field="form"
        :is-new="isNew"
        :type="type"
      >
        <template #default="{ templateVars }">
          <v-form>
            <link-rule-form
              v-model="form"
              :template-vars="templateVars"
            />
          </v-form>
        </template>
      </template-testing-test-variables-wrapper>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.cancel') }}
      </v-btn>
      <v-btn
        :disabled="isDisabled"
        :loading="submitting"
        class="primary"
        type="submit"
        @click="submit"
      >
        {{ $t('common.submit') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS, TEMPLATE_TESTING_TEST_TYPES, VALIDATION_DELAY } from '@/constants';

import { linkRuleToForm, formToLinkRule } from '@/helpers/entities/link/form';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useI18n } from '@/hooks/i18n';

import LinkRuleForm from '@/components/other/link-rule/form/link-rule-form.vue';
import TemplateTestingTestVariablesWrapper from '@/components/other/template-testing/template-testing-test-variables-wrapper.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createLinkRule,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    LinkRuleForm,
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
    const type = TEMPLATE_TESTING_TEST_TYPES.linkRule;

    const { config, close } = useInnerModal(props);
    const { t } = useI18n();

    const form = ref(linkRuleToForm(config.value.linkRule));

    const isNew = computed(() => !config.value.linkRule?._id);
    const title = computed(() => config.value.title ?? t('modals.createLinkRule.create.title'));

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToLinkRule(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      config,
      isNew,
      type,
      title,
      isDisabled,
      submitting,
      submit,
    };
  },
};
</script>
