<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ config.title }}</span>
    </template>
    <template #text="">
      <c-information-block :title="$t('view.sharedViewUrl')">
        <v-text-field
          :value="config.url"
          readonly
        >
          <template #append="">
            <c-copy-btn
              :value="config.url"
              :tooltip="$t('common.copyLink')"
              top
              @success="showCopySuccessPopup"
              @error="showCopyErrorPopup"
            />
          </template>
        </v-text-field>
      </c-information-block>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('common.close') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.shareView,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
  ],
  methods: {
    showCopySuccessPopup() {
      this.$popups.success({ text: this.$t('success.linkCopied') });
    },

    showCopyErrorPopup() {
      this.$popups.error({ text: this.$t('errors.default') });
    },
  },
};
</script>
