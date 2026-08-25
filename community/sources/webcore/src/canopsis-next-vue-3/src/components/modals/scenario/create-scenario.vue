<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper text-class="position-relative" close>
      <template #title>
        <span>{{ title }}</span>
      </template>

      <template #text>
        <scenario-form :rule-id="scenarioId" />
      </template>

      <template #actions>
        <v-btn variant="text" @click="closeModal">
          {{ t('common.cancel') }}
        </v-btn>

        <v-btn
          :disabled="isSubmitting"
          :loading="isSubmitting"
          color="primary"
          type="submit"
        >
          {{ submitLabel }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed, inject } from 'vue';

import {
  formToScenario,
  scenarioErrorToForm,
  scenarioToForm,
} from '@/helpers/entities/scenario/form';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useEntityInfoPropertyFetching } from '@/hooks/store/modules/entity-info-property';

import ScenarioForm from '@/components/other/scenario/form/scenario-form.vue';

export default {
  components: {
    ScenarioForm,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();

    const system = inject('$system', {});

    const config = computed(() => props.modal.config ?? {});
    const scenario = computed(() => config.value.scenario);
    const scenarioId = computed(() => scenario.value?._id);
    const isNew = computed(() => !scenarioId.value);
    const title = computed(() => config.value.title || t('modals.createScenario.create.title'));
    const submitLabel = computed(() => t(isNew.value ? 'common.create' : 'common.save'));

    const form = useForm({
      initialValues: scenarioToForm(scenario.value, system.timezone),
      errorsToValidation: errors => scenarioErrorToForm(errors, form.values.value),
    });
    const { values, handleSubmit, isSubmitting } = form;

    const closeModal = () => modals.hide({ id: props.modal.id });

    const submit = handleSubmit(async () => {
      const result = await config.value.action?.(formToScenario(values.value, system.timezone));

      await config.value.afterSubmit?.(result);

      closeModal();

      return result;
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close: closeModal,
    });
    
    useEntityInfoPropertyFetching();

    return {
      t,
      title,
      scenarioId,
      submitLabel,
      isSubmitting,
      submit,
      closeModal,
    };
  },
};
</script>
