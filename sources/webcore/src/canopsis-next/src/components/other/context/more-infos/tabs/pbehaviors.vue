<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      v-data-table.ma-0.pbehaviorsTable(:items="pbehaviors", :headers="pbehaviorsTableHeaders")
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
            v-btn(v-if="hasAccessToDeletePbehavior", @click="deletePbehavior(props.item._id)", icon, small)
              v-icon(color="error") delete
</template>

<script>
import { MODALS, USERS_RIGHTS } from '@/constants';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
import pbehaviorEntityMixin from '@/mixins/entities/pbehavior';

export default {
  mixins: [authMixin, modalMixin, pbehaviorEntityMixin],
  props: {
    itemId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      pbehaviorsTableHeaders: [
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
      ],
    };
  },
  computed: {
    hasAccessToDeletePbehavior() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.pbehaviorDelete);
    },
  },
  mounted() {
    this.fetchItems();
  },
  methods: {
    async deletePbehavior(itemId) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removePbehavior({ id: itemId });
            await this.fetchItems();
          },
        },
      });
    },
    fetchItems() {
      this.fetchPbehaviorsByEntityId({ id: this.itemId });
    },
  },
};
</script>

