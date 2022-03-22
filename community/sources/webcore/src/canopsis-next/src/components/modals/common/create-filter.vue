<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        patterns-form(
          v-model="form",
          :name="config.name",
          :entity="config.entity",
          :pbehavior="config.pbehavior",
          :alarm="config.alarm",
          :event="config.event"
        )
      template(#actions="")
        v-btn(
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
import { MODALS } from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/forms/filter';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import PatternsForm from '@/components/forms/patterns-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createFilter,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PatternsForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: filterToForm(this.modal.config.filter),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createFilter.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToFilter(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
