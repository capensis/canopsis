<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.clickOutsideConfirmation.title') }}
    template(slot="text")
      span.subheading {{ $t('modals.clickOutsideConfirmation.text') }}
    template(slot="actions")
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

import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Click outside confirmation modal
 */
export default {
  name: MODALS.clickOutsideConfirmation,
  components: { ModalWrapper },
  mixins: [submittableMixin()],
  methods: {
    async submit(confirmed) {
      await this.config.action(confirmed);

      this.$modals.hide();
    },
  },
};
</script>
