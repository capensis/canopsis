<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      v-data-table.ma-0.pbehaviorsTable(:items="pbehaviors", :headers="headers")
        template(slot="items", slot-scope="props")
          td {{ props.item.name }}
          td {{ props.item.author }}
          td {{ props.item.connector }}
          td {{ props.item.connector_name }}
          td
            v-icon(
              small,
              :color="props.item.enabled ? 'primary' : 'error'"
            ) {{ props.item.enabled ? 'check' : 'clear'}}
          td {{ props.item.tstart | date('long') }}
          td {{ props.item.tstop | date('long') }}
          td {{ props.item.type_ }}
          td {{ props.item.reason }}
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

import authMixin from '@/mixins/auth';
import queryMixin from '@/mixins/query';
import entitiesPbehaviorMixin from '@/mixins/entities/pbehavior';
import entitiesPbehaviorCommentMixin from '@/mixins/entities/pbehavior/comment';

export default {
  mixins: [
    authMixin,
    queryMixin,
    entitiesPbehaviorMixin,
    entitiesPbehaviorCommentMixin,
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
          sortable: false,
        },
        {
          text: this.$t('common.author'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.connector'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.connectorName'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.isEnabled'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.begins'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.ends'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.type'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.reason'),
          sortable: false,
        },
        {
          text: this.$t('pbehaviors.rrule'),
          sortable: false,
        },
        {
          text: this.$t('common.actionsLabel'),
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
        name: MODALS.createPbehavior,
        config: {
          pbehavior,

          action: async (data) => {
            const { comments, ...preparedData } = data;

            await this.updatePbehavior({ data: preparedData, id: pbehavior._id });
            await this.updateSeveralPbehaviorComments({ pbehavior, comments });

            this.fetchList();
            this.$popups.addSuccess({ text: this.$t('success.default') });
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

