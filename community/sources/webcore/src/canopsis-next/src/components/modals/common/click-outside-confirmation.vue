<template lang="pug">
  modal-wrapper
    template(#title="")
      span {{ $t('modals.clickOutsideConfirmation.title') }}
    template(#text="")
      span.subheading {{ $t('modals.clickOutsideConfirmation.text') }}
    template(#actions="")
      v-btn(
        depressed,
        flat,
        @click="$modals.hide"
      ) {{ $t('modals.clickOutsideConfirmation.buttons.backToForm') }}
      v-btn.warning(
        :loading="submitting",
        :disabled="isDisabled",
        @click.prevent="submit(false)"
      ) {{ $t('modals.clickOutsideConfirmation.buttons.dontSave') }}
      v-btn.primary(
        :loading="submitting",
        :disabled="isDisabled",
        @click.prevent="submit(true)"
      ) {{ $t('modals.clickOutsideConfirmation.buttons.save') }}
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
