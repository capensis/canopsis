<template lang="pug">
  v-card
    v-card-title.blue.darken-4.white--text
      h2 {{ $t('modals.view.select.title') }}
    v-card-text
      v-list.py-0(dark)
        v-list-group(v-for="group in groups", :key="group._id")
          v-list-tile(slot="activator")
            v-list-tile-title {{ group.name }}
          v-list-tile(v-for="view in group.views", :key="view._id", @click.stop="selectView(view._id)")
            v-list-tile-title.pl-2 {{ view.name }}
</template>

<script>
import { MODALS } from '@/constants';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import entitiesViewsGroupsMixin from '@/mixins/entities/view/group';

export default {
  name: MODALS.selectView,
  mixins: [modalInnerMixin, entitiesViewsGroupsMixin],
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

