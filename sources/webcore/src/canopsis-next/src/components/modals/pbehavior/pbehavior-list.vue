<template lang="pug">
  div
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline {{ $t('alarmList.actions.titles.pbehaviorList') }}
      v-card-text
        v-data-table(:headers="headers", :items="filteredPbehaviors", disable-initial-sort)
          template(slot="items", slot-scope="props")
            td(v-for="key in fields")
              span(
                v-if="key === 'tstart' || key === 'tstop'",
                key="key"
              ) {{ props.item[key] | date('long') }}
              span(v-else) {{ props.item[key] }}
            td
              v-btn.mx-0(
                v-for="action in availableActions",
                :key="action.name",
                @click="() => action.action(props.item)",
                icon
              )
                v-icon {{ action.icon }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn.primary(@click="$modals.hide") {{ $t('common.ok') }}
</template>

<script>

import { MODALS, CRUD_ACTIONS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import entitiesPbehaviorCommentMixin from '@/mixins/entities/pbehavior/comment';

/**
 * Modal showing a list of an alarm's pbehaviors
 */
export default {
  name: MODALS.pbehaviorList,

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
