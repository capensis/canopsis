<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <pbehavior-patterns-form v-model="form" />
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import PbehaviorPatternsForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-patterns-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.pbehaviorPatterns,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    ModalWrapper,
    PbehaviorPatternsForm,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { patterns } = this.modal.config;

    return {
      form: cloneDeep(patterns),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.patterns.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
