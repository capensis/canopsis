<template lang="pug">
  modal-wrapper(data-test="pbehaviorListModal", close)
    template(slot="title")
      span {{ $t('alarmList.actions.titles.pbehaviorList') }}
    template(slot="text")
      c-advanced-data-table(:headers="headers", :items="filteredPbehaviors", expand)
        template(slot="enabled", slot-scope="props")
          c-enabled(:value="props.item.enabled")
        template(slot="tstart", slot-scope="props") {{ props.item.tstart | timezone($system.timezone) }}
        template(slot="tstop", slot-scope="props") {{ props.item.tstop | timezone($system.timezone) }}
        template(slot="rrule", slot-scope="props")
          v-icon {{ props.item.rrule ? 'check' : 'clear' }}
        template(slot="expand", slot-scope="props")
          pbehaviors-list-expand-item(:pbehavior="props.item")
        template(slot="actions", slot-scope="props")
          v-layout(row)
            c-action-btn(
              v-if="hasAccessToEditPbehavior",
              type="edit",
              @click="showEditPbehaviorModal(props.item)"
            )
            c-action-btn(
              v-if="hasAccessToDeletePbehavior",
              type="delete",
              @click="showRemovePbehaviorModal(props.item._id)"
            )
        template(slot="is_active_status", slot-scope="props")
          v-icon(:color="props.item.is_active_status ? 'primary' : 'error'") $vuetify.icons.settings_sync
      v-layout(v-if="showAddButton", justify-end)
        v-btn(
          icon,
          fab,
          small,
          color="primary",
          @click="showCreatePbehaviorModal"
        )
          v-icon add
    template(slot="actions")
      v-btn.primary(data-test="pbehaviorListConfirmButton", @click="$modals.hide") {{ $t('common.ok') }}
</template>

<script>
import { MODALS, CRUD_ACTIONS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import PbehaviorsListExpandItem from '@/components/other/pbehavior/exploitation/pbehaviors-list-expand-item.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal showing a list of an alarm's pbehaviors
 */
export default {
  inject: ['$system'],
  components: {
    PbehaviorsListExpandItem,
    ModalWrapper,
  },
  mixins: [modalInnerMixin, entitiesPbehaviorMixin],
  computed: {
    headers() {
      return [
        { text: this.$t('common.name'), value: 'name' },
        { text: this.$t('common.author'), value: 'author' },
        { text: this.$t('pbehaviors.isEnabled'), value: 'enabled' },
        { text: this.$t('pbehaviors.begins'), value: 'tstart' },
        { text: this.$t('pbehaviors.ends'), value: 'tstop' },
        { text: this.$t('pbehaviors.type'), value: 'type.name' },
        { text: this.$t('pbehaviors.reason'), value: 'reason.name' },
        { text: this.$t('pbehaviors.rrule'), value: 'rrule' },
        { text: this.$t('common.status'), value: 'is_active_status', sortable: false },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
    availableActions() {
      return this.modal.config.availableActions || [];
    },

    hasAccessToEditPbehavior() {
      return this.availableActions.includes(CRUD_ACTIONS.update);
    },

    hasAccessToDeletePbehavior() {
      return this.availableActions.includes(CRUD_ACTIONS.delete);
    },

    filteredPbehaviors() {
      if (this.modal.config.onlyActive) {
        return this.pbehaviors.filter(value => value.enabled);
      }

      return this.pbehaviors;
    },
    showAddButton() {
      return this.modal.config.availableActions.includes(CRUD_ACTIONS.create);
    },
  },
  mounted() {
    this.fetchPbehaviorsByEntityId({ id: this.modal.config.entityId });
  },
  methods: {
    showRemovePbehaviorModal(pbehaviorId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.removePbehavior({ id: pbehaviorId }),
        },
      });
    },

    showCreatePbehaviorModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          filter: {
            _id: { $in: [this.modal.config.entityId] },
          },
        },
      });
    },

    showEditPbehaviorModal(pbehavior) {
      this.$modals.show({
        name: MODALS.createPbehavior,
        config: {
          pbehavior,
          noFilter: true,
          timezone: this.$system.timezone,
          action: data => this.updatePbehavior({ data, id: pbehavior._id }),
        },
      });
    },
  },
};
</script>
