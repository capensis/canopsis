<template lang="pug">
  v-layout.white(column)
    v-layout(row, justify-end)
      c-action-fab-btn.ma-2(
        :tooltip="$t('modals.pbehaviorsCalendar.title')",
        icon="calendar_today",
        color="secondary",
        small,
        left,
        @click="showPbehaviorsCalendarModal"
      )
    c-advanced-data-table.ma-0(
      :items="pbehaviors",
      :headers="headers",
      :loading="pending"
    )
      template(#enabled="{ item }")
        c-enabled(:value="item.enabled")
      template(#tstart="{ item }") {{ item.tstart | timezone($system.timezone) }}
      template(#tstop="{ item }") {{ item.tstop | timezone($system.timezone) }}
      template(#rrule="{ item }")
        v-icon {{ item.rrule ? 'check' : 'clear' }}
      template(#icon="{ item }")
        v-icon(color="primary") {{ item.type.icon_name }}
      template(#actions="{ item }")
        v-layout(row)
          c-action-btn(v-if="editable", type="edit", @click="showEditPbehaviorModal(item)")
          c-action-btn(v-if="deletable", type="delete", @click="showDeletePbehaviorModal(item._id)")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import Observer from '@/services/observer';

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
  props: {
    entity: {
      type: Object,
      required: true,
    },
    deletable: {
      type: Boolean,
      default: false,
    },
    editable: {
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
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.author'), value: 'author.name' },
        { text: this.$t('pbehaviors.isEnabled'), value: 'enabled' },
        { text: this.$t('pbehaviors.begins'), value: 'tstart' },
        { text: this.$t('pbehaviors.ends'), value: 'tstop' },
        { text: this.$t('pbehaviors.type'), value: 'type.name' },
        { text: this.$t('pbehaviors.reason'), value: 'reason.name' },
        { text: this.$t('pbehaviors.rrule'), value: 'rrule' },
        { text: this.$t('common.icon'), value: 'icon' },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
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
      removePbehavior: 'removeWithoutStore',
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
