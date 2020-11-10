<template lang="pug">
  modal-wrapper(data-test="confirmationModal")
    template(v-if="!config.hideTitle", slot="title")
      span {{ $t('common.confirmation') }}
    template(v-if="config.text", slot="text")
      span.subheading {{ config.text }}
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
          @click="cancel"
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
  inject: ['$clickOutside'],
  components: { ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  created() {
    this.$clickOutside.register(this.cancelHandler);
  },
  beforeDestroy() {
    this.$clickOutside.unregister(this.cancelHandler);
  },
  methods: {
    async cancelHandler() {
      if (this.config.cancel) {
        await this.config.cancel();
      }
    },

    async submit() {
      if (this.config.action) {
        await this.config.action();
      }

      this.$modals.hide();
    },

    async cancel() {
      await this.cancelHandler();
      this.$modals.hide();
    },
  },
};
</script>

