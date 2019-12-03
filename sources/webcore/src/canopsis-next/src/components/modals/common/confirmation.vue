<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('common.confirmation') }}
    template(slot="actions")
      v-layout(wrap, justify-center)
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          data-test="submitButton",
          @click.prevent="submit"
        ) {{ $t('common.yes') }}
        v-btn.error(
          data-test="cancelButton",
          @click="$modals.hide"
        ) {{ $t('common.no') }}
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
  name: MODALS.confirmation,
  components: { ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action();
      }

      this.$modals.hide();
    },
  },
};
</script>

