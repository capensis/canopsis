<template lang="pug">
  v-card(data-test="selectViewModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.view.select.title') }}
    v-card-text
      v-list.py-0(dark)
        v-list-group(v-for="group in availableGroups", :key="group._id")
          v-list-tile(
          slot="activator",
          :data-test="`selectView-group-${group._id}`"
          )
            v-list-tile-title {{ group.name }}
          v-list-tile(
          v-for="view in group.views",
          :key="view._id",
          :data-test="`selectView-view-${view._id}`",
          @click="selectView(view._id)"
          )
            v-list-tile-title.pl-2 {{ view.title }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesViewsGroupsMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';

export default {
  name: MODALS.selectView,
  mixins: [modalInnerMixin, entitiesViewsGroupsMixin, rightsEntitiesGroupMixin],
  methods: {
    async selectView(viewId) {
      if (this.config.action) {
        await this.config.action(viewId);
      }

      this.hideModal();
    },
  },
};
</script>

