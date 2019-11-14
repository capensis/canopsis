<template lang="pug">
  modal-wrapper
    template(slot="title")
      span {{ $t('common.confirmation') }}
    template(slot="actions")
      v-btn.primary(
        :loading="submitting",
        :disabled="submitting",
        data-test="submitButton",
        @click.prevent="submit"
      ) {{ $t('common.yes') }}
      v-btn.error(@click="$modals.hide") {{ $t('common.no') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.confirmation,
  components: { ModalWrapper },
  mixins: [modalInnerMixin],
  data() {
    return {
      submitting: false,
    };
  },
  methods: {
    async submit() {
      this.submitting = true;

      if (this.config.action) {
        await this.config.action();
      }

      this.$modals.hide();
      this.submitting = false;
    },
  },
};
</script>

