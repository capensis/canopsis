<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createCommentEvent.title') }}</span>
      </template>
      <template #text="">
        <v-layout class="gap-4" column>
          <template v-if="items.length">
            <alarm-general-table :items="items" />
          </template>
          <alarm-comment-template-form v-field="form" :templates="templates" />
        </v-layout>
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
          :loading="submitting"
          :disabled="isDisabled"
          class="primary"
          type="submit"
        >
          {{ $t('common.saveChanges') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { computed, ref } from 'vue';

import { MODALS } from '@/constants';

import { createCommentFormToCreateCommentEvent } from '@/helpers/entities/widget/forms/alarm';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import AlarmCommentTemplateForm from '@/components/other/comment-template/form/alarm-comment-template-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create comment an alarm
 */
export default {
  name: MODALS.createCommentEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, AlarmCommentTemplateForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const form = ref({
      template: null,
      comment: '',
    });

    const { config, close } = useInnerModal(props);

    const items = computed(() => config.value.items ?? []);
    const templates = computed(() => config.value.templates ?? []);

    const { submitting, isDisabled, submit } = useSubmittableForm({
      form,
      method: async () => {
        await config.value.action(createCommentFormToCreateCommentEvent(form.value));
        close();
      },
    });

    useFormConfirmableCloseModal({
      form,
      submit,
      close,
    });

    return {
      form,
      items,
      templates,
      submitting,
      isDisabled,
      submit,
    };
  },
};
</script>
