<template lang="pug">
  modal-wrapper(data-test="selectViewModal", close)
    template(slot="title")
      span {{ $t('modals.view.select.title') }}
    template(slot="text")
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

import entitiesViewsGroupsMixin from '@/mixins/entities/view/group';
import rightsEntitiesGroupMixin from '@/mixins/rights/entities/group';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.selectView,
  components: { ModalWrapper },
  mixins: [entitiesViewsGroupsMixin, rightsEntitiesGroupMixin],
  methods: {
    async selectView(viewId) {
      if (this.config.action) {
        await this.config.action(viewId);
      }

      this.$modals.hide();
    },
  },
};
</script>

