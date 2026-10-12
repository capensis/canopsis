<template>
  <v-form @submit.prevent="submit">
    <c-modal-wrapper close>
      <template #title>
        <span>{{ t('modals.createCommentEvent.title') }}</span>
      </template>

      <template #text>
        <c-layout gap="4" column>
          <alarm-general-table v-if="items.length" :items="items" />
          <alarm-comment-template-form
            :form="values"
            :templates="templates"
            @update:form="updateForm"
          />
        </c-layout>
      </template>

      <template #actions>
        <v-btn :disabled="submitting" variant="text" @click="close">
          {{ t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="submitting"
          color="primary"
          type="submit"
        >
          {{ t('common.saveChanges') }}
        </v-btn>
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { computed } from 'vue';

import { MODALS } from '@/constants';

import { createCommentFormToCreateCommentEvent } from '@/helpers/entities/widget/forms/alarm';

import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';
import { useForm } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';

import AlarmCommentTemplateForm from '@/components/other/comment-template/form/alarm-comment-template-form.vue';
import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

/**
 * Modal to create comment an alarm
 */
export default {
  name: MODALS.createCommentEvent,

  components: { AlarmGeneralTable, AlarmCommentTemplateForm },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = useForm({
      initialValues: {
        template: config.value.templates?.length ? config.value.templates[0] : null,
        comment: '',
      },
    });
    const { values, handleSubmit, isSubmitting: submitting } = form;

    const items = computed(() => config.value.items ?? []);
    const templates = computed(() => config.value.templates ?? []);

    const updateForm = (nextForm) => {
      values.value = nextForm;
    };

    const submit = handleSubmit(async () => {
      await config.value.action(createCommentFormToCreateCommentEvent(values.value));
      close();
    });

    useFormConfirmableCloseModal({
      form: values,
      submit,
      close,
    });

    return {
      t,
      values,
      items,
      templates,
      submitting,
      submit,
      close,
      updateForm,
    };
  },
};
</script>
