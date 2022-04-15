<template lang="pug">
  v-layout.white(column)
    v-layout(row, justify-end)
      c-action-fab-btn.ma-2(
        :tooltip="$t('modals.pbehaviorsCalendar.title')",
        icon="calendar_today",
        color="secondary",
        small,
        @click="showPbehaviorsCalendarModal"
      )
    v-data-table.ma-0(:items="pbehaviors", :headers="headers", :loading="pending")
      template(#items="{ item }")
        td {{ item.name }}
        td {{ item.author }}
        td
          c-enabled(:value="item.enabled")
        td {{ item.tstart | timezone($system.timezone) }}
        td {{ item.tstop | timezone($system.timezone) }}
        td {{ item.type.name }}
        td {{ item.reason.name }}
        td {{ item.rrule }}
        td
          v-layout(v-if="hasAccessToDeletePbehavior", row)
            c-action-btn(type="edit", @click="showEditPbehaviorModal(item)")
            c-action-btn(type="delete", @click="showDeletePbehaviorModal(item._id)")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { queryMixin } from '@/mixins/query';

const { mapActions } = createNamespacedHelpers('pbehavior');

export default {
  inject: ['$system'],
  mixins: [
    authMixin,
    queryMixin,
  ],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    tabId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      pbehaviors: [],
    };
  },
  computed: {
    hasAccessToDeletePbehavior() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.pbehaviorDelete);
    },

    headers() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.author'), value: 'author' },
        { text: this.$t('pbehaviors.isEnabled'), value: 'enabled' },
        { text: this.$t('pbehaviors.begins'), value: 'tstart' },
        { text: this.$t('pbehaviors.ends'), value: 'tstop' },
        { text: this.$t('pbehaviors.type'), value: 'type.type' },
        { text: this.$t('pbehaviors.reason'), value: 'reason.name' },
        { text: this.$t('pbehaviors.rrule'), value: 'rrule' },
        { text: this.$t('common.actionsLabel'), value: 'actionsLabel', sortable: false },
      ];
    },

    queryNonce() {
      return this.getQueryNonceById(this.tabId);
    },
  },
  watch: {
    queryNonce(value, oldValue) {
      if (value > oldValue) {
        this.fetchList();
      }
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchPbehaviorsByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
    }),

    showEditPbehaviorModal(pbehavior) {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          pbehaviors: [pbehavior],
          afterSubmit: this.fetchList,
        },
      });
    },

    showPbehaviorsCalendarModal() {
      this.$modals.show({
        name: MODALS.pbehaviorsCalendar,
        config: {
          title: this.$t('modals.pbehaviorsCalendar.entity.title', { name: this.entity.name }),
          entityId: this.entity._id,
        },
      });
    },

    showDeletePbehaviorModal(pbehaviorId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removePbehavior({ id: pbehaviorId });

            return this.fetchList();
          },
        },
      });
    },

    async fetchList() {
      try {
        this.pending = true;

        this.pbehaviors = await this.fetchPbehaviorsByEntityIdWithoutStore({ id: this.entity._id });
      } catch (err) {
        console.warn(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
