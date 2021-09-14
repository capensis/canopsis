<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.eventFilterRule.editActions') }}
    template(slot="text")
      event-filter-rule-actions-form(v-model="form.actions")
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import EventFilterRuleActionsForm from '@/components/other/event-filter/form/event-filter-rule-actions-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.eventFilterRuleActions,
  $_veeValidate: {
    validator: 'new',
  },
  components: { EventFilterRuleActionsForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { actions = [] } = this.modal.config;

    return {
      form: {
        actions: cloneDeep(actions),
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form.actions);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
