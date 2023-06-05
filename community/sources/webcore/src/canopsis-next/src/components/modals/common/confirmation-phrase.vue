<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ config.title }}
      template(#text="")
        v-alert(:value="true", color="warning", icon="info")
          span.pre-line(v-html="config.text")
        div.my-3
          p.mb-2 {{ config.phraseText }}
          pre.black--text.grey.lighten-2.d-inline.pa-1 {{ originalPhrase }}
        v-text-field.mt-2(v-model="phrase", :label="$t('modals.confirmationPhrase.phrase')")
      template(#actions="")
        v-btn(@click="$modals.hide", flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled || !phrasesEqual",
          type="submit"
        ) {{ $t('common.yes') }}
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
