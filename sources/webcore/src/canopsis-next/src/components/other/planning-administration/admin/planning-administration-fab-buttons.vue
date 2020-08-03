<template lang="pug">
  fab-expand-buttons(@refresh="$emit('refresh')")
    v-tooltip(v-if="hasCreateAnyTypeAccess", top)
      v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateTypeModal")
        v-icon list
      span {{ $t('modals.createPbehaviorType.title') }}
</template>

<script>
import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';

import FabExpandButtons from '@/components/other/fab-buttons/fab-expand-buttons.vue';
import { MODALS } from '@/constants';

export default {
  components: {
    FabExpandButtons,
  },
  mixins: [
    rightsTechnicalPbehaviorTypesMixin,
    entitiesPbehaviorTypesMixin,
  ],
  methods: {
    showCreateTypeModal() {
      this.$modals.show({
        name: MODALS.createPbehaviorType,
        config: {
          action: data => this.createPbehaviorType({ data }),
        },
      });
    },
  },
};
</script>
