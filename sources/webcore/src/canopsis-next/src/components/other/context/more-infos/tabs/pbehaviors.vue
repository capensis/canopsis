<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      v-data-table.ma-0(:items="pbehaviors", :headers="headers")
        template(slot="items", slot-scope="props")
          td {{ props.item.name }}
          td {{ props.item.author }}
          td
            enabled-column(:value="props.item.enabled")
          td {{ props.item.tstart | timezone($system.timezone, 'long', true) }}
          td {{ props.item.tstop | timezone($system.timezone, 'long', true) }}
          td {{ props.item.type.name }}
          td {{ props.item.reason.name }}
          td {{ props.item.rrule }}
          td
            template(v-if="hasAccessToDeletePbehavior")
              v-btn(
                icon,
                small,
                @click.stop="showEditPbehaviorModal(props.item)"
              )
                v-icon edit
              v-btn(
                icon,
                small,
                @click="showDeletePbehaviorModal(props.item._id)"
              )
                v-icon(color="error") delete
</template>

<script>
import { MODALS, USERS_RIGHTS } from '@/constants';

import EnabledColumn from '@/components/tables/enabled-column.vue';

import authMixin from '@/mixins/auth';
import queryMixin from '@/mixins/query';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';

export default {
  components: {
    EnabledColumn,
  },
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
          value: 'isEnabled',
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
          value: 'type_',
        },
        {
          text: this.$t('pbehaviors.reason'),
          value: 'reason',
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
          afterSubmit: () => {
            this.fetchList();
            this.$popups.success({
              text: this.$t('success.default'),
            });
          },
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

