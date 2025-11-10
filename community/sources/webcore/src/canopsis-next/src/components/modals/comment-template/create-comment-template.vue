<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <comment-template-form v-model="form" />
      </template>
      <template #actions="">
        <v-btn
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

import { MODALS } from '@/constants';

import { commentTemplateToForm, formToCommentTemplate } from '@/helpers/entities/comment-template/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import CommentTemplateForm from '@/components/other/comment-template/form/comment-template-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createCommentTemplate,
  $_veeValidate: {
    validator: 'new',
  },
  components: { CommentTemplateForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(commentTemplateToForm(config.value.template));

    const title = computed(() => config.value.title || t('modals.createCommentTemplate.create.title'));

    const { submitting, isDisabled, submit } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action?.(formToCommentTemplate(form.value));

        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      submitting,
      isDisabled,

      submit,
      close,
    };
  },
};
</script>
