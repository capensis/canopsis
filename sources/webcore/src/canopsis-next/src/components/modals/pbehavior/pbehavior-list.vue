<template lang="pug">
  modal-wrapper(data-test="pbehaviorListModal", close)
    template(slot="title")
      span {{ $t('alarmList.actions.titles.pbehaviorList') }}
    template(slot="text")
      advanced-data-table(:headers="headers", :items="filteredPbehaviors", expand)
        template(slot="enabled", slot-scope="props")
          enabled-column(:value="props.item.enabled")
        template(slot="tstart", slot-scope="props") {{ props.item.tstart | timezone($system.timezone, 'long', true) }}
        template(slot="tstop", slot-scope="props") {{ props.item.tstop | timezone($system.timezone, 'long', true) }}
        template(slot="type", slot-scope="props") {{ props.item.type | get('name', null, '') }}
        template(slot="reason", slot-scope="props") {{ props.item.reason | get('name', null, '') }}
        template(slot="expand", slot-scope="props")
          pbehaviors-list-expand-item(:pbehavior="props.item")
        template(slot="actions", slot-scope="props")
          v-layout(row)
            action-btn(
              v-if="hasAccessToEditPbehavior",
              type="edit",
              @click="showEditPbehaviorModal(props.item)"
            )
            action-btn(
              v-if="hasAccessToDeletePbehavior",
              type="delete",
              @click="showRemovePbehaviorModal(props.item._id)"
            )
        template(slot="is_active_status", slot-scope="props")
          v-icon(:color="props.item.is_active_status ? 'primary' : 'error'") $vuetify.icons.settings_sync
      v-layout(v-if="showAddButton", row)
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

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

import ActionBtn from '@/components/tables/action-btn.vue';
import EnabledColumn from '@/components/tables/enabled-column.vue';
import PbehaviorsListExpandItem from '@/components/other/pbehavior/exploitation/pbehaviors-list-expand-item.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal showing a list of an alarm's pbehaviors
 */
export default {
  components: {
    PbehaviorsListExpandItem,
    ActionBtn,
    EnabledColumn,
    ModalWrapper,
  },
  mixins: [modalInnerMixin, entitiesPbehaviorMixin],
  inject: ['$system'],
  computed: {
    headers() {
      return [
        { sortable: false, text: this.$t('tables.pbehaviorList.name'), value: 'name' },
        { sortable: false, text: this.$t('tables.pbehaviorList.author'), value: 'author' },
        { sortable: false, text: this.$t('tables.pbehaviorList.enabled'), value: 'enabled' },
        { sortable: false, text: this.$t('tables.pbehaviorList.tstart'), value: 'tstart' },
        { sortable: false, text: this.$t('tables.pbehaviorList.tstop'), value: 'tstop' },
        { sortable: false, text: this.$t('tables.pbehaviorList.type'), value: 'type' },
        { sortable: false, text: this.$t('tables.pbehaviorList.reason'), value: 'reason' },
        { sortable: false, text: this.$t('common.status'), value: 'is_active_status' },
        { sortable: false, text: this.$t('common.actionsLabel'), value: 'actions' },
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
        name: MODALS.pbehaviorPlanning,
        config: {
          pbehaviors: [pbehavior],

          afterSubmit: () => this.fetchPbehaviorsByEntityId({ id: this.modal.config.entityId }),
        },
      });
    },
  },
};
</script>
