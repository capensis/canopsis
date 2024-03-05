<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <patterns-form
          v-model="form"
          v-bind="patternsProps"
        />
      </template>
      <template #actions="">
        <v-btn
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
import { omit } from 'lodash';

import { MODALS, PATTERNS_FIELDS, VALIDATION_DELAY } from '@/constants';

import { filterToForm, formToFilter } from '@/helpers/entities/filter/form';

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
          await this.config.action(formToFilter(this.form, this.getPatternsFields(), this.config.corporate));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
