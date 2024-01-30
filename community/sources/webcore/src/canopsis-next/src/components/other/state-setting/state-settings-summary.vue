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
      <template v-else>
        {{ stateMethodName }}
      </template>
    </span>
    <v-expand-transition>
      <span v-if="!stateSettingPending">{{ stateMethodSummaryText }}</span>
    </v-expand-transition>
  </v-layout>
</template>

<script>

import { createNamespacedHelpers } from 'vuex';

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
    stateMethodName() {
      /**
       * TODO: Should be changed on real data
       */
      return '%rule name%';
    },

    stateMethodSummaryText() {
      /**
       * TODO: Should be changed on real data
       */
      return 'TODO: Should be implemented';
    },
  },
  mounted() {
    this.checkStateSetting({
      // name: this.entity.name,
      // type: this.entity.type,
      // infos: this.entity.infos,
      impact_level: this.entity.impact_level,
    });
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
