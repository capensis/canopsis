<template lang="pug">
  modal-wrapper(data-test="confirmationModal", close)
    template(v-if="!config.hideTitle", slot="title")
      span {{ title }}
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.confirmation,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
  ],
  data() {
    return {
      submitted: false,
      cancelled: false,
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('common.confirmation');
    },
  },
  beforeDestroy() {
    if (!this.submitted && this.config.cancel) {
      this.config.cancel(this.cancelled);
    }
  },
  methods: {
    cancel() {
      this.cancelled = true;

      this.$modals.hide();
    },
    async submit() {
      if (this.config.action) {
        await this.config.action();
      }

      this.submitted = true;
      this.$modals.hide();
    },
  },
};
</script>
