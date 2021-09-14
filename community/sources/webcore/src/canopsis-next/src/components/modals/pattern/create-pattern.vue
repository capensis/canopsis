<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.eventFilterRule.editPattern') }}
      template(slot="text")
        pattern-form(
          v-model="form",
          :operators="operators",
          :only-simple-rule="config.onlySimpleRule"
        )
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled || patternWasChanged",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, get } from 'lodash';

import { MODALS, EVENT_FILTER_RULE_OPERATORS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import PatternForm from '@/components/other/pattern/form/pattern-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPattern,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PatternForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { pattern = {} } = this.modal.config;

    return {
      form: cloneDeep(pattern),
      activeTab: 0,
    };
  },
  computed: {
    operators() {
      return this.config.operators || EVENT_FILTER_RULE_OPERATORS;
    },

    patternWasChanged() {
      return get(this.fields, ['pattern', 'changed']);
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
