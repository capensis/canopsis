<template lang="pug">
  div
    v-card
      v-card-title
        span.headline {{ $t('modals.createPbehavior.title') }}
      v-card-text
        v-data-table(:headers="headers", :items="item.pbehaviors", disable-initial-sort, hide-actions)
          template(slot="headerCell", slot-scope="props")
            span {{ $t(props.header.text) }}
          template(slot="items", slot-scope="props")
            td(v-for="key in fields")
              span(
              v-if="key === 'tstart' || key === 'tstop'",
              key="key"
              ) {{ props.item[key] | moment("DD/MM/YYYY HH:mm:ss") }}
              span(v-else) {{ props.item[key] }}
            td
              v-btn.mx-0(@click="remove(props.item._id)", icon)
                v-icon delete
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapGetters: modalMapGetters } = createNamespacedHelpers('modal');
const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');
const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('entities/pbehavior');

export default {
  data() {
    const fields = [
      'name',
      'author',
      'connector',
      'connector_name',
      'enabled',
      'tstart',
      'tstop',
      'type_',
      'reason',
      'rrule',
    ];

    const headers = fields.map(v => ({ sortable: false, text: `tables.pbehaviorList.${v}` }));

    headers.push({ sortable: false, text: 'common.actionsLabel' });

    return {
      fields,
      headers,
    };
  },
  computed: {
    ...modalMapGetters(['config']),
    ...entitiesMapGetters(['getItem']),
    item() {
      return this.getItem(this.config.itemType, this.config.itemId);
    },
  },
  methods: {
    ...pbehaviorMapActions({
      removePbehavior: 'remove',
    }),

    async remove(id) {
      await this.removePbehavior({ id });
    },
  },
};
</script>
