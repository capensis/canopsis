<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        patterns-form(v-model="form", v-bind="patternsProps")
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
import { omit } from 'lodash';

import { MODALS, PATTERNS_FIELDS, VALIDATION_DELAY } from '@/constants';

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
    delay: VALIDATION_DELAY,
  },
  components: { PatternsForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: filterToForm(this.modal.config.filter, this.getPatternsFields()),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createFilter.create.title');
    },

    patternsProps() {
      return omit(this.config, ['title', 'action']);
    },

    patternsFields() {
      return this.getPatternsFields();
    },
  },
  methods: {
    getPatternsFields() {
      const { withAlarm, withEntity, withPbehavior, withEvent, withServiceWeather } = this.modal.config;

      return [
        withAlarm && PATTERNS_FIELDS.alarm,
        withEntity && PATTERNS_FIELDS.entity,
        withPbehavior && PATTERNS_FIELDS.pbehavior,
        withEvent && PATTERNS_FIELDS.event,
        withServiceWeather && PATTERNS_FIELDS.serviceWeather,
      ].filter(Boolean);
    },

    async submit() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToFilter(this.form, this.patternsFields, this.config.corporate));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
