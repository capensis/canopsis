<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      v-data-table.ma-0(:items="pbehaviors", :headers="headers")
        template(slot="items", slot-scope="props")
          td {{ props.item.name }}
          td {{ props.item.author }}
          td
            c-enabled(:value="props.item.enabled")
          td {{ props.item.tstart | timezone($system.timezone, 'long', true) }}
          td {{ props.item.tstop | timezone($system.timezone, 'long', true) }}
          td {{ props.item.type.name }}
          td {{ props.item.reason.name }}
          td {{ props.item.rrule }}
          td
            v-layout(v-if="hasAccessToDeletePbehavior", row)
              c-action-btn(type="edit", @click="showEditPbehaviorModal(props.item)")
              c-action-btn(type="delete", @click="showDeletePbehaviorModal(props.item._id)")
</template>

<script>
import { MODALS, USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';
import queryMixin from '@/mixins/query';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

export default {
  inject: ['$system'],
  mixins: [
    authMixin,
    queryMixin,
    entitiesPbehaviorMixin,
  ],
  props: {
    itemId: {
      type: String,
      required: true,
    },
    tabId: {
      type: String,
      required: true,
    },
  },
  computed: {
    hasAccessToDeletePbehavior() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.pbehaviorDelete);
    },

    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.author'),
          value: 'author',
        },
        {
          text: this.$t('pbehaviors.isEnabled'),
          value: 'enabled',
        },
        {
          text: this.$t('pbehaviors.begins'),
          value: 'tstart',
        },
        {
          text: this.$t('pbehaviors.ends'),
          value: 'tstop',
        },
        {
          text: this.$t('pbehaviors.type'),
          value: 'type.type',
        },
        {
          text: this.$t('pbehaviors.reason'),
          value: 'reason.name',
        },
        {
          text: this.$t('pbehaviors.rrule'),
          value: 'rrule',
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actionsLabel',
          sortable: false,
        },
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
    showEditPbehaviorModal(pbehavior) {
      this.$modals.show({
        name: MODALS.pbehaviorPlanning,
        config: {
          pbehaviors: [pbehavior],
          afterSubmit: this.fetchList,
        },
      });
    },

    showDeletePbehaviorModal(pbehaviorId) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removePbehavior({ id: pbehaviorId });

            this.fetchList();
          },
        },
      });
    },
    fetchList() {
      this.fetchPbehaviorsByEntityId({ id: this.itemId });
    },
  },
};
</script>

