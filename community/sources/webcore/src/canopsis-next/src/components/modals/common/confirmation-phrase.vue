<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <v-alert
          color="warning"
          icon="info"
        >
          <span
            class="pre-line"
            v-html="config.text"
          />
        </v-alert>
        <div class="my-3">
          <p class="mb-2">
            {{ config.phraseText }}
          </p>
          <pre class="black--text grey lighten-2 d-inline pa-1">{{ originalPhrase }}</pre>
        </div>
        <v-text-field
          class="mt-2"
          v-model="phrase"
          :label="$t('modals.confirmationPhrase.phrase')"
        />
      </template>
      <template #actions="">
        <v-btn
          @click="$modals.hide"
          text
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :loading="submitting"
          :disabled="isDisabled || !phrasesEqual"
          type="submit"
        >
          {{ $t('common.yes') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS } from '@/constants';

import { submittableMixinCreator } from '@/mixins/submittable';
import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation phrase modal
 */
export default {
  name: MODALS.confirmationPhrase,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
  ],
  data() {
    return {
      phrase: '',
    };
  },
  computed: {
    originalPhrase() {
      return this.config.phrase;
    },

    phrasesEqual() {
      return this.phrase === this.originalPhrase;
    },
  },
  methods: {
    async submit() {
      if (this.phrasesEqual) {
        if (this.config.action) {
          await this.config.action();
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
