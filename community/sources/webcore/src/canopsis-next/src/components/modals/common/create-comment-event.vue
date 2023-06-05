<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createCommentEvent.title') }}
      template(#text="")
        v-layout(column)
          template(v-if="items.length")
            alarm-general-table(:items="items")
            v-divider.my-3
          v-text-field(
            v-model="form.output",
            v-validate="'required'",
            :label="$tc('common.comment')",
            :error-messages="errors.collect('comment')",
            name="comment"
          )
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          type="submit"
        ) {{ $t('common.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create comment an alarm
 */
export default {
  name: MODALS.createCommentEvent,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: {
        output: '',
      },
    };
  },
  computed: {
    items() {
      return this.config.items ?? [];
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action({ ...this.form });

        if (this.config && this.config.afterSubmit) {
          await this.config.afterSubmit();
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
