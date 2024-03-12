<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ config.title }}
      template(#text="")
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-text-field(
              v-model="form.comment",
              v-validate="'required'",
              :label="$t('modals.createEvent.fields.output')",
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
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { mapIds } from '@/helpers/entities';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to remove alarms from meta alarm
 */
export default {
  name: MODALS.removeAlarmsFromMetaAlarm,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [
    modalInnerMixin,
    modalInnerItemsMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: {
        comment: '',
      },
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          ...this.form,

          alarms: mapIds(this.items),
        };

        await this.config?.action?.(data);

        this.$modals.hide();
      }
    },
  },
};
</script>
