<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('modals.formConfirmation.title') }}
    template(slot="text")
      span.subheading {{ $t('modals.formConfirmation.text') }}
    template(slot="actions")
      v-btn(
        @click="$modals.hide"
      ) {{ $t('modals.formConfirmation.buttons.backToForm') }}
      v-btn.warning(
        :loading="submitting",
        :disabled="isDisabled",
        @click.prevent="submit(false)"
      ) {{ $t('modals.formConfirmation.buttons.dontSave') }}
      v-btn.primary(
        :loading="submitting",
        :disabled="isDisabled",
        @click.prevent="submit(true)"
      ) {{ $t('modals.formConfirmation.buttons.save') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.formConfirmation,
  components: { ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  methods: {
    async submit(submitted) {
      if (this.config.action) {
        await this.config.action(submitted);
      }

      this.$modals.hide();
    },
  },
};
</script>

