<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createManualMetaAlarm.title') }}
      template(#text="")
        v-layout(row)
          v-flex.text-xs-center
            alarm-general-table(:items="alarms")
        v-layout(row)
          v-divider.my-3
        v-layout(row)
          v-flex(xs12)
            manual-meta-alarm-form(v-model="form")
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

import { mapIds } from '@/helpers/array';
import { isAlarmStateNotOk } from '@/helpers/entities/alarm/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import ManualMetaAlarmForm from '@/components/widgets/alarm/forms/manual-meta-alarm-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to manage alarms in manual meta alarm
 */
export default {
  name: MODALS.createManualMetaAlarm,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    AlarmGeneralTable,
    ManualMetaAlarmForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: {
        metaAlarm: null,
        comment: '',
      },
    };
  },
  computed: {
    alarms() {
      return this.config.items.filter(isAlarmStateNotOk);
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          comment: this.form.comment,
          alarms: mapIds(this.alarms),
        };

        if (this.form.metaAlarm?._id) {
          data.id = this.form.metaAlarm?._id;
        } else {
          data.name = this.form.metaAlarm;
        }

        await this.config?.action?.(data);

        this.$modals.hide();
      }
    },
  },
};
</script>
