<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        patterns-form(v-model="form", alarm, entity)
          template(slot="additionalTabs")
            v-tab Pbehaviors types
            v-tab-item
              v-layout(row)
                pbehavior-type-field(
                  v-model="form.active_on_pbh",
                  label="Active on types",
                  :is-item-disabled="isActiveItemDisabled",
                  chips,
                  multiple
                )
              v-layout(row)
                pbehavior-type-field(
                  v-model="form.disabled_on_pbh",
                  label="Disabled on types",
                  :is-item-disabled="isDisabledItemDisabled",
                  chips,
                  multiple
                )
      template(slot="actions")
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

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';
import validationErrorsMixin from '@/mixins/form/validation-errors';

import PatternsForm from '@/components/forms/patterns.vue';
import PbehaviorTypeField from '@/components/other/pbehavior/calendar/partials/pbehavior-type-field.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.patterns,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PbehaviorTypeField, ModalWrapper, PatternsForm },
  mixins: [
    modalInnerMixin,
    submittableMixin(),
    validationErrorsMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: this.modal.config.patterns ? cloneDeep(this.modal.config.patterns) : {},
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.patterns.title');
    },
  },
  methods: {
    isActiveItemDisabled(item) {
      return this.form.disabled_on_pbh.includes(item._id);
    },

    isDisabledItemDisabled(item) {
      return this.form.active_on_pbh.includes(item._id);
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(this.form);
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>
