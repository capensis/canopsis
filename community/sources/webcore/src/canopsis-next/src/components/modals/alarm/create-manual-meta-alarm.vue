<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.createManualMetaAlarm.title') }}
      template(#text="")
        v-container
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

import { isWarningAlarmState, mapIds } from '@/helpers/entities';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { entitiesManualMetaAlarmMixin } from '@/mixins/entities/manual-meta-alarm';
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
    modalInnerItemsMixin,
    entitiesManualMetaAlarmMixin,
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
      return this.items.filter(isWarningAlarmState);
    },

    metaAlarmId() {
      return this.form.metaAlarm?._id;
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

        if (this.metaAlarmId) {
          await this.addAlarmsIntoManualMetaAlarm({ id: this.metaAlarmId, data });
        } else {
          data.name = this.form.metaAlarm;

          await this.createManualMetaAlarm({ data });
        }

        if (this.config.afterSubmit) {
          await this.config.afterSubmit();
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
