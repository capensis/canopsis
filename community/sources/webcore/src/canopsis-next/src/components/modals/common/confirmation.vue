<template>
  <modal-wrapper close>
    <template
      v-if="!config.hideTitle"
      #title=""
    >
      <span>{{ title }}</span>
    </template>
    <template
      v-if="config.text"
      #text=""
    >
      <span class="text-subtitle-1 pre-wrap">{{ config.text }}</span>
    </template>
    <template #actions="">
      <v-layout
        wrap
        justify-center
      >
        <v-btn
          :outlined="$system.dark"
          color="error"
          @click="cancel"
        >
          {{ $t('common.no') }}
        </v-btn>
        <v-btn
          :loading="submitting"
          :disabled="isDisabled"
          class="ml-2"
          color="primary"
          @click.prevent="submit"
        >
          {{ $t('common.yes') }}
        </v-btn>
      </v-layout>
    </template>
  </modal-wrapper>
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
