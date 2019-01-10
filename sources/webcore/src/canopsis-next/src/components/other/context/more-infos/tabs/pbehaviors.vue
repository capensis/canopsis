<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      v-data-table.ma-0.pbehaviorsTable(:items="items", :headers="pbehaviorsTableHeaders")
        template(slot="items", slot-scope="props")
          td {{ props.item.name }}
          td {{ props.item.author }}
          td {{ props.item.connector }}
          td {{ props.item.connector_name }}
          td {{ props.item.enabled }}
          td {{ props.item.tstart | date('long') }}
          td {{ props.item.tstop | date('long') }}
          td {{ props.item.type_ }}
          td {{ props.item.reason }}
          td {{ props.item.rrule }}
          td
            v-btn.error--text(@click="deletePbehavior(props.item._id)", icon, small)
              v-icon delete
</template>

<script>
import modalMixin from '@/mixins/modal';
import pbehaviorEntityMixin from '@/mixins/entities/pbehavior';
import { MODALS } from '@/constants';

export default {
  mixins: [modalMixin, pbehaviorEntityMixin],
  props: {
    itemId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      items: [],
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
    async fetchItems() {
      await this.fetchPbehaviorsByEntityId({ id: this.itemId });

      if (this.pbehaviorItems) {
        this.items = [...this.pbehaviorItems];
      }
    },
  },
};
</script>

