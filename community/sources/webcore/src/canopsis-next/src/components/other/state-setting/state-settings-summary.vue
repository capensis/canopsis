<template>
  <v-layout column>
    <span>
      {{ $t('stateSetting.computeMethod') }}:
      <v-progress-circular
        v-if="stateSettingPending"
        class="ml-1"
        color="primary"
        size="20"
        width="3"
        indeterminate
      />
      <b v-else>{{ stateMethodName }}</b>
    </span>
    <v-expand-transition>
      <span v-if="!stateSettingPending && stateMethodSummaryText">{{ stateMethodSummaryText }}</span>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { JUNIT_STATE_SETTING_METHODS } from '@/constants';

import { infosToArray } from '@/helpers/entities/shared/form';
import { isEntityEventsStateSettings } from '@/helpers/entities/entity/entity';

const { mapActions: mapEntityActions } = createNamespacedHelpers('entity');

export default {
  props: {
    entity: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      stateSettingPending: false,
      stateSetting: {},
    };
  },
  computed: {
    isEventsStateSettings() {
      return isEntityEventsStateSettings(this.entity);
    },

    stateMethodName() {
      if (this.isEventsStateSettings) {
        return this.$tc('common.event', 2);
      }

      return this.stateSetting.title || this.$t(`stateSetting.junit.methods.${JUNIT_STATE_SETTING_METHODS.worst}`);
    },

    stateMethodSummaryText() {
      if (this.stateSetting.title) {
        return 'TODO: Should be implemented';
      }

      return '';
    },
  },
  mounted() {
    if (!this.isEventsStateSettings) {
      this.checkStateSetting({
        name: this.entity.name,
        type: this.entity.type,
        infos: infosToArray(this.entity.infos),
        impact_level: this.entity.impact_level,
      });
    }
  },
  methods: {
    ...mapEntityActions({
      checkEntityStateSetting: 'checkStateSetting',
    }),

    async checkStateSetting(data) {
      try {
        this.stateSettingPending = true;
        this.stateSetting = await this.checkEntityStateSetting({ data });
      } catch (err) {
        console.error(err);
      } finally {
        this.stateSettingPending = false;
      }
    },
  },
};
</script>
