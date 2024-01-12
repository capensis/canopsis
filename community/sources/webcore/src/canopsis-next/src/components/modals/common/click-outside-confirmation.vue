<template>
  <modal-wrapper>
    <template #title="">
      <span>{{ $t('modals.clickOutsideConfirmation.title') }}</span>
    </template>
    <template #text="">
      <span class="text-subtitle-1">{{ $t('modals.clickOutsideConfirmation.text') }}</span>
    </template>
    <template #actions="">
      <v-btn
        depressed
        text
        @click="$modals.hide"
      >
        {{ $t('modals.clickOutsideConfirmation.buttons.backToForm') }}
      </v-btn>
      <v-btn
        class="warning"
        :loading="submitting"
        :disabled="isDisabled"
        @click.prevent="submit(false)"
      >
        {{ $t('modals.clickOutsideConfirmation.buttons.dontSave') }}
      </v-btn>
      <v-btn
        class="primary"
        :loading="submitting"
        :disabled="isDisabled"
        @click.prevent="submit(true)"
      >
        {{ $t('modals.clickOutsideConfirmation.buttons.save') }}
      </v-btn>
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Click outside confirmation modal
 */
export default {
  name: MODALS.clickOutsideConfirmation,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
  ],
  methods: {
    async submit(confirmed) {
      await this.config.action(confirmed);

      this.$modals.hide();
    },
  },
};
</script>
