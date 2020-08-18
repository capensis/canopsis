<template lang="pug">
  modal-wrapper(data-test="pbehaviorListModal")
    template(slot="title")
      span {{ $t('alarmList.actions.titles.pbehaviorList') }}
    template(slot="text")
      v-data-table(:headers="headers", :items="filteredPbehaviors", disable-initial-sort)
        template(slot="items", slot-scope="props")
          td {{ props.item.name }}
          td {{ props.item.author }}
          td
            enabled-column(:value="props.item.enabled")
          td {{ props.item.tstart | timezone($system.timezone, 'long', true) }}
          td {{ props.item.tstop | timezone($system.timezone, 'long', true) }}
          td {{ props.item.type.name }}
          td {{ props.item.reason.name }}
          td
            v-btn.mx-0(
              v-for="action in availableActions",
              :key="action.name",
              :data-test="`pbehaviorRow-${props.item._id}-action-${action.name}`",
              @click="() => action.action(props.item)",
              icon
            )
              v-icon {{ action.icon }}
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

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import entitiesPbehaviorCommentMixin from '@/mixins/entities/pbehavior/comment';

import EnabledColumn from '@/components/tables/enabled-column.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal showing a list of an alarm's pbehaviors
 */
export default {
  components: { EnabledColumn, ModalWrapper },
  mixins: [modalInnerMixin, entitiesPbehaviorMixin, entitiesPbehaviorCommentMixin],
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
        { sortable: false, text: this.$t('common.actionsLabel') },
      ];
    },
    availableActions() {
      const availableActions = this.modal.config.availableActions || [];

      return availableActions.reduce((acc, action) => {
        if (action === CRUD_ACTIONS.delete) {
          acc.push({
            name: CRUD_ACTIONS.delete,
            icon: 'delete',
            action: pbehavior => this.showRemovePbehaviorModal(pbehavior._id),
          });
        }

        if (action === CRUD_ACTIONS.update) {
          acc.push({
            name: CRUD_ACTIONS.update,
            icon: 'edit',
            action: pbehavior => this.showEditPbehaviorModal(pbehavior),
          });
        }

        return acc;
      }, []);
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
