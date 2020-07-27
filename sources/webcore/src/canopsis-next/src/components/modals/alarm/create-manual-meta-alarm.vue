<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
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
import { isObject } from 'lodash';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  META_ALARM_EVENT_DEFAULT_FIELDS, ENTITIES_STATES,
} from '@/constants';
import { isWarningAlarmState } from '@/helpers/entities';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import submittableMixin from '@/mixins/submittable';

import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';
import ManualMetaAlarmForm from '@/components/other/alarm/forms/manual-meta-alarm-form.vue';

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
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin, submittableMixin()],
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
    /**
     * Function for data preparation
     *
     * @param {string} type - type of the event
     * @param {Array} items - item of the entity | Array of items of entity
     * @param {Object} data - data for the event
     * @returns {Object[]}
     */
    prepareData(type, items, data = {}) {
      return [{
        event_type: type,
        ma_children: items.map(({ entity }) => entity._id),

        ...data,
      }];
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          ...META_ALARM_EVENT_DEFAULT_FIELDS,

          output: this.form.output,
          state: ENTITIES_STATES.minor,
        };

        if (this.eventType === EVENT_ENTITY_TYPES.manualMetaAlarmUpdate) {
          data.ma_parents = [this.form.metaAlarm.d];
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
