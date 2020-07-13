<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ config.title }}
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
              :items="metaAlarms",
              label="Meta alarm",
              :error-messages="errors.collect('metaAlarm')",
              :loading="pending",
              item-value="id",
              item-text="title",
              name="metaAlarm",
              return-object,
              blur-on-create
            )
              template(slot="no-data")
                v-list-tile
                  v-list-tile-content
                    v-list-tile-title(v-html="$t('modals.createMetaAlarm.noData')")
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
import { isString } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { MODALS, MANUAL_META_ALARMS_REQUEST_FILTER } from '@/constants';

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
      metaAlarms: [],
      form: {
        metaAlarm: null,
      },
    };
  },
  mounted() {
    this.fetchMetaAlarms();
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
        connector: 'conn',
        connector_name: 'conname',
        source_type: 'resource',
        event_type: type,
        component: 'test-01',
        resource: 'rew1',
        state: 1,
        output: data.output,
        ma_children: items.map(({ entity }) => entity._id),
      }];
    },

    async fetchMetaAlarms() {
      this.pending = true;

      const { alarms = [] } = await this.fetchAlarmsListWithoutStore({
        params: { filter: MANUAL_META_ALARMS_REQUEST_FILTER },
      });

      this.metaAlarms = alarms;
      this.pending = false;
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (isString(this.form.metaAlarm)) {
          // TODO: create a new one meta alarm
        } else {
          // TODO: add alarms to meta alarm
        }

        // await this.createEvent(this.config.eventType, this.items, data);

        // this.$modals.hide();
      }
    },
  },
};
</script>
