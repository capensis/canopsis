<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ $t('modals.createManualMetaAlarm.title') }}
      template(slot="text")
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="alarms")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-flex(xs12)
              manual-meta-alarm-form(v-model="form")
      template(slot="actions")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          type="submit"
        ) {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { get, isObject } from 'lodash';

import { MODALS, EVENT_ENTITY_TYPES } from '@/constants';

import { isWarningAlarmState } from '@/helpers/entities';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import { submittableMixin } from '@/mixins/submittable';
import { confirmableModalMixin } from '@/mixins/confirmable-modal';

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
  },
  components: {
    AlarmGeneralTable,
    ManualMetaAlarmForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerItemsMixin,
    eventActionsAlarmMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
      form: {
        manualMetaAlarm: null,
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
          data.ma_parents = [get(this.form.metaAlarm, 'entity._id')];
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