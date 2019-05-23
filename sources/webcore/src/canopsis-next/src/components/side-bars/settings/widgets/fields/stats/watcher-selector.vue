<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsWatcherSelector.title') }}
    v-container(fluid)
      v-btn(@click="showWatchersListModal") Select watcher
</template>

<script>
import { MODALS, WIDGET_TYPES, CONTEXT_ENTITIES_TYPES } from '@/constants';

import { generateWidgetByType } from '@/helpers/entities';

import modalMixin from '@/mixins/modal';
import formMixin from '@/mixins/form';

export default {
  mixins: [modalMixin, formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    stats: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    showWatchersListModal() {
      const widget = generateWidgetByType(WIDGET_TYPES.context);
      const filter = { $and: [{ _id: 'watcher_a5ecaa17-833c-4617-bf49-09b04f347880' }] };

      const watcherFilter = {
        title: 'watcher_a5ecaa17-833c-4617-bf49-09b04f347880',
        filter,
      };

      widget.parameters.mainFilter = watcherFilter;
      widget.parameters.viewFilters = [watcherFilter];

      const query = {
        typesFilter: { type: CONTEXT_ENTITIES_TYPES.watcher },
        mainFilter: watcherFilter.filter,
      };

      this.showModal({
        name: MODALS.contextEntitiesList,
        config: {
          widget,
          query,
        },
      });
    },
  },
};
</script>

