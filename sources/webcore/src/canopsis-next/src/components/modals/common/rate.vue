<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        v-layout(justify-center)
          span.subheading {{ config.text }}
        rate-form(v-model="form")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import { MODALS } from '@/constants';

import RateForm from '@/components/forms/rate.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.rate,
  components: { ModalWrapper, RateForm },
  mixins: [
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: {
        comment: '',
        rating: 5,
      },
    };
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(this.form);
      }

      this.$modals.hide();
    },
  },
};
</script>

