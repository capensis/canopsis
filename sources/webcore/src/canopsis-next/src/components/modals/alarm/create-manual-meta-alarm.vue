<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createManualMetaAlarm.title') }}
      template(slot="text")
        v-container
          v-layout(row)
            v-flex.text-xs-center
              alarm-general-table(:items="items")
          v-layout(row)
            v-divider.my-3
          v-layout(row)
            v-combobox(
              v-model="form.metaAlarm",
              v-validate="'required'",
              :items="manualMetaAlarms",
              :label="$t('modals.createManualMetaAlarm.fields.metaAlarm')",
              :error-messages="errors.collect('manualMetaAlarm')",
              :loading="pending",
              item-value="d",
              item-text="v.display_name",
              name="manualMetaAlarm",
              return-object,
              blur-on-create
            )
              template(slot="no-data")
                v-list-tile
                  v-list-tile-content
                    v-list-tile-title(v-html="$t('modals.createManualMetaAlarm.noData')")
          v-layout(row)
            v-text-field(
              v-model="form.output",
              :label="$t('modals.createManualMetaAlarm.fields.output')"
            )
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
import { createNamespacedHelpers } from 'vuex';

import {
  MODALS,
  EVENT_ENTITY_TYPES,
  MANUAL_META_ALARMS_REQUEST_FILTER,
  META_ALARM_EVENT_DEFAULT_FIELDS,
} from '@/constants';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import eventActionsAlarmMixin from '@/mixins/event-actions/alarm';
import submittableMixin from '@/mixins/submittable';

import AlarmGeneralTable from '@/components/other/alarm/alarm-general-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('alarm');

/**
 * Modal to cancel an alarm
 */
export default {
  name: MODALS.createManualMetaAlarm,
  $_veeValidate: {
    validator: 'new',
  },
  components: { AlarmGeneralTable, ModalWrapper },
  mixins: [modalInnerItemsMixin, eventActionsAlarmMixin, submittableMixin()],
  data() {
    return {
      pending: false,
      manualMetaAlarms: [],
      form: {
        manualMetaAlarm: null,
        output: '',
      },
    };
  },
  computed: {
    eventType() {
      return isObject(this.form.metaAlarm)
        ? EVENT_ENTITY_TYPES.manualMetaAlarmUpdate
        : EVENT_ENTITY_TYPES.manualMetaAlarmGroup;
    },
  },
  mounted() {
    this.fetchManualMetaAlarms();
  },
  methods: {
    ...mapActions({
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

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

    async fetchManualMetaAlarms() {
      this.pending = true;

      const { alarms = [] } = await this.fetchAlarmsListWithoutStore({
        params: { filter: MANUAL_META_ALARMS_REQUEST_FILTER },
      });

      this.manualMetaAlarms = alarms;
      this.pending = false;
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          ...META_ALARM_EVENT_DEFAULT_FIELDS,

          output: this.form.output,
        };

        if (this.eventType === EVENT_ENTITY_TYPES.manualMetaAlarmUpdate) {
          data.ma_parents = [this.form.metaAlarm.d];
        } else {
          data.display_name = this.form.metaAlarm;
        }

        await this.createEvent(this.eventType, this.items, data);

        this.$modals.hide();
      }
    },
  },
};
</script>
