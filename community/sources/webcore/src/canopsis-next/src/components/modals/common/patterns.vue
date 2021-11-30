<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title)
        span {{ title }}
      template(#text)
        patterns-form(
          v-model="form",
          :name="config.name",
          :entity="config.entity",
          :alarm="config.alarm",
          :event="config.event",
          :total-entity="config.totalEntity"
        )
      template(#actions)
        v-btn(
          :disabled="submitting",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import PatternsForm from '@/components/forms/patterns-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.patterns,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PatternsForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: this.modal.config.patterns ? cloneDeep(this.modal.config.patterns) : {},
    };
  },
  computed: {
    title() {
      return this.config.title || this.$tc('common.pattern', 2);
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
