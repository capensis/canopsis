<template lang="pug">
  modal-wrapper(data-test="pbehaviorListModal")
    template(slot="title")
      span {{ $t('alarmList.actions.titles.pbehaviorList') }}
    template(slot="text")
      v-data-table(:headers="headers", :items="filteredPbehaviors", disable-initial-sort)
        template(slot="items", slot-scope="props")
          td(v-for="key in fields", :key="key")
            span(
              v-if="key === 'tstart' || key === 'tstop'",
            ) {{ props.item[key] | date('long') }}
            span(v-else) {{ props.item[key] }}
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

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal showing a list of an alarm's pbehaviors
 */
export default {
  name: MODALS.pbehaviorList,
  components: { ModalWrapper },
  mixins: [modalInnerMixin, entitiesPbehaviorMixin, entitiesPbehaviorCommentMixin],
  data() {
    const fields = [
      'name',
      'author',
      'enabled',
      'tstart',
      'tstop',
      'type_',
      'reason',
    ];

    const headers = fields.map(v => ({ sortable: false, text: this.$t(`tables.pbehaviorList.${v}`) }));

    headers.push({ sortable: false, text: this.$t('common.actionsLabel') });

    return {
      fields,
      headers,
    };
  },
  computed: {
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
        return this.pbehaviors.filter(value => value.isActive);
      }

      return this.pbehaviors;
    },
    showAddButton() {
      return this.modal.config.showAddButton;
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
        name: MODALS.createPbehavior,
        config: {
          pbehavior: {
            filter: {
              _id: { $in: [this.modal.config.entityId] },
            },
          },
          action: data => this.createPbehavior({ data }),
        },
      });
    },

    showEditPbehaviorModal(pbehavior) {
      this.$modals.show({
        name: MODALS.createPbehavior,
        config: {
          pbehavior,

          action: async (data) => {
            const { comments, ...preparedData } = data;

            await this.updatePbehavior({ data: preparedData, id: pbehavior._id });
            await this.updateSeveralPbehaviorComments({ pbehavior, comments });

            await this.fetchPbehaviorsByEntityId({ id: this.modal.config.entityId });
          },
        },
      });
    },
  },
};
</script>
