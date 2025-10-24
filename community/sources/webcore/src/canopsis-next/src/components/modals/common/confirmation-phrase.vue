<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <c-alert
          color="warning"
          icon="info"
        >
          <span
            v-html="config.text"
            class="pre-line text-body-2"
          />
        </c-alert>
        <div class="my-3">
          <p v-html="sanitizedPhraseText" class="mb-2" />
          <pre class="black--text grey lighten-2 d-inline pa-1">{{ originalPhrase }}</pre>
        </div>
        <v-text-field
          v-model="phrase"
          :label="$t('modals.confirmationPhrase.phrase')"
          class="mt-2"
          autofocus
        />
        <component
          v-if="config.additionalForm?.component"
          :is="config.additionalForm.component"
          v-model="additionalForm"
        />
      </template>
      <template #actions="">
        <v-btn
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="isDisabled || !phrasesEqual"
          class="primary"
          type="submit"
        >
          {{ $t('common.yes') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { ref, computed } from 'vue';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';

import { useInnerModal } from '@/hooks/modals';
import { useSubmittableForm } from '@/hooks/submittable-form';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation phrase modal
 */
export default {
  name: MODALS.confirmationPhrase,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ModalWrapper },
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { config, close } = useInnerModal(props);

    const phrase = ref('');
    const additionalForm = ref(config.value.additionalForm?.form);

    const originalPhrase = computed(() => config.value.phrase);
    const phrasesEqual = computed(() => phrase.value === originalPhrase.value);
    const sanitizedPhraseText = computed(() => sanitizeHtml(config.value.phraseText || ''));

    const { submitting, isDisabled, submit } = useSubmittableForm({
      form: { phrase },
      method: async () => {
        if (phrasesEqual.value) {
          await config.value.action?.(additionalForm.value);

          close();
        }
      },
    });

    return {
      phrase,
      additionalForm,
      config,
      originalPhrase,
      phrasesEqual,
      sanitizedPhraseText,
      submitting,
      isDisabled,
      submit,
    };
  },
};
</script>
