<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <tag-form
          v-model="form"
          :is-imported="isImported"
          :is-new="isNew"
        />
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
import { ref, computed } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { tagToForm, formToTag } from '@/helpers/entities/tag/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';
import { useFormConfirmableCloseModal } from '@/hooks/confirmable-modal';

import TagForm from '@/components/other/tag/form/tag-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createTag,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { TagForm, ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { config, close } = useInnerModal(props);

    const form = ref(tagToForm(config.value.tag));

    const isNew = computed(() => !config.value.tag?._id);
    const title = computed(() => (
      config.value.title || t('modals.createTag.create.title')
    ));

    const isImported = computed(() => config.value.isImported);

    const { submit, isDisabled, submitting } = useSubmittableForm({
      form,
      method: async () => {
        if (config.value.action) {
          await config.value.action(formToTag(form.value));
        }
        close();
      },
    });

    useFormConfirmableCloseModal({ form, submit, close });

    return {
      form,
      title,
      isNew,
      isImported,
      isDisabled,
      submitting,
      submit,
      close,
    };
  },
};
</script>
