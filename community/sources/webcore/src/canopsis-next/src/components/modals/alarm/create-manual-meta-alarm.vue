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
import { isObject } from 'lodash';

import { MODALS, EVENT_ENTITY_TYPES, VALIDATION_DELAY } from '@/constants';

import { isWarningAlarmState } from '@/helpers/entities';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { modalInnerItemsMixin } from '@/mixins/modal/inner-items';
import { eventActionsAlarmMixin } from '@/mixins/event-actions/alarm';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import AlarmGeneralTable from '@/components/widgets/alarm/alarm-general-list.vue';
import ManualMetaAlarmForm from '@/components/widgets/alarm/forms/manual-meta-alarm-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to manage alarms in meta alarm
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
    eventActionsAlarmMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: {
        metaAlarm: null,
        output: '',
      },
    };
  },
  computed: {
    alarms() {
      return this.items.filter(isWarningAlarmState);
    },

    eventType() {
      return isObject(this.form.metaAlarm)
        ? EVENT_ENTITY_TYPES.manualMetaAlarmUpdate
        : EVENT_ENTITY_TYPES.manualMetaAlarmGroup;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = { output: this.form.output };

        if (this.eventType === EVENT_ENTITY_TYPES.manualMetaAlarmUpdate) {
          data.ma_parents = [this.form.metaAlarm?._id];
        } else {
          data.display_name = this.form.metaAlarm;
        }

        await this.createEvent(this.eventType, this.alarms, data);

        this.$modals.hide();
      }
    },
  },
};
</script>
