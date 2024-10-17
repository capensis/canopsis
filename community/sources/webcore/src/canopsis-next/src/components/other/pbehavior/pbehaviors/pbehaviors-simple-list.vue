<template>
  <v-layout column>
    <v-layout class="gap-2" justify-end>
      <c-action-fab-btn
        v-if="addable"
        :tooltip="$t('modals.createPbehavior.create.title')"
        icon="add"
        color="primary"
        left
        @click="showCreatePbehaviorModal"
      />
      <c-action-fab-btn
        :tooltip="$t('modals.pbehaviorsCalendar.title')"
        icon="calendar_today"
        color="secondary"
        left
        @click="showPbehaviorsCalendarModal"
      />
    </v-layout>
    <c-advanced-data-table
      :items="pbehaviors"
      :headers="headers"
      :loading="pending"
      :dense="dense"
    >
      <template #enabled="{ item }">
        <c-enabled :value="item.enabled" />
      </template>
      <template #tstart="{ item }">
        {{ formatIntervalDate(item, 'tstart') }}
      </template>
      <template #tstop="{ item }">
        {{ formatIntervalDate(item, 'tstop') }}
      </template>
      <template #rrule_end="{ item }">
        {{ formatRruleEndDate(item) }}
      </template>
      <template #rrule="{ item }">
        <v-icon>{{ item.rrule ? 'check' : 'clear' }}</v-icon>
      </template>
      <template #icon="{ item }">
        <v-icon color="primary">
          {{ item.type.icon_name }}
        </v-icon>
      </template>
      <template #status="{ item }">
        <v-icon :color="item.is_active_status ? 'primary' : 'error'">
          $vuetify.icons.settings_sync
        </v-icon>
      </template>
      <template #actions="{ item }">
        <pbehavior-actions
          :pbehavior="item"
          :removable="removable"
          :updatable="updatable"
          @refresh="fetchList"
        />
      </template>
    </c-advanced-data-table>
  </v-layout>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import Observer from '@/services/observer';

import { createEntityIdPatternByValue } from '@/helpers/entities/pattern/form';

import { pbehaviorsDateFormatMixin } from '@/mixins/pbehavior/pbehavior-date-format';

import PbehaviorActions from './partials/pbehavior-actions.vue';

const { mapActions } = createNamespacedHelpers('pbehavior');

export default {
  inject: {
    $system: {},
    $periodicRefresh: {
      default() {
        return new Observer();
      },
    },
  },
  components: { PbehaviorActions },
  mixins: [pbehaviorsDateFormatMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: false,
    },
    withActiveStatus: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: false,
      pbehaviors: [],
    };
  },
  computed: {
    headers() {
      const headers = [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.author'), value: 'author.display_name' },
        { text: this.$t('pbehavior.isEnabled'), value: 'enabled' },
        { text: this.$t('pbehavior.begins'), value: 'tstart' },
        { text: this.$t('pbehavior.ends'), value: 'tstop' },
        { text: this.$t('pbehavior.rruleEnd'), value: 'rrule_end' },
        { text: this.$t('common.recurrence'), value: 'rrule' },
        { text: this.$t('common.type'), value: 'type.name' },
        { text: this.$t('common.reason'), value: 'reason.name' },
        { text: this.$tc('common.icon', 1), value: 'icon' },
      ];

      if (this.withActiveStatus) {
        headers.push({ text: this.$t('common.status'), value: 'status', sortable: false });
      }

      if (this.updatable || this.removable) {
        headers.push({ text: this.$t('common.actionsLabel'), value: 'actions', sortable: false });
      }

      return headers;
    },
  },
  mounted() {
    this.fetchList();

    this.$periodicRefresh.register(this.fetchList);
  },
  beforeDestroy() {
    this.$periodicRefresh.unregister(this.fetchList);
  },
  methods: {
    ...mapActions({
      fetchPbehaviorsByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
    }),

    showPbehaviorsCalendarModal() {
      this.$modals.show({
        name: MODALS.pbehaviorsCalendar,
        config: {
          title: this.$t('modals.pbehaviorsCalendar.entity.title', { name: this.entity.name }),
          entityId: this.entity._id,
        },
      });
    },

    showCreatePbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          entityPattern: createEntityIdPatternByValue(this.entity._id),
          entities: [this.entity],
          afterSubmit: this.fetchList,
        },
      });
    },

    async fetchList() {
      try {
        this.pending = true;

        this.pbehaviors = await this.fetchPbehaviorsByEntityIdWithoutStore({
          id: this.entity._id,
          params: {
            with_flags: true,
          },
        });
      } catch (err) {
        console.warn(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
