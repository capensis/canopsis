<template lang="pug">
  modal-wrapper(close)
    template(v-if="!config.hideTitle", #title="")
      span {{ title }}
    template(v-if="config.text", #text="")
      span.subheading {{ config.text }}
    template(#actions="")
      v-layout(wrap, justify-center)
        v-btn(
          :outline="$system.dark",
          color="error",
          @click="cancel"
        ) {{ $t('common.no') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          @click.prevent="submit"
        ) {{ $t('common.yes') }}
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
  inject: ['$system'],
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
      return this.config.title ?? this.$t('common.confirmation');
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
