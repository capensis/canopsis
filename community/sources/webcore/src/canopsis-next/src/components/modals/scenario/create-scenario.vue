<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        {{ title }}
      </template>
      <template #text="">
        <scenario-form v-model="form" />
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
import { computed, ref } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { formToScenario, scenarioToForm, scenarioErrorToForm } from '@/helpers/entities/scenario/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';
import { useInfo } from '@/hooks/store/modules/info';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { usePopups } from '@/hooks/popups';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';

import ScenarioForm from '@/components/other/scenario/form/scenario-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createScenario,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ScenarioForm,
    ModalWrapper,
  },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);
    const { timezone } = useInfo();
    const popups = usePopups();

    const form = ref(scenarioToForm(config.value.scenario, timezone.value));

    const { setFormErrors } = useValidationFormErrors(form);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        try {
          if (config.value.action) {
            await config.value.action(formToScenario(form.value, timezone.value));
          }

          close();
        } catch (err) {
          if (err.error) {
            popups.error({ text: err.error });
          } else {
            setFormErrors(scenarioErrorToForm(err, form.value));
          }

          throw err;
        }
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });
    useEntityInfoPropertyFetching();

    const title = computed(() => config.value.title ?? t('modals.createScenario.create.title'));

    return {
      config,

      form,

      isDisabled,
      submitting,

      title,

      submit,
      close,
    };
  },
};
</script>
