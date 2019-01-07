<template lang="pug">
  div
    v-card
      v-card-title.primary.white--text
        v-layout(justify-space-between, align-center)
          span.headline {{ $t('modals.createPbehavior.title') }}
      v-card-text
        v-data-table(:headers="headers", :items="firstItem.pbehaviors", disable-initial-sort, hide-actions)
          template(slot="items", slot-scope="props")
            td(v-for="key in fields")
              span(
              v-if="key === 'tstart' || key === 'tstop'",
              key="key"
              ) {{ props.item[key] | date('long') }}
              span(v-else) {{ props.item[key] }}
            td
              v-btn.mx-0(@click="showRemovePbehaviorModal(props.item._id)", icon)
                v-icon delete
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import modalInnerItemsMixin from '@/mixins/modal/inner-items';
import { MODALS } from '@/constants';

const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('pbehavior');

/**
 * Modal showing a list of an alarm's pbehaviors
 */
export default {
  name: MODALS.pbehaviorList,

  mixins: [modalInnerItemsMixin],
  data() {
    const fields = [
      'name',
      'author',
      'enabled',
      'tstart',
      'tstop',
      'type_',
      'reason',
      'rrule',
    ];

    const headers = fields.map(v => ({ sortable: false, text: this.$t(`tables.pbehaviorList.${v}`) }));

    headers.push({ sortable: false, text: this.$t('common.actionsLabel') });

    return {
      fields,
      headers,
    };
  },
  methods: {
    ...pbehaviorMapActions({
      removePbehavior: 'remove',
    }),

    showRemovePbehaviorModal(pbehaviorId) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.removePbehavior({ id: pbehaviorId }),
        },
      });
    },
  },
};
</script>
